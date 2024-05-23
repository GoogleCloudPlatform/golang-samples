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
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const name = "work"

var (
	meter          = otel.Meter(name)
	tracer         = otel.Tracer(name)
	sleepHistogram metric.Int64Histogram
)

func init() {
	var err error
	sleepHistogram, err = meter.Int64Histogram("example.sleep.histogram",
		metric.WithDescription("Sample histogram to measure time spent in sleeping"),
		metric.WithExplicitBucketBoundaries(50, 75, 100, 125, 150, 200),
		metric.WithUnit("ms"))

	if err != nil {
		panic(err)
	}
}

// randomSleep simulates a some job being triggerred in response to an API call to the server.
// This function computes 10 random values and records them into a histogram which can be
// later visualized as a distribution.
func randomSleep(r *http.Request) time.Duration {
	ctx, span := tracer.Start(r.Context(), "randomSleep")
	defer span.End()

	hostValue := attribute.String("host.value", r.Host)

	// simulate the work by sleeping
	sleepTime := time.Duration(100+rand.Intn(100)) * time.Millisecond
	time.Sleep(sleepTime)

	// record time slept
	sleepHistogram.Record(ctx, int64(sleepTime/time.Millisecond), metric.WithAttributes(hostValue))
	return sleepTime
}
