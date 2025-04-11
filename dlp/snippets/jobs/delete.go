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

package jobs

// [START dlp_delete_job]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// deleteJob deletes the job with the given name.
func deleteJob(w io.Writer, jobName string) error {
	// jobName := "job-example"
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %w", err)
	}
	defer client.Close()
	req := &dlppb.DeleteDlpJobRequest{
		Name: jobName,
	}
	if err = client.DeleteDlpJob(ctx, req); err != nil {
		return fmt.Errorf("DeleteDlpJob: %w", err)
	}
	fmt.Fprintf(w, "Successfully deleted job")
	return nil
}

// [END dlp_delete_job]
