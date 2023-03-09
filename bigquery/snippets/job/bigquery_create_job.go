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

package job

// [START bigquery_create_job]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"

	"github.com/google/uuid"
)

// createJob demonstrates running an arbitrary SQL statement as a query job.
func createJob(projectID, sql string) error {
	// sql := "SELECT country_name from `bigquery-public-data.utility_us.country_code_iso`:"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// Demonstrate adding a label to the job.
	q := client.Query(sql)
	q.Labels = map[string]string{"example-label": "example-value"}

	// The library will create job IDs for you automatically, but this can be overridden by
	// setting the Job ID explicitly.  Job IDs are unique within a project and cannot be
	// reused.
	q.JobID = fmt.Sprintf("my_job_prefix_%s", uuid.New().String())

	// Start job execution.
	job, err := q.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	return nil
}

// [END bigquery_create_job]
