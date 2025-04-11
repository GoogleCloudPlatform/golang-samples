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

package kms

// [START kms_check_state_import_job]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
)

// checkStateImportJob checks the state of an ImportJob in KMS.
func checkStateImportJob(w io.Writer, name string) error {
	// name := "projects/PROJECT_ID/locations/global/keyRings/my-key-ring/importJobs/my-import-job"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Call the API.
	result, err := client.GetImportJob(ctx, &kmspb.GetImportJobRequest{
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("failed to get import job: %w", err)
	}
	fmt.Fprintf(w, "Current state of import job %q: %s\n", result.Name, result.State)
	return nil
}

// [END kms_check_state_import_job]
