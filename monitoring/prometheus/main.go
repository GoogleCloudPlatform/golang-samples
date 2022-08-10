// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

// [START monitoring_sli_metrics_prometheus_setup]
import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// [END monitoring_sli_metrics_prometheus_setup]
// [START monitoring_sli_metrics_prometheus_create_metrics]
// Sets up metrics.
var (
	requestCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "go_request_count",
		Help: "total request count",
	})
	failedRequestCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "go_failed_request_count",
		Help: "failed request count",
	})
	responseLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "go_response_latency",
		Help: "response latencies",
	})
)

// [END monitoring_sli_metrics_prometheus_create_metrics]

func main() {
	http.HandleFunc("/", handle)
	// [START monitoring_sli_metrics_prometheus_metrics_endpoint]
	http.Handle("/metrics", promhttp.Handler())
	// [END monitoring_sli_metrics_prometheus_metrics_endpoint]
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	// [START monitoring_sli_metrics_prometheus_latency]
	requestReceived := time.Now()
	defer func() {
		responseLatency.Observe(time.Since(requestReceived).Seconds())
	}()
	// [END monitoring_sli_metrics_prometheus_latency]
	// [START monitoring_sli_metrics_prometheus_counts]
	requestCount.Inc()

	// Fails 10% of the time.
	if rand.Intn(100) >= 90 {
		log.Printf("intentional failure encountered")
		failedRequestCount.Inc()
		http.Error(w, "intentional error!", http.StatusInternalServerError)
		return
	}
	// [END monitoring_sli_metrics_prometheus_counts]
	delay := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(delay)
	fmt.Fprintf(w, "Succeeded after %v", delay)
	return
}
