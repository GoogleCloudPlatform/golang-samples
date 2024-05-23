// Copyright 2024 Google LLC
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
	tracer        = otel.Tracer(name)
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

// doWork simulates a some job being triggerred in response to an API call to the server.
// This function computes 10 random values and records them into a histogram which can be
// later visualized as a distribution.
func doWork(ctx context.Context, host string) time.Duration {
	start := time.Now()
	hostValue := attribute.String("host.value", host)

	// simulate the overall work by sleeping
	sleepTime := time.Duration(100+rand.Intn(100)) * time.Millisecond
	time.Sleep(sleepTime)

	// wrap the random number generation in a span - to better visualize the time spent in this part
	traceCtx, span := tracer.Start(ctx, "doWork")
	for i := 0; i < 10; i++ {
		randomNum := rand.Intn(100)
		workHistogram.Record(traceCtx, int64(randomNum), metric.WithAttributes(hostValue))
	}
	span.End()

	elapsedTime := time.Since(start)
	return elapsedTime
}
