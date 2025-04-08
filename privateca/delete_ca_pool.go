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

package snippets

// [START privateca_delete_ca_pool]
import (
	"context"
	"fmt"
	"io"

	privateca "cloud.google.com/go/security/privateca/apiv1"
	"cloud.google.com/go/security/privateca/apiv1/privatecapb"
)

// Delete the CA pool as mentioned by the ca_pool_name.
// Before deleting the pool, all CAs in the pool MUST BE deleted.
func deleteCaPool(w io.Writer, projectId string, location string, caPoolId string) error {
	// projectId := "your_project_id"
	// location := "us-central1"	// For a list of locations, see: https://cloud.google.com/certificate-authority-service/docs/locations.
	// caPoolId := "ca-pool-id"		// A unique id/name for the ca pool.

	ctx := context.Background()
	caClient, err := privateca.NewCertificateAuthorityClient(ctx)
	if err != nil {
		return fmt.Errorf("NewCertificateAuthorityClient creation failed: %w", err)
	}
	defer caClient.Close()

	fullCaPoolName := fmt.Sprintf("projects/%s/locations/%s/caPools/%s", projectId, location, caPoolId)

	// See https://pkg.go.dev/cloud.google.com/go/security/privateca/apiv1/privatecapb#DeleteCaPoolRequest.
	req := &privatecapb.DeleteCaPoolRequest{
		Name: fullCaPoolName,
	}

	op, err := caClient.DeleteCaPool(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteCaPool failed: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("DeleteCaPool failed during wait: %w", err)
	}

	fmt.Fprintf(w, "CA Pool deleted")

	return nil
}

// [END privateca_delete_ca_pool]
