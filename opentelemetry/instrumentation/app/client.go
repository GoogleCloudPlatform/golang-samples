package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func callSingle(ctx context.Context) error {
	res, err := otelhttp.Get(ctx, "http://localhost:8080/single")
	if err != nil {
		return fmt.Errorf("error invoking /single: %w", err)
	}

	err = res.Body.Close()
	if err != nil {
		return fmt.Errorf("error closing response body: %w", err)
	}

	return nil
}
