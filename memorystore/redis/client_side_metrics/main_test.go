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

package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type mockRedisConn struct {
	dialCalls     int
	doCalls       int
	closeCalls    int
	activeConnErr error // Sticky connection error for the active socket lifecycle
	doErrs        []error
	doReplies     []interface{}
	doRepliesIdx  int
	doErrsIdx     int
}

func (m *mockRedisConn) Close() error {
	m.closeCalls++
	return nil
}

func (m *mockRedisConn) Err() error {
	return m.activeConnErr
}

func (m *mockRedisConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	// Ignore the internal Redigo pool Close sentinel command
	if commandName == "" {
		return nil, nil
	}

	m.doCalls++
	if m.doErrsIdx < len(m.doErrs) {
		err := m.doErrs[m.doErrsIdx]
		m.doErrsIdx++
		if err != nil {
			return nil, err
		}
	}
	if m.doRepliesIdx < len(m.doReplies) {
		reply := m.doReplies[m.doRepliesIdx]
		m.doRepliesIdx++
		return reply, nil
	}
	return "OK", nil
}

func (m *mockRedisConn) Send(commandName string, args ...interface{}) error { return nil }
func (m *mockRedisConn) Flush() error                                       { return nil }
func (m *mockRedisConn) Receive() (interface{}, error)                      { return nil, nil }

func setupTestTelemetry(t *testing.T) (*MetricClient, func()) {
	ctx := context.Background()
	tp := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	tracer := otel.Tracer("redis.test")

	mp := sdkmetric.NewMeterProvider()
	otel.SetMeterProvider(mp)
	meter := mp.Meter("redigo.test")

	rttHist, _ := meter.Float64Histogram("redis_client_rtt")
	clientBlockHist, _ := meter.Float64Histogram("redis_client_blocking_latency")
	appBlockHist, _ := meter.Float64Histogram("redis_application_blocking_latency")
	retryCounter, _ := meter.Int64Counter("redis_retry_count")
	connErrorCounter, _ := meter.Int64Counter("redis_connectivity_error_count")

	metricClient := &MetricClient{
		Tracer:           tracer,
		RTTBridge:        rttHist,
		ClientBlockBridge: clientBlockHist,
		AppBlockBridge:    appBlockHist,
		RetryCounter:     retryCounter,
		ConnErrorCounter: connErrorCounter,
	}

	cleanup := func() {
		tp.Shutdown(ctx)
		mp.Shutdown(ctx)
	}

	return metricClient, cleanup
}

func TestSmartRedisCallTable(t *testing.T) {
	client, cleanup := setupTestTelemetry(t)
	defer cleanup()

	tests := []struct {
		name            string
		operation       string
		command         string
		args            []interface{}
		dialErrs        []error
		connErrs        []error
		doErrs          []error
		doReplies       []interface{}
		expectErr       bool
		expectedErrSub  string
		expectedDials   int
		expectedCloses  int
		expectedDoCalls int
		expectedVal     interface{}
	}{
		{
			name:            "Success_SET_GET",
			operation:       "set_user",
			command:         "SET",
			args:            []interface{}{"user:123", "active"},
			dialErrs:        nil,
			connErrs:        nil,
			doErrs:          []error{nil},
			doReplies:       []interface{}{"OK"},
			expectErr:       false,
			expectedDials:   1,
			expectedCloses:  0, // Recycled to pool, no TCP socket close called
			expectedDoCalls: 1,
			expectedVal:     "OK",
		},
		{
			name:      "RetryOnDialError_ThenSuccess",
			operation: "get_user",
			command:   "GET",
			args:      []interface{}{"user:123"},
			dialErrs: []error{
				fmt.Errorf("connection refused attempt 1"),
				nil,
			},
			connErrs:        nil,
			doErrs:          []error{nil},
			doReplies:       []interface{}{"active"},
			expectErr:       false,
			expectedDials:   2,
			expectedCloses:  0, // Recycled to pool upon success
			expectedDoCalls: 1,
			expectedVal:     "active",
		},
		{
			name:      "RetryOnStaleConn_ThenSuccess",
			operation: "get_user",
			command:   "GET",
			args:      []interface{}{"user:123"},
			dialErrs:  nil,
			connErrs: []error{
				fmt.Errorf("stale socket read error attempt 1"),
				nil,
			},
			doErrs:          []error{nil},
			doReplies:       []interface{}{"active"},
			expectErr:       false,
			expectedDials:   2,
			expectedCloses:  1, // 1 for the dead pool checkout (destroyed!), 0 for the final success (recycled)
			expectedDoCalls: 1,
			expectedVal:     "active",
		},
		{
			name:      "PermanentFailure_AfterMaxRetries",
			operation: "set_user",
			command:   "SET",
			args:      []interface{}{"user:123", "active"},
			dialErrs: []error{
				fmt.Errorf("redis cluster offline A"),
				fmt.Errorf("redis cluster offline B"),
				fmt.Errorf("redis cluster offline C"),
			},
			connErrs:        nil,
			doErrs:          nil,
			expectErr:       true,
			expectedErrSub:  "max retries reached for set_user",
			expectedDials:   3,
			expectedCloses:  0, // Dial failures mean GetContext returns nil (no TCP sockets to close)
			expectedDoCalls: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			sharedConn := &mockRedisConn{
				doErrs:    tc.doErrs,
				doReplies: tc.doReplies,
			}

			pool := &redis.Pool{
				MaxIdle:     10,
				IdleTimeout: 240 * time.Second,
				Dial: func() (redis.Conn, error) {
					sharedConn.dialCalls++
					if sharedConn.dialCalls <= len(tc.dialErrs) {
						err := tc.dialErrs[sharedConn.dialCalls-1]
						if err != nil {
							return nil, err
						}
					}
					// Update the sticky active connection error for this socket's lifecycle
					if sharedConn.dialCalls <= len(tc.connErrs) {
						sharedConn.activeConnErr = tc.connErrs[sharedConn.dialCalls-1]
					} else {
						sharedConn.activeConnErr = nil
					}
					return sharedConn, nil
				},
			}

			val, err := client.smartRedisCall(ctx, tc.operation, pool, tc.command, tc.args...)

			if tc.expectErr {
				if err == nil {
					t.Errorf("%s: expected a permanent failure error but it returned nil", tc.name)
				} else if !strings.Contains(err.Error(), tc.expectedErrSub) {
					t.Errorf("%s: expected error to contain '%s', got '%v'", tc.name, tc.expectedErrSub, err)
				}
			} else {
				if err != nil {
					t.Errorf("%s: unexpected error: %v", tc.name, err)
				}
				if val != tc.expectedVal {
					t.Errorf("%s: expected value %v, got %v", tc.name, tc.expectedVal, val)
				}
			}

			if sharedConn.dialCalls != tc.expectedDials {
				t.Errorf("%s: expected %d Dial attempts, got %d", tc.name, tc.expectedDials, sharedConn.dialCalls)
			}

			if sharedConn.doCalls != tc.expectedDoCalls {
				t.Errorf("%s: expected %d Do calls, got %d", tc.name, tc.expectedDoCalls, sharedConn.doCalls)
			}

			if sharedConn.closeCalls != tc.expectedCloses {
				t.Errorf("%s: expected %d Close calls, got %d", tc.name, tc.expectedCloses, sharedConn.closeCalls)
			}
		})
	}
}