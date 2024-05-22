package main

import (
	"context"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const name = "work"

var (
	meter   = otel.Meter(name)
	workCnt metric.Int64Counter
)

func init() {
	var err error
	workCnt, err = meter.Int64Counter("example.counter",
		metric.WithDescription("Processed Jobs"),
		metric.WithUnit("{job_processed}"))
	if err != nil {
		panic(err)
	}
}

func doWork(ctx context.Context, host string) time.Duration {
	start := time.Now()
	sleepTime := time.Duration(100+rand.Intn(100)) * time.Millisecond
	time.Sleep(sleepTime)

	jobAttrs := attribute.String("host.value", host)
	workCnt.Add(ctx, 1, metric.WithAttributes(jobAttrs))

	elapsedTime := time.Since(start)
	return elapsedTime
}
