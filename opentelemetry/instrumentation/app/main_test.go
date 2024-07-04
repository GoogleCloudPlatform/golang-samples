// Copyright 2023 Google LLC
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
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// OTEL_BSP_SCHEDULE_DELAY, OTEL_METRIC_EXPORT_INTERVAL
const timeoutSeconds = 3

func TestWriteTelemetry(t *testing.T) {
	ctx := context.Background()
	// Tests build.
	m := testutil.BuildMain(t)
	if !m.Built() {
		t.Fatalf("failed to build app")
	}
	// Listen for telemetry from the application before starting it, so the
	// application will successfully connect right away during the test.
	listener, err := net.Listen("tcp", "localhost:")
	if err != nil {
		t.Fatal(err)
	}

	// Handle traces from the application by comparing them against
	// expectations
	ts := &traceServer{
		t:               t,
		expectationsMet: make(chan struct{}),
		expectations: []*spanExpectation{
			{name: "/multi"},
			{name: "/single"},
		},
	}
	http.HandleFunc("/v1/traces", ts.handleTraces)
	// Handle metrics from the application by comparing them against
	// expectations
	ms := &metricsServer{
		t:               t,
		expectationsMet: make(chan struct{}),
		expectations: []*metricExpectation{
			{name: "http.server.duration"},
		},
	}
	http.HandleFunc("/v1/metrics", ms.handleMetrics)
	srv := &http.Server{}
	go srv.Serve(listener)
	defer srv.Shutdown(ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	// Run the application, and verify stdout logs contain expected structured logs
	go func() {
		defer wg.Done()
		stdout, stderr, _ := m.Run(map[string]string{
			// Point the OTLP exporter at our test handler
			"OTEL_EXPORTER_OTLP_ENDPOINT": "http://" + listener.Addr().String(),
			"OTEL_EXPORTER_OTLP_INSECURE": "true",
			// Export metrics and traces after 1 second
			"OTEL_BSP_SCHEDULE_DELAY":     "1000",
			"OTEL_METRIC_EXPORT_INTERVAL": "1000",
		}, timeoutSeconds*time.Second)
		t.Logf("stdout: %v", string(stdout))
		t.Logf("stderr: %v", string(stderr))
		verifyStdoutLogs(t, stdout)
	}()
	// Send requests to our application to generate telemetry.
	testutil.Retry(t, 2*timeoutSeconds, 500*time.Millisecond, func(r *testutil.R) {
		resp, err := http.Get("http://localhost:8080/multi")
		if err != nil {
			r.Errorf(err.Error())
			r.Fail()
		} else if resp == nil {
			r.Errorf("status code was nil")
			r.Fail()
		} else if resp.StatusCode != 200 {
			r.Errorf("expected status code to be 200, was %v", resp.StatusCode)
			r.Fail()
		}
	})
	wg.Add(1)
	go func() {
		// Wait for the expected OTLP traces to be sent
		select {
		case <-time.After(timeoutSeconds * time.Second):
			t.Error("Timeout waiting for traces")
		case <-ts.expectationsMet:
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		// Wait for the expected OTLP metrics to be sent
		select {
		case <-time.After(timeoutSeconds * time.Second):
			t.Error("Timeout waiting for metrics")
		case <-ms.expectationsMet:
		}
		wg.Done()
	}()
	// Wait for all checks to complete
	wg.Wait()
}

type expectedLogFormat struct {
	Timestamp    string `json:"timestamp"`
	Severity     string `json:"severity"`
	Message      string `json:"message"`
	SubRequests  int    `json:"subRequests"`
	TraceID      string `json:"logging.googleapis.com/trace"`
	SpanID       string `json:"logging.googleapis.com/spanId"`
	TraceSampled bool   `json:"logging.googleapis.com/trace_sampled"`
}

const expectedLogMessage = "handle /multi request"

// verifyStdoutLogs checks for the expected log message in json format, and
// ensures all expected fields are set.
func verifyStdoutLogs(t *testing.T, stdout []byte) {
	for _, line := range strings.Split(string(stdout), "\n") {
		var contents expectedLogFormat
		if err := json.Unmarshal([]byte(line), &contents); err != nil {
			t.Error(err)
			continue
		}
		t.Logf("stdout line: %v; parsed: %+v", line, contents)
		if contents.Message != expectedLogMessage {
			continue
		}
		assert.NotEmpty(t, contents.Timestamp)
		assert.NotEmpty(t, contents.Severity)
		assert.NotEmpty(t, contents.Message)
		assert.NotEmpty(t, contents.SubRequests)
		assert.NotEmpty(t, contents.TraceID)
		assert.NotEmpty(t, contents.SpanID)
		assert.NotEmpty(t, contents.TraceSampled)
		return
	}
	t.Errorf("Did not find log message: %v", expectedLogMessage)
}

// traceServer implements OTLP trace receiver interfaces
type traceServer struct {
	t *testing.T

	lock            sync.Mutex
	expectations    []*spanExpectation
	expectationsMet chan struct{}
}

type spanExpectation struct {
	met  bool
	name string
}

func (t *traceServer) handleTraces(w http.ResponseWriter, req *http.Request) {
	t.lock.Lock()
	defer t.lock.Unlock()
	body, err := readAndCloseBody(req)
	if err != nil {
		t.t.Error(err)
		return
	}
	r, err := unmarshalTracesRequest(body)
	if err != nil {
		t.t.Error(err)
		return
	}
	t.t.Logf("trace request: %+v", r.Traces())
	// Iterate through each span sent, and update expectations with whether
	// they have been met.
	resourceSpans := r.Traces().ResourceSpans()
	for i := 0; i < resourceSpans.Len(); i++ {
		resourceSpan := resourceSpans.At(i)
		for j := 0; j < resourceSpan.ScopeSpans().Len(); j++ {
			scopeSpan := resourceSpan.ScopeSpans().At(j)
			for k := 0; k < scopeSpan.Spans().Len(); k++ {
				span := scopeSpan.Spans().At(k)
				allMet := true
				for _, expected := range t.expectations {
					expected.update(span)
					if !expected.met {
						allMet = false
					}
				}
				// If all expectations are met, notify the test.
				if allMet {
					select {
					case <-t.expectationsMet:
						// expectationsMet is already closed
					default:
						close(t.expectationsMet)
					}
				}
			}
		}
	}
}

func unmarshalTracesRequest(buf []byte) (ptraceotlp.ExportRequest, error) {
	req := ptraceotlp.NewExportRequest()
	err := req.UnmarshalProto(buf)
	return req, err
}

func readAndCloseBody(req *http.Request) ([]byte, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	if err = req.Body.Close(); err != nil {
		return nil, err
	}
	return body, nil
}

func (s *spanExpectation) update(span ptrace.Span) {
	if span.Name() != s.name {
		return
	}
	s.met = true
}

// metricsServer implements OTLP metrics receiver interfaces
type metricsServer struct {
	t *testing.T

	lock            sync.Mutex
	expectations    []*metricExpectation
	expectationsMet chan struct{}
}

type metricExpectation struct {
	met  bool
	name string
}

func (s *metricExpectation) update(metric pmetric.Metric) {
	if metric.Name() != s.name {
		return
	}
	s.met = true
}

func (t *metricsServer) handleMetrics(w http.ResponseWriter, req *http.Request) {
	t.lock.Lock()
	defer t.lock.Unlock()
	body, err := readAndCloseBody(req)
	if err != nil {
		t.t.Error(err)
		return
	}
	r, err := unmarshalMetricsRequest(body)
	if err != nil {
		t.t.Error(err)
		return
	}
	t.t.Logf("metrics request: %+v", r.Metrics())

	// Iterate through each metric sent, and update expectations with whether
	// they have been met.
	resourceMetrics := r.Metrics().ResourceMetrics()
	for i := 0; i < resourceMetrics.Len(); i++ {
		resourceMetric := resourceMetrics.At(i)
		for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
			scopeMetric := resourceMetric.ScopeMetrics().At(j)
			for k := 0; k < scopeMetric.Metrics().Len(); k++ {
				metric := scopeMetric.Metrics().At(k)
				t.t.Logf("metric: %+v", metric.Name())
				allMet := true
				for _, expected := range t.expectations {
					expected.update(metric)
					if !expected.met {
						allMet = false
					}
				}
				// If all expectations are met, notify the test.
				if allMet {
					select {
					case <-t.expectationsMet:
						// expectationsMet is already closed
					default:
						close(t.expectationsMet)
					}
				}
			}
		}
	}
}

func unmarshalMetricsRequest(buf []byte) (pmetricotlp.ExportRequest, error) {
	req := pmetricotlp.NewExportRequest()
	err := req.UnmarshalProto(buf)
	return req, err
}
