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

package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// mockRedisConn is a dependency-free fake implementing the redis.Conn interface
type mockRedisConn struct {
	data          map[string]interface{}
	errOnDo       error
	errOnDoCount  int
	errOnErr      error
	errOnErrCount int
	closeCalls    int
	doCalls       int
	dialCalls     int
}

func (m *mockRedisConn) Close() error {
	m.closeCalls++
	return nil
}

func (m *mockRedisConn) Err() error {
	if m.errOnErrCount > 0 {
		m.errOnErrCount--
		return m.errOnErr
	}
	return nil
}

func (m *mockRedisConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	// Ignore Redigo's internal Close() sentinel to keep test counters accurate
	if commandName == "" {
		return nil, nil
	}

	m.doCalls++
	if m.errOnDoCount > 0 {
		m.errOnDoCount--
		return nil, m.errOnDo
	}

	if commandName == "SET" && len(args) >= 2 {
		key := fmt.Sprintf("%v", args[0])
		val := args[1]
		if m.data == nil {
			m.data = make(map[string]interface{})
		}
		m.data[key] = val
		return "OK", nil
	}

	if commandName == "GET" && len(args) >= 1 {
		key := fmt.Sprintf("%v", args[0])
		val, ok := m.data[key]
		if !ok {
			return nil, nil
		}
		return val, nil
	}

	return nil, nil
}

func (m *mockRedisConn) Send(commandName string, args ...interface{}) error {
	return nil
}

func (m *mockRedisConn) Flush() error {
	return nil
}

func (m *mockRedisConn) Receive() (interface{}, error) {
	return nil, nil
}

// initTestTelemetry safely binds a MetricClient to hermetic No-op implementations
func initTestTelemetry() *MetricClient {
	tp := trace.NewNoopTracerProvider()

	meter := otel.GetMeterProvider().Meter("noop")
	rttHist, _ := meter.Float64Histogram("redis_client_rtt", metric.WithUnit("ms"))
	clientBlockHist, _ := meter.Float64Histogram("redis_client_blocking_latency", metric.WithUnit("ms"))
	appBlockHist, _ := meter.Float64Histogram("redis_application_blocking_latency", metric.WithUnit("ms"))
	retryCounter, _ := meter.Int64Counter("redis_retry_count")
	connErrorCounter, _ := meter.Int64Counter("redis_connectivity_error_count")

	// Disable sleeping during unit tests to make the test suite lightning fast
	sleep = func(d time.Duration) {}

	return &MetricClient{
		tracer:           tp.Tracer("noop"),
		rttHist:          rttHist,
		clientBlockHist:  clientBlockHist,
		appBlockHist:     appBlockHist,
		retryCounter:     retryCounter,
		connErrorCounter: connErrorCounter,
	}
}

type redisOperation struct {
	opName      string
	command     string
	args        []interface{}
	expectedVal interface{}
	expectErr   bool
}

func TestSmartRedisCallTable(t *testing.T) {
	client := initTestTelemetry()

	tests := []struct {
		name            string
		maxIdle         int
		errOnDo         error
		errOnDoCount    int
		errOnErr        error
		errOnErrCount   int
		errOnDial       error
		errOnDialCount  int
		operations      []redisOperation
		expectedDials   int
		expectedDoCalls int
		expectedErrSub  string
	}{
		{
			name:    "Success_SET_and_GET",
			maxIdle: 1,
			operations: []redisOperation{
				{opName: "set_user", command: "SET", args: []interface{}{"user:123", "active"}, expectedVal: "OK", expectErr: false},
				{opName: "get_user", command: "GET", args: []interface{}{"user:123"}, expectedVal: "active", expectErr: false},
			},
			expectedDials:   1,
			expectedDoCalls: 2,
		},
		{
			name:          "RetryOnConnError_ThenSuccess",
			maxIdle:       0,
			errOnErr:      errors.New("connection reset by peer"),
			errOnErrCount: 1,
			operations: []redisOperation{
				{opName: "get_user", command: "GET", args: []interface{}{"user:123"}, expectedVal: "active", expectErr: false},
			},
			expectedDials:   2,
			expectedDoCalls: 1,
		},
		{
			name:         "PermanentFailure_AfterMaxRetries",
			maxIdle:      0,
			errOnDo:      errors.New("Redis cluster unavailable permanently"),
			errOnDoCount: 3,
			operations: []redisOperation{
				{opName: "set_user", command: "SET", args: []interface{}{"user:123", "active"}, expectedVal: nil, expectErr: true},
			},
			expectedDials:   3,
			expectedDoCalls: 3,
			expectedErrSub:  "max retries reached for set_user",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			sharedConn := &mockRedisConn{
				data:          make(map[string]interface{}),
				errOnDo:       tc.errOnDo,
				errOnDoCount:  tc.errOnDoCount,
				errOnErr:      tc.errOnErr,
				errOnErrCount: tc.errOnErrCount,
			}

			// Pre-populate data for GET operations following connection retries
			if tc.name == "RetryOnConnError_ThenSuccess" {
				sharedConn.data["user:123"] = "active"
			}

			dialErrsLeft := tc.errOnDialCount
			pool := &redis.Pool{
				MaxIdle: tc.maxIdle,
				Dial: func() (redis.Conn, error) {
					sharedConn.dialCalls++
					if tc.errOnDial != nil && dialErrsLeft > 0 {
						dialErrsLeft--
						return nil, tc.errOnDial
					}
					return sharedConn, nil
				},
			}
			defer pool.Close()

			for _, op := range tc.operations {
				val, err := client.smartRedisCall(ctx, pool, op.opName, op.command, op.args...)

				if op.expectErr {
					if err == nil {
						t.Errorf("%s: expected a permanent failure error but it returned nil", tc.name)
					} else if tc.expectedErrSub != "" && !strings.Contains(err.Error(), tc.expectedErrSub) {
						t.Errorf("%s: expected error to contain '%s', got '%v'", tc.name, tc.expectedErrSub, err)
					}
					continue
				}

				if err != nil {
					t.Fatalf("%s: unexpected smartRedisCall error: %v", tc.name, err)
				}

				if val != op.expectedVal {
					t.Errorf("%s: smartRedisCall return got %v, want %v", tc.name, val, op.expectedVal)
				}
			}

			if tc.expectedDials > 0 && sharedConn.dialCalls != tc.expectedDials {
				t.Errorf("%s: expected %d Dial attempts, got %d", tc.name, tc.expectedDials, sharedConn.dialCalls)
			}

			if tc.expectedDoCalls > 0 && sharedConn.doCalls != tc.expectedDoCalls {
				t.Errorf("%s: expected %d Do calls, got %d", tc.name, tc.expectedDoCalls, sharedConn.doCalls)
			}
		})
	}
}
