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

// handerWithSpanContext adds attributes from the span context
// [START opentelemetry_instrumentation_spancontext_logger]
func handerWithSpanContext(handler slog.Handler) *spanContextLogHandler {
	return &spanContextLogHandler{Handler: handler}
}

// spanContextLogHandler is an slog.Handler which adds attributes from the
// span context.
type spanContextLogHandler struct {
	slog.Handler
}

// Handle overrides slog.Handler's Handle method. This adds attributes from the
// span context to the slog.Record.
func (t *spanContextLogHandler) Handle(ctx context.Context, record slog.Record) error {
	// Get the SpanContext from the golang Context.
	if s := trace.SpanContextFromContext(ctx); s.IsValid() {
		// Add the trace_id attribute from the SpanContext.
		if s.HasTraceID() {
			record.AddAttrs(
				slog.Any("trace_id", s.TraceID()),
			)
		}

		// Add the span_id attribute from the SpanContext.
		if s.HasSpanID() {
			record.AddAttrs(
				slog.Any("span_id", s.SpanID()),
			)
		}

		// Add the trace_flags attribute from the SpanContext.
		// This includes whether or not the trace is sampled.
		record.AddAttrs(
			slog.Any("trace_flags", s.TraceFlags()),
		)
	}
	return t.Handler.Handle(ctx, record)
}

// [END opentelemetry_instrumentation_spancontext_logger]
