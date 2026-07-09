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
	data           map[string]interface{}
	errOnDo        error
	errOnDoCount   int
	errOnErr       error
	errOnErrCount  int
	errOnDial      error
	errOnDialCount int
	closeCalls     int
	doCalls        int
	dialCalls      int
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

// initTestTelemetry safely binds package-level OTel globals to hermetic No-op implementations
func initTestTelemetry() {
	tp := trace.NewNoopTracerProvider()
	tracer = tp.Tracer("noop")

	meter := otel.GetMeterProvider().Meter("noop")
	rttHist, _ = meter.Float64Histogram("redis_client_rtt", metric.WithUnit("ms"))
	clientBlockHist, _ = meter.Float64Histogram("redis_client_blocking_latency", metric.WithUnit("ms"))
	appBlockHist, _ = meter.Float64Histogram("redis_application_blocking_latency", metric.WithUnit("ms"))
	retryCounter, _ = meter.Int64Counter("redis_retry_count")
	connErrorCounter, _ = meter.Int64Counter("redis_connectivity_error_count")

	// Disable sleeping during unit tests to make the test suite lightning fast
	sleep = func(d time.Duration) {}
}

func TestSmartRedisCall_Success(t *testing.T) {
	initTestTelemetry()

	sharedConn := &mockRedisConn{
		data: make(map[string]interface{}),
	}

	pool := &redis.Pool{
		MaxIdle: 1, // Redigo will pool and reuse this connection
		Dial: func() (redis.Conn, error) {
			sharedConn.dialCalls++
			return sharedConn, nil
		},
	}
	defer pool.Close()

	ctx := context.Background()

	// 1. Test SET
	setReply, err := smartRedisCall(ctx, pool, "set_user", "SET", "user:123", "active")
	if err != nil {
		t.Fatalf("smartRedisCall SET failed: %v", err)
	}
	if s, ok := setReply.(string); !ok || s != "OK" {
		t.Errorf("SET reply got %v, want 'OK'", setReply)
	}

	// 2. Test GET
	getReply, err := smartRedisCall(ctx, pool, "get_user", "GET", "user:123")
	if err != nil {
		t.Fatalf("smartRedisCall GET failed: %v", err)
	}
	if s, ok := getReply.(string); !ok || s != "active" {
		t.Errorf("GET reply got %v, want 'active'", getReply)
	}

	// Adjusted for Redigo pooling: 1 Dial since the GET reuses the SET's idle connection
	if sharedConn.dialCalls != 1 {
		t.Errorf("expected 1 Dial call (connection reuse), got %d", sharedConn.dialCalls)
	}
	// Pool intercept means Close is deferred until pool.Close()
	if sharedConn.closeCalls != 0 {
		t.Errorf("expected 0 Close calls on the underlying mock during runtime, got %d", sharedConn.closeCalls)
	}
}

func TestSmartRedisCall_RetryOnConnError(t *testing.T) {
	initTestTelemetry()

	sharedConn := &mockRedisConn{
		data: make(map[string]interface{}),
		// Force the FIRST conn.Err() call to fail, subsequent calls to succeed
		errOnErr:      errors.New("connection reset by peer"),
		errOnErrCount: 1,
	}

	pool := &redis.Pool{
		MaxIdle: 0, // Disable pooling to accurately track dial/close attempts per retry
		Dial: func() (redis.Conn, error) {
			sharedConn.dialCalls++
			return sharedConn, nil
		},
	}
	defer pool.Close()

	ctx := context.Background()

	// Configure expected SET response
	sharedConn.data["user:123"] = "active"

	// Call GET which should trigger the retry logic upon seeing the conn.Err()
	getReply, err := smartRedisCall(ctx, pool, "get_user", "GET", "user:123")
	if err != nil {
		t.Fatalf("smartRedisCall GET failed on retry: %v", err)
	}
	if s, ok := getReply.(string); !ok || s != "active" {
		t.Errorf("GET reply got %v, want 'active'", getReply)
	}

	t.Logf("RetryOnConnError: dialCalls=%d, closeCalls=%d", sharedConn.dialCalls, sharedConn.closeCalls)
}

func TestSmartRedisCall_PermanentFailure(t *testing.T) {
	initTestTelemetry()

	sharedConn := &mockRedisConn{
		errOnDo:      errors.New("Redis cluster unavailable permanently"),
		errOnDoCount: 3, // Force all 3 retries to fail
	}

	pool := &redis.Pool{
		MaxIdle: 0, // Disable pooling to avoid reuse during retries
		Dial: func() (redis.Conn, error) {
			sharedConn.dialCalls++
			return sharedConn, nil
		},
	}
	defer pool.Close()

	ctx := context.Background()

	reply, err := smartRedisCall(ctx, pool, "set_user", "SET", "user:123", "active")

	// Print precise state to identify the root cause of the nil error
	t.Logf("PermanentFailure Debug: reply=%v, err=%v, dialCalls=%d, doCalls=%d, remainingErrs=%d",
		reply, err, sharedConn.dialCalls, sharedConn.doCalls, sharedConn.errOnDoCount)

	if err == nil {
		t.Error("expected smartRedisCall to return a permanent failure error but it returned nil")
	}

	if sharedConn.doCalls != 3 {
		t.Errorf("expected 3 Do attempts before permanent failure, got %d", sharedConn.doCalls)
	}
}
