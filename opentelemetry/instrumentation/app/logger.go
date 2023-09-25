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
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

// handerWithTraceContext adds attributes from the trace context
func handerWithTraceContext(handler slog.Handler) *traceContextLogHandler {
	return &traceContextLogHandler{Handler: handler}
}

type traceContextLogHandler struct {
	slog.Handler
}

func (t *traceContextLogHandler) Handle(ctx context.Context, record slog.Record) error {
	if s := trace.SpanContextFromContext(ctx); s.IsValid() {
		if s.HasTraceID() {
			record.AddAttrs(slog.Any("trace_id", s.TraceID()))
		}
		if s.HasSpanID() {
			record.AddAttrs(slog.Any("span_id", s.SpanID()))
		}
		record.AddAttrs(slog.Any("trace_flags", s.TraceFlags()))
	}
	return t.Handler.Handle(ctx, record)
}
