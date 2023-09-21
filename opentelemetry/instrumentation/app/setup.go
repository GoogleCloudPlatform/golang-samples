package main

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

func withOpenTelemetry(ctx context.Context, fn func() error) error {
	otel.SetTextMapPropagator(autoprop.NewTextMapPropagator())

	texporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return err
	}
	tp := trace.NewTracerProvider(trace.WithBatcher(texporter))
	defer tp.Shutdown(context.TODO())
	otel.SetTracerProvider(tp)

	// TODO(https://github.com/open-telemetry/opentelemetry-go-contrib/issues/4131) - use env-based configuration instead of hardcoding
	mexporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return err
	}
	mp := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(mexporter)))
	defer mp.Shutdown(ctx)
	otel.SetMeterProvider(mp)

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	return fn()
}
