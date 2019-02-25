// Copyright 2019 Google LLC
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

// [START functions_log_retrieve]

// Package log contains logging examples.
package log

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/logging/apiv2"
	"google.golang.org/api/iterator"
	loggingpb "google.golang.org/genproto/googleapis/logging/v2"
)

// logEntries retrieves log entries from projectID and writes them to the
// passed io.Writer.
func logEntries(w io.Writer, projectID string) error {
	ctx := context.Background()
	client, err := logging.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("logging.NewClient: %v", err)
	}

	req := &loggingpb.ListLogEntriesRequest{
		ResourceNames: []string{"projects/" + projectID},
		PageSize:      10,
	}

	fmt.Fprintln(w, "Entries:")
	it := client.ListLogEntries(ctx, req)
	// Wrap in a for loop to get all available log entries.
	resp, err := it.Next()
	if err == iterator.Done {
		return nil
	}
	if err != nil {
		return fmt.Errorf("it.Next: %v", err)
	}
	fmt.Fprintln(w, resp.GetPayload())
	return nil
}

// [END functions_log_retrieve]
