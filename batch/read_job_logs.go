// Copyright 2022 Google LLC
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

package snippets

// [START batch_job_logs]
import (
	"context"
	"fmt"
	"io"

	batch "cloud.google.com/go/batch/apiv1"
	"cloud.google.com/go/batch/apiv1/batchpb"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"google.golang.org/api/iterator"
)

// Retrieve the logs written by the given job to Cloud Logging
func printJobLogs(w io.Writer, projectID string, job *batchpb.Job) error {
	// projectID := "your_project_id"

	ctx := context.Background()
	batchClient, err := batch.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer batchClient.Close()

	adminClient, err := logadmin.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("Failed to create logadmin client: %w", err)
	}
	defer adminClient.Close()

	const name = "batch_task_logs"

	iter := adminClient.Entries(ctx,
		// Only get entries from the "batch_task_logs" log for the job with the given UID
		logadmin.Filter(fmt.Sprintf(`logName = "projects/%s/logs/%s" AND labels.job_uid=%s`, projectID, name, job.Uid)),
	)

	var entries []*logging.Entry

	for {
		logEntry, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("unable to fetch log entry: %w", err)
		}
		entries = append(entries, logEntry)
		fmt.Fprintf(w, "%s\n", logEntry.Payload)
	}

	fmt.Fprintf(w, "Successfully fetched %d log entries\n", len(entries))

	return nil
}

// [END batch_job_logs]
