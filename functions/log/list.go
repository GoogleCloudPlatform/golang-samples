// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
