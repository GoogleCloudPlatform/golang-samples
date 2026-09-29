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

package observability

// [START go_observability_logging]
import (
	"context"
	"fmt"
	"log/slog"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"google.golang.org/api/option"
)

// enableLogging demonstrates how to configure structured actionable error
// logging for Google Cloud Go client libraries. Set GOOGLE_SDK_GO_LOGGING=true
// in the environment before starting the application to enable error logging.
func enableLogging() error {
	ctx := context.Background()

	// Configure slog to output JSON to stdout.
	// Use slog.LevelWarn to capture terminal client request failures, or
	// slog.LevelDebug to also capture per-attempt network failures across retries.
	opts := &slog.HandlerOptions{Level: slog.LevelWarn}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	client, err := secretmanager.NewClient(ctx, option.WithLogger(logger))
	if err != nil {
		return fmt.Errorf("secretmanager.NewClient: %w", err)
	}
	defer client.Close()

	// Use the client to make requests...

	return nil
}

// [END go_observability_logging]
