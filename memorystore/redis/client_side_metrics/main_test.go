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

// redisOp represents a single smartRedisCall invocation within a test case.
type redisOp struct {
	operation string
	command   string
	args      []interface{}
	wantReply interface{}
	wantErr   bool
}

func TestSmartRedisCall(t *testing.T) {
	tests := []struct {
		name           string
		maxIdle        int
		errOnErr       error
		errOnErrCount  int
		errOnDo        error
		errOnDoCount   int
		presetData     map[string]interface{}
		ops            []redisOp
		checkDial      bool
		wantDialCalls  int
		checkDo        bool
		wantDoCalls    int
		checkClose     bool
		wantCloseCalls int
	}{
		{
			name:    "Success_SET_then_GET",
			maxIdle: 1, // Redigo will pool and reuse this connection
			ops: []redisOp{
				{operation: "set_user", command: "SET", args: []interface{}{"user:123", "active"}, wantReply: "OK"},
				{operation: "get_user", command: "GET", args: []interface{}{"user:123"}, wantReply: "active"},
			},
			// 1 Dial since the GET reuses the SET's idle connection
			checkDial:     true,
			wantDialCalls: 1,
			// Pool intercept means Close is deferred until pool.Close()
			checkClose:     true,
			wantCloseCalls: 0,
		},
		{
			name:          "RetryOnConnError",
			maxIdle:       0, // Disable pooling to accurately track dial attempts per retry
			errOnErr:      errors.New("connection reset by peer"),
			errOnErrCount: 1, // Force the FIRST conn.Err() call to fail
			presetData:    map[string]interface{}{"user:123": "active"},
			ops: []redisOp{
				{operation: "get_user", command: "GET", args: []interface{}{"user:123"}, wantReply: "active"},
			},
			// Assert the retry actually happened: 1 failed dial + 1 successful retry
			checkDial:     true,
			wantDialCalls: 2,
		},
		{
			name:         "PermanentFailure",
			maxIdle:      0, // Disable pooling to avoid reuse during retries
			errOnDo:      errors.New("Redis cluster unavailable permanently"),
			errOnDoCount: 3, // Force all 3 retries to fail
			ops: []redisOp{
				{operation: "set_user", command: "SET", args: []interface{}{"user:123", "active"}, wantErr: true},
			},
			checkDo:     true,
			wantDoCalls: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := initTestTelemetry()

			data := tc.presetData
			if data == nil {
				data = make(map[string]interface{})
			}
			sharedConn := &mockRedisConn{
				data:          data,
				errOnErr:      tc.errOnErr,
				errOnErrCount: tc.errOnErrCount,
				errOnDo:       tc.errOnDo,
				errOnDoCount:  tc.errOnDoCount,
			}

			pool := &redis.Pool{
				MaxIdle: tc.maxIdle,
				Dial: func() (redis.Conn, error) {
					sharedConn.dialCalls++
					return sharedConn, nil
				},
			}
			defer pool.Close()

			ctx := context.Background()

			for _, op := range tc.ops {
				reply, err := client.smartRedisCall(ctx, pool, op.operation, op.command, op.args...)

				if op.wantErr {
					if err == nil {
						t.Errorf("%s: expected an error but got nil", op.operation)
					}
					continue
				}
				if err != nil {
					t.Fatalf("%s: unexpected error: %v", op.operation, err)
				}
				if op.wantReply != nil {
					s, ok := reply.(string)
					if !ok || s != op.wantReply {
						t.Errorf("%s reply got %v, want %v", op.operation, reply, op.wantReply)
					}
				}
			}

			if tc.checkDial && sharedConn.dialCalls != tc.wantDialCalls {
				t.Errorf("expected %d Dial calls (1 fail, 1 retry), got %d", tc.wantDialCalls, sharedConn.dialCalls)
			}
			if tc.checkDo && sharedConn.doCalls != tc.wantDoCalls {
				t.Errorf("expected %d Do attempts, got %d", tc.wantDoCalls, sharedConn.doCalls)
			}
			if tc.checkClose && sharedConn.closeCalls != tc.wantCloseCalls {
				t.Errorf("expected %d Close calls, got %d", tc.wantCloseCalls, sharedConn.closeCalls)
			}
		})
	}
}