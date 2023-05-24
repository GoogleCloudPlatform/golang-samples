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

// [START privateca_create_ca_pool]
import (
	"context"
	"fmt"
	"io"

	privateca "cloud.google.com/go/security/privateca/apiv1"
	"cloud.google.com/go/security/privateca/apiv1/privatecapb"
)

// Create a Certificate Authority pool. All certificates created under this CA pool will
// follow the same issuance policy, IAM policies, etc.
func createCaPool(w io.Writer, projectId string, location string, caPoolId string) error {
	// projectId := "your_project_id"
	// location := "us-central1"	// For a list of locations, see: https://cloud.google.com/certificate-authority-service/docs/locations.
	// caPoolId := "ca-pool-id"		// A unique id/name for the ca pool.

	ctx := context.Background()
	caClient, err := privateca.NewCertificateAuthorityClient(ctx)
	if err != nil {
		return fmt.Errorf("NewCertificateAuthorityClient creation failed: %w", err)
	}
	defer caClient.Close()

	caPool := &privatecapb.CaPool{
		// Set the tier (see: https://cloud.google.com/certificate-authority-service/docs/tiers).
		Tier: privatecapb.CaPool_ENTERPRISE,
	}

	locationPath := fmt.Sprintf("projects/%s/locations/%s", projectId, location)

	// See https://pkg.go.dev/cloud.google.com/go/security/privateca/apiv1/privatecapb#CreateCaPoolRequest.
	req := &privatecapb.CreateCaPoolRequest{
		Parent:   locationPath,
		CaPoolId: caPoolId,
		CaPool:   caPool,
	}

	op, err := caClient.CreateCaPool(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateCaPool failed: %w", err)
	}

	if _, err = op.Wait(ctx); err != nil {
		return fmt.Errorf("CreateCaPool failed during wait: %w", err)
	}

	fmt.Fprintf(w, "CA Pool created")

	return nil
}

// [END privateca_create_ca_pool]
