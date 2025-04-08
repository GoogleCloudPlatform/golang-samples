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

// [START kms_create_import_job]
import (
	"context"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
)

// createImportJob creates a new job for importing keys into KMS.
func createImportJob(w io.Writer, parent, id string) error {
	// parent := "projects/PROJECT_ID/locations/global/keyRings/my-key-ring"
	// id := "my-import-job"

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &kmspb.CreateImportJobRequest{
		Parent:      parent,
		ImportJobId: id,
		ImportJob: &kmspb.ImportJob{
			// See allowed values and their descriptions at
			// https://cloud.google.com/kms/docs/algorithms#protection_levels
			ProtectionLevel: kmspb.ProtectionLevel_HSM,
			// See allowed values and their descriptions at
			// https://cloud.google.com/kms/docs/key-wrapping#import_methods
			ImportMethod: kmspb.ImportJob_RSA_OAEP_3072_SHA1_AES_256,
		},
	}

	// Call the API.
	result, err := client.CreateImportJob(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create import job: %w", err)
	}
	fmt.Fprintf(w, "Created import job: %s\n", result.Name)
	return nil
}

// [END kms_create_import_job]
