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
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHandlerSeverity(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		expectSeverity string
		logFunc        func(*slog.Logger)
	}{
		{
			expectSeverity: "DEBUG",
			logFunc:        func(l *slog.Logger) { l.DebugContext(ctx, "debug") },
		},
		{
			expectSeverity: "INFO",
			logFunc:        func(l *slog.Logger) { l.InfoContext(ctx, "info") },
		},
		{
			expectSeverity: "WARNING",
			logFunc:        func(l *slog.Logger) { l.WarnContext(ctx, "warn") },
		},
		{
			expectSeverity: "ERROR",
			logFunc:        func(l *slog.Logger) { l.ErrorContext(ctx, "error") },
		},
	}

	for _, tc := range tests {
		t.Run(tc.expectSeverity, func(t *testing.T) {
			buf := &bytes.Buffer{}
			jsonHandler := slog.NewJSONHandler(
				buf,
				&slog.HandlerOptions{ReplaceAttr: replacer, Level: slog.LevelDebug},
			)
			logger := slog.New(jsonHandler)
			tc.logFunc(logger)

			line := &expectedLogFormat{}
			require.NoError(t, json.Unmarshal(buf.Bytes(), line))
			require.Equal(t, tc.expectSeverity, line.Severity)
		})
	}
}

func TestHandlerTimestamp(t *testing.T) {
	ctx := context.Background()

	buf := &bytes.Buffer{}
	jsonHandler := slog.NewJSONHandler(
		buf,
		&slog.HandlerOptions{ReplaceAttr: replacer, Level: slog.LevelDebug},
	)
	logger := slog.New(jsonHandler)
	logger.InfoContext(ctx, "foo")

	line := &expectedLogFormat{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), line))
	require.NotEmpty(t, line.Timestamp)

	_, err := time.Parse(time.RFC3339Nano, line.Timestamp)
	require.NoErrorf(t, err, "could not parse timestamp as RFC3339 with nanos")
}
