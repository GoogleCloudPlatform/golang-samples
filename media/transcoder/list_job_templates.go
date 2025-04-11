// Copyright 2020 Google LLC
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

package transcoder

// [START transcoder_list_job_templates]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/api/iterator"

	transcoder "cloud.google.com/go/video/transcoder/apiv1"
	"cloud.google.com/go/video/transcoder/apiv1/transcoderpb"
)

// listJobTemplates gets all previously-created job templates for a given
// location. See
// https://cloud.google.com/transcoder/docs/how-to/job-templates#list_job_template
// for more information.
func listJobTemplates(w io.Writer, projectID string, location string) error {
	// projectID := "my-project-id"
	// location := "us-central1"
	ctx := context.Background()
	client, err := transcoder.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &transcoderpb.ListJobTemplatesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}

	it := client.ListJobTemplates(ctx, req)
	fmt.Fprintln(w, "Job templates:")
	for {
		response, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListJobTemplates: %w", err)
		}
		fmt.Fprintln(w, response.GetName())
	}

	return nil
}

// [END transcoder_list_job_templates]
