// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START memorystore_redis_client_side_metrics]
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	gcpmetric "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	gcptrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// MetricClient encapsulates the tracer and metric histograms to avoid package-level globals.
type MetricClient struct {
	tracer           trace.Tracer
	rttHist          metric.Float64Histogram
	clientBlockHist  metric.Float64Histogram
	appBlockHist     metric.Float64Histogram
	retryCounter     metric.Int64Counter
	connErrorCounter metric.Int64Counter
}

// sleep hook enables lightning-fast unit tests by stubbing out real time.Sleep
var sleep = time.Sleep

// sinceMs calculates elapsed time in fractional milliseconds to avoid truncating sub-millisecond durations.
func sinceMs(start time.Time) float64 {
	return float64(time.Since(start).Microseconds()) / 1000.0
}

func initTelemetry(ctx context.Context) (*MetricClient, func(), error) {
	traceExporter, err := gcptrace.New()
	if err != nil {
		return nil, nil, fmt.Errorf("gcptrace.New: %w", err)
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(traceExporter))
	otel.SetTracerProvider(tp)
	tracer := tp.Tracer("redigo.client")

	metricExporter, err := gcpmetric.New()
	if err != nil {
		return nil, nil, fmt.Errorf("gcpmetric.New: %w", err)
	}
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(10*time.Second))))
	otel.SetMeterProvider(mp)
	meter := mp.Meter("redigo.metrics")

	rttHist, err := meter.Float64Histogram("redis_client_rtt", metric.WithUnit("ms"))
	if err != nil {
		return nil, nil, fmt.Errorf("redis_client_rtt histogram: %w", err)
	}
	clientBlockHist, err := meter.Float64Histogram("redis_client_blocking_latency", metric.WithUnit("ms"))
	if err != nil {
		return nil, nil, fmt.Errorf("redis_client_blocking_latency histogram: %w", err)
	}
	appBlockHist, err := meter.Float64Histogram("redis_application_blocking_latency", metric.WithUnit("ms"))
	if err != nil {
		return nil, nil, fmt.Errorf("redis_application_blocking_latency histogram: %w", err)
	}
	retryCounter, err := meter.Int64Counter("redis_retry_count")
	if err != nil {
		return nil, nil, fmt.Errorf("redis_retry_count counter: %w", err)
	}
	connErrorCounter, err := meter.Int64Counter("redis_connectivity_error_count")
	if err != nil {
		return nil, nil, fmt.Errorf("redis_connectivity_error_count counter: %w", err)
	}

	client := &MetricClient{
		tracer:           tracer,
		rttHist:          rttHist,
		clientBlockHist:  clientBlockHist,
		appBlockHist:     appBlockHist,
		retryCounter:     retryCounter,
		connErrorCounter: connErrorCounter,
	}

	initAttrs := metric.WithAttributes(attribute.String("operation", "startup"))
	client.retryCounter.Add(ctx, 0, initAttrs)
	client.connErrorCounter.Add(ctx, 0, initAttrs)

	shutdown := func() {
		tp.Shutdown(ctx)
		mp.Shutdown(ctx)
	}

	return client, shutdown, nil
}

func (c *MetricClient) smartRedisCall(ctx context.Context, pool *redis.Pool, operationName string, commandName string, args ...interface{}) (interface{}, error) {
	// Create a dedicated child span for the Redis command
	ctx, span := c.tracer.Start(ctx, operationName)
	span.SetAttributes(attribute.String("redis.command", commandName))
	defer span.End()

	maxRetries := 3
	attempt := 0
	metricOpts := metric.WithAttributes(attribute.String("operation", operationName))
	var lastErr error

	for attempt < maxRetries {
		poolStart := time.Now()
		// Use GetContext to respect context deadlines and cancellation
		conn, err := pool.GetContext(ctx)
		c.clientBlockHist.Record(ctx, sinceMs(poolStart), metricOpts)

		if err != nil {
			c.connErrorCounter.Add(ctx, 1, metricOpts)
			c.retryCounter.Add(ctx, 1, metricOpts)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			lastErr = err
			attempt++
			if attempt >= maxRetries {
				break
			}
			sleep(time.Duration(100<<attempt) * time.Millisecond)
			continue
		}

		// Check if the connection is dead
		if err := conn.Err(); err != nil {
			conn.Close()
			c.connErrorCounter.Add(ctx, 1, metricOpts)
			c.retryCounter.Add(ctx, 1, metricOpts)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			lastErr = err
			attempt++
			if attempt >= maxRetries {
				break
			}
			sleep(time.Duration(100<<attempt) * time.Millisecond)
			continue
		}

		reqStart := time.Now()
		// Redigo has no native DoContext; pass timeouts using redis.DoWithTimeout when context has a deadline
		var reply interface{}
		if deadline, ok := ctx.Deadline(); ok {
			reply, err = redis.DoWithTimeout(conn, time.Until(deadline), commandName, args...)
		} else {
			reply, err = conn.Do(commandName, args...)
		}
		c.rttHist.Record(ctx, sinceMs(reqStart), metricOpts)
		conn.Close()

		if err != nil {
			c.retryCounter.Add(ctx, 1, metricOpts)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			lastErr = err
			attempt++
			if attempt >= maxRetries {
				break
			}
			sleep(time.Duration(100<<attempt) * time.Millisecond)
			continue
		}

		appStart := time.Now()
		// Replace fmt.Sprintf to remove unnecessary string formatting overhead
		sleep(2 * time.Millisecond)
		c.appBlockHist.Record(ctx, sinceMs(appStart), metricOpts)

		// Reset span status to Ok if the retry or execution eventually succeeds
		span.SetStatus(codes.Ok, "")

		return reply, nil
	}
	return nil, fmt.Errorf("max retries reached for %s: %w", operationName, lastErr)
}

func main() {
	ctx := context.Background()
	client, shutdown, err := initTelemetry(ctx)
	if err != nil {
		log.Printf("Failed to initialize telemetry: %v", err)
		os.Exit(1)
	}
	defer shutdown()

	redisHost := os.Getenv("REDISHOST")
	redisPort := os.Getenv("REDISPORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	pool := &redis.Pool{
		MaxIdle:     10,
		MaxActive:   20,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", redisHost, redisPort))
		},
	}
	defer pool.Close()

	ctx, span := client.tracer.Start(ctx, "fetch_data_span")
	defer span.End()

	// Simple write and read operations
	_, err = client.smartRedisCall(ctx, pool, "set_user", "SET", "user:123", "active")
	if err != nil {
		log.Printf("Error setting data: %v", err)
	}
	val, err := client.smartRedisCall(ctx, pool, "get_user", "GET", "user:123")
	if err != nil {
		log.Printf("Error fetching data: %v", err)
	} else {
		log.Printf("Retrieved value: %s", val)
	}
}

// [END memorystore_redis_client_side_metrics]
