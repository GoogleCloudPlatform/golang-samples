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

// [START privateca_enable_ca]
import (
	"context"
	"fmt"
	"io"

	privateca "cloud.google.com/go/security/privateca/apiv1"
	"cloud.google.com/go/security/privateca/apiv1/privatecapb"
)

// Enable the Certificate Authority present in the given ca pool.
// CA cannot be enabled if it has been already deleted.
func enableCa(w io.Writer, projectId string, location string, caPoolId string, caId string) error {
	// projectId := "your_project_id"
	// location := "us-central1"	// For a list of locations, see: https://cloud.google.com/certificate-authority-service/docs/locations.
	// caPoolId := "ca-pool-id"		// The id of the CA pool under which the CA is present.
	// caId := "ca-id"				// The id of the CA to be enabled.

	ctx := context.Background()
	caClient, err := privateca.NewCertificateAuthorityClient(ctx)
	if err != nil {
		return fmt.Errorf("NewCertificateAuthorityClient creation failed: %w", err)
	}
	defer caClient.Close()

	fullCaName := fmt.Sprintf("projects/%s/locations/%s/caPools/%s/certificateAuthorities/%s",
		projectId, location, caPoolId, caId)

	// Create the EnableCertificateAuthorityRequest.
	// See https://pkg.go.dev/cloud.google.com/go/security/privateca/apiv1/privatecapb#EnableCertificateAuthorityRequest.
	req := &privatecapb.EnableCertificateAuthorityRequest{Name: fullCaName}

	op, err := caClient.EnableCertificateAuthority(ctx, req)
	if err != nil {
		return fmt.Errorf("EnableCertificateAuthority failed: %w", err)
	}

	var caResp *privatecapb.CertificateAuthority
	if caResp, err = op.Wait(ctx); err != nil {
		return fmt.Errorf("EnableCertificateAuthority failed during wait: %w", err)
	}

	if caResp.State != privatecapb.CertificateAuthority_ENABLED {
		return fmt.Errorf("unable to enable Certificate Authority. Current state: %s", caResp.State.String())
	}

	fmt.Fprintf(w, "Successfully enabled Certificate Authority: %s.", caId)
	return nil
}

// [END privateca_enable_ca]
