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
	meter         = otel.Meter(name)
	workHistogram metric.Int64Histogram
)

func init() {
	var err error
	workHistogram, err = meter.Int64Histogram("example.histogram",
		metric.WithDescription("Sample histogram"),
		metric.WithUnit("1"))

	if err != nil {
		panic(err)
	}
}

func doWork(ctx context.Context, host string) time.Duration {
	start := time.Now()

	hostValue := attribute.String("host.value", host)

	sleepTime := time.Duration(100+rand.Intn(100)) * time.Millisecond
	time.Sleep(sleepTime)
	for i := 0; i < 10; i++ {
		randomNum := rand.Intn(100)
		workHistogram.Record(ctx, int64(randomNum), metric.WithAttributes(hostValue))
	}

	elapsedTime := time.Since(start)
	return elapsedTime
}
