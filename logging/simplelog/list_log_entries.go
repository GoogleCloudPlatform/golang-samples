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

// Sample simplelog writes some entries, lists them, then deletes the log.
package main

// [START logging_list_log_entries]
import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"google.golang.org/api/iterator"
)

func getEntries(projectID string) ([]*logging.Entry, error) {
	ctx := context.Background()
	adminClient, err := logadmin.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create logadmin client: %v", err)
	}
	defer adminClient.Close()

	var entries []*logging.Entry
	const name = "log-example"
	lastHour := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)

	iter := adminClient.Entries(ctx,
		// Only get entries from the "log-example" log within the last hour.
		logadmin.Filter(fmt.Sprintf(`logName = "projects/%s/logs/%s" AND timestamp > "%s"`, projectID, name, lastHour)),
		// Get most recent entries first.
		logadmin.NewestFirst(),
	)

	// Fetch the most recent 20 entries.
	for len(entries) < 20 {
		entry, err := iter.Next()
		if err == iterator.Done {
			return entries, nil
		}
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// [END logging_list_log_entries]
