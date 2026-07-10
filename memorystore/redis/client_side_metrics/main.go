// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
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
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// MetricClient encapsulates OpenTelemetry handlers and the Redis connection pool
type MetricClient struct {
	Tracer           trace.Tracer
	RTTBridge        metric.Float64Histogram
	ClientBlockBridge metric.Float64Histogram
	AppBlockBridge    metric.Float64Histogram
	RetryCounter     metric.Int64Counter
	ConnErrorCounter metric.Int64Counter
	Pool             *redis.Pool
}

// sinceMs calculates elapsed time in fractional milliseconds to preserve sub-millisecond measurements.
func sinceMs(start time.Time) float64 {
	return float64(time.Since(start).Microseconds()) / 1000.0
}

func initTelemetry(ctx context.Context) (*MetricClient, *sdktrace.TracerProvider, *sdkmetric.MeterProvider, error) {
	// 1. Setup Tracing
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
	)
	otel.SetTracerProvider(tp)
	tracer := otel.Tracer("redis.client")

	// 2. Setup Metrics
	metricExporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(10*time.Second))),
	)
	otel.SetMeterProvider(mp)
	meter := mp.Meter("redigo.metrics")

	rttHist, err := meter.Float64Histogram("redis_client_rtt", metric.WithUnit("ms"))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create rttHist: %w", err)
	}

	clientBlockHist, err := meter.Float64Histogram("redis_client_blocking_latency", metric.WithUnit("ms"))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create clientBlockHist: %w", err)
	}

	appBlockHist, err := meter.Float64Histogram("redis_application_blocking_latency", metric.WithUnit("ms"))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create appBlockHist: %w", err)
	}

	retryCounter, err := meter.Int64Counter("redis_retry_count")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create retryCounter: %w", err)
	}

	connErrorCounter, err := meter.Int64Counter("redis_connectivity_error_count")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create connErrorCounter: %w", err)
	}

	metricOpts := metric.WithAttributes(attribute.String("operation", "startup"))
	retryCounter.Add(ctx, 0, metricOpts)
	connErrorCounter.Add(ctx, 0, metricOpts)

	client := &MetricClient{
		Tracer:           tracer,
		RTTBridge:        rttHist,
		ClientBlockBridge: clientBlockHist,
		AppBlockBridge:    appBlockHist,
		RetryCounter:     retryCounter,
		ConnErrorCounter: connErrorCounter,
	}

	return client, tp, mp, nil
}

func initRedisPool() *redis.Pool {
	redisHost := os.Getenv("REDISHOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDISPORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", redisHost, redisPort))
		},
	}
}

func (c *MetricClient) smartRedisCall(ctx context.Context, operationName string, pool *redis.Pool, commandName string, args ...interface{}) (interface{}, error) {
	maxRetries := 3
	attempt := 0
	metricOpts := metric.WithAttributes(attribute.String("operation", operationName))
	var lastErr error

	span := trace.SpanFromContext(ctx)

	for attempt < maxRetries {
		poolStart := time.Now()
		// Get connection using context to avoid blocking indefinitely
		conn, err := pool.GetContext(ctx)
		c.ClientBlockBridge.Record(ctx, sinceMs(poolStart), metricOpts)

		if err != nil {
			c.ConnErrorCounter.Add(ctx, 1, metricOpts)
			c.RetryCounter.Add(ctx, 1, metricOpts)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			lastErr = err
			attempt++
			if attempt >= maxRetries {
				break
			}
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
			continue
		}

		// Check if the connection is dead
		if err := conn.Err(); err != nil {
			conn.Close()
			c.ConnErrorCounter.Add(ctx, 1, metricOpts)
			c.RetryCounter.Add(ctx, 1, metricOpts)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			lastErr = err
			attempt++
			if attempt >= maxRetries {
				break
			}
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
			continue
		}

		reqStart := time.Now()
		reply, err := conn.Do(commandName, args...)
		c.RTTBridge.Record(ctx, sinceMs(reqStart), metricOpts)

		if err != nil {
			conn.Close()
			c.ConnErrorCounter.Add(ctx, 1, metricOpts)
			c.RetryCounter.Add(ctx, 1, metricOpts)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			lastErr = err
			attempt++
			if attempt >= maxRetries {
				break
			}
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
			continue
		}

		appStart := time.Now()
		// Simulate application processing overhead (replacing FMT formatting)
		time.Sleep(2 * time.Millisecond)
		c.AppBlockBridge.Record(ctx, sinceMs(appStart), metricOpts)

		// Reset span status to Ok if the operation eventually succeeds
		span.SetStatus(codes.Ok, "")

		conn.Close()
		return reply, nil
	}
	return nil, fmt.Errorf("max retries reached for %s: %w", operationName, lastErr)
}

func main() {
	ctx := context.Background()
	client, tp, mp, err := initTelemetry(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}
	defer tp.Shutdown(ctx)
	defer mp.Shutdown(ctx)

	pool := initRedisPool()
	defer pool.Close()

	if client.Tracer != nil {
		var span trace.Span
		ctx, span = client.Tracer.Start(ctx, "process_user_span")
		defer span.End()

		trySet, err := client.smartRedisCall(ctx, "set_user", pool, "SET", "user:123", "active")
		if err != nil {
			fmt.Printf("Error setting key: %v\n", err)
		} else {
			fmt.Printf("Set Response: %v\n", trySet)
		}

		result, err := client.smartRedisCall(ctx, "get_user", pool, "GET", "user:123")
		if err != nil {
			fmt.Printf("Error getting key: %v\n", err)
		} else {
			fmt.Printf("Retrieved: %v\n", result)
		}
	}
}
// [END memorystore_redis_client_side_metrics]