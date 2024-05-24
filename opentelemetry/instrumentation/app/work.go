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
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const scopeName = "github.com/GoogleCloudPlatform/golang-samples/opentelemetry/instrumentation/app/work"

var (
	meter                = otel.Meter(scopeName)
	tracer               = otel.Tracer(scopeName)
	sleepHistogram       metric.Int64Histogram
	subRequestsHistogram metric.Int64Histogram
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

	subRequestsHistogram, err = meter.Int64Histogram("example.subrequests.histogram",
		metric.WithDescription("Sample histogram to measure time spent in sleeping"),
		metric.WithExplicitBucketBoundaries(1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
		metric.WithUnit("1"))
	if err != nil {
		panic(err)
	}
}

// randomSleep simulates a some job being triggerred in response to an API call to the server.
// This function records the time spent in sleeping in a histogram which can later be
// visualized as a distribution.
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

// computeSubrequests performs the task of making 3-7 http of requests to /single endpoint on localhost:8080.
// This function records the number of subrequests made in a histogram which can later be visualized
// as a distribution.
func computeSubrequests(r *http.Request) error {
	ctx, span := tracer.Start(r.Context(), "subrequests")
	defer span.End()

	subRequests := 3 + rand.Intn(4)
	// Write a structured log with the request context, which allows the log to
	// be linked with the trace for this request.
	slog.InfoContext(ctx, "computing multiple requests", slog.Int("subRequests", subRequests))

	// Make 3-7 http requests to the /single endpoint.
	for i := 0; i < subRequests; i++ {
		if err := callSingle(ctx); err != nil {
			return err
		}
	}

	// record number of sub-requests made
	subRequestsHistogram.Record(ctx, int64(subRequests))
	return nil
}
