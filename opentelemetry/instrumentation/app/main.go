package main

import (
	"context"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()
	if err := withOpenTelemetry(ctx, runServer); err != nil {
		slog.ErrorContext(ctx, "server exited with error", slog.Any("error", err))
		os.Exit(1)
	}
}
