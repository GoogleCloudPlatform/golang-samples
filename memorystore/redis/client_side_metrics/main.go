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
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	gcpmetric "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	gcptrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	tracer           trace.Tracer
	rttHist          metric.Float64Histogram
	clientBlockHist  metric.Float64Histogram
	appBlockHist     metric.Float64Histogram
	retryCounter     metric.Int64Counter
	connErrorCounter metric.Int64Counter
)

func initTelemetry(ctx context.Context) func() {
	traceExporter, err := gcptrace.New()
	if err != nil {
		log.Fatalf("Failed to create trace exporter: %v", err)
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(traceExporter))
	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("redigo.client")

	metricExporter, err := gcpmetric.New()
	if err != nil {
		log.Fatalf("Failed to create metric exporter: %v", err)
	}
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(10*time.Second))))
	otel.SetMeterProvider(mp)
	meter := mp.Meter("redigo.metrics")

	rttHist, _ = meter.Float64Histogram("redis_client_rtt", metric.WithUnit("ms"))
	clientBlockHist, _ = meter.Float64Histogram("redis_client_blocking_latency", metric.WithUnit("ms"))
	appBlockHist, _ = meter.Float64Histogram("redis_application_blocking_latency", metric.WithUnit("ms"))
	retryCounter, _ = meter.Int64Counter("redis_retry_count")
	connErrorCounter, _ = meter.Int64Counter("redis_connectivity_error_count")

	initAttrs := metric.WithAttributes(attribute.String("operation", "startup"))
	retryCounter.Add(ctx, 0, initAttrs)
	connErrorCounter.Add(ctx, 0, initAttrs)

	return func() {
		tp.Shutdown(ctx)
		mp.Shutdown(ctx)
	}
}

func smartRedisCall(ctx context.Context, pool *redis.Pool, operationName string, commandName string, args ...interface{}) (interface{}, error) {
	// Create a dedicated child span for the Redis command
	ctx, span := tracer.Start(ctx, operationName)
	span.SetAttributes(attribute.String("redis.command", commandName))
	defer span.End()

	maxRetries := 3
	attempt := 0
	metricOpts := metric.WithAttributes(attribute.String("operation", operationName))

	for attempt < maxRetries {
		poolStart := time.Now()
		conn := pool.Get()
		clientBlockHist.Record(ctx, float64(time.Since(poolStart).Milliseconds()), metricOpts)

		if err := conn.Err(); err != nil {
			conn.Close()
			connErrorCounter.Add(ctx, 1, metricOpts)
			retryCounter.Add(ctx, 1, metricOpts)
			span.RecordError(err) // Attach error to trace
			attempt++
			time.Sleep(time.Duration(100<<attempt) * time.Millisecond)
			continue
		}

		reqStart := time.Now()
		reply, err := conn.Do(commandName, args...)
		rttHist.Record(ctx, float64(time.Since(reqStart).Milliseconds()), metricOpts)
		conn.Close()

		if err != nil {
			retryCounter.Add(ctx, 1, metricOpts)
			span.RecordError(err) // Attach error to trace
			attempt++
			time.Sleep(time.Duration(100<<attempt) * time.Millisecond)
			continue
		}

		appStart := time.Now()
		_ = fmt.Sprintf("%v", reply)
		appBlockHist.Record(ctx, float64(time.Since(appStart).Milliseconds()), metricOpts)

		return reply, nil
	}
	return nil, fmt.Errorf("max retries reached for %s", operationName)
}

func main() {
	ctx := context.Background()
	shutdown := initTelemetry(ctx)
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

	ctx, span := tracer.Start(ctx, "fetch_data_span")
	defer span.End()

	// Simple write and read operations
	_, err := smartRedisCall(ctx, pool, "set_user", "SET", "user:123", "active")
	if err != nil {
		log.Printf("Error setting data: %v", err)
	}

	val, err := smartRedisCall(ctx, pool, "get_user", "GET", "user:123")
	if err != nil {
		log.Printf("Error fetching data: %v", err)
	} else {
		log.Printf("Retrieved value: %s", val)
	}
}

// [END memorystore_redis_client_side_metrics]
