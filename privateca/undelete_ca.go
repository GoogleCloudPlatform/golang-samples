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

// [START privateca_undelete_ca]
import (
	"context"
	"fmt"
	"io"

	privateca "cloud.google.com/go/security/privateca/apiv1"
	"cloud.google.com/go/security/privateca/apiv1/privatecapb"
)

// Undelete a Certificate Authority from the specified CA pool.
func unDeleteCa(w io.Writer, projectId string, location string, caPoolId string, caId string) error {
	// projectId := "your_project_id"
	// location := "us-central1"	// For a list of locations, see: https://cloud.google.com/certificate-authority-service/docs/locations.
	// caPoolId := "ca-pool-id"		// The id of the CA pool under which the CA is present.
	// caId := "ca-id"				// The id of the CA to be undeleted.

	ctx := context.Background()
	caClient, err := privateca.NewCertificateAuthorityClient(ctx)
	if err != nil {
		return fmt.Errorf("NewCertificateAuthorityClient creation failed: %w", err)
	}
	defer caClient.Close()

	fullCaName := fmt.Sprintf("projects/%s/locations/%s/caPools/%s/certificateAuthorities/%s",
		projectId, location, caPoolId, caId)

	// Check if the CA is deleted.
	// See https://pkg.go.dev/cloud.google.com/go/security/privateca/apiv1/privatecapb#GetCertificateAuthorityRequest.
	caReq := &privatecapb.GetCertificateAuthorityRequest{Name: fullCaName}
	caResp, err := caClient.GetCertificateAuthority(ctx, caReq)
	if err != nil {
		return fmt.Errorf("GetCertificateAuthority failed: %w", err)
	}

	if caResp.State != privatecapb.CertificateAuthority_DELETED {
		return fmt.Errorf("you can only undelete deleted Certificate Authorities. %s is not deleted", caId)
	}

	// Create the UndeleteCertificateAuthority.
	// See https://pkg.go.dev/cloud.google.com/go/security/privateca/apiv1/privatecapb#UndeleteCertificateAuthorityRequest.
	req := &privatecapb.UndeleteCertificateAuthorityRequest{Name: fullCaName}

	op, err := caClient.UndeleteCertificateAuthority(ctx, req)
	if err != nil {
		return fmt.Errorf("UndeleteCertificateAuthority failed: %w", err)
	}

	if caResp, err = op.Wait(ctx); err != nil {
		return fmt.Errorf("UndeleteCertificateAuthority failed during wait: %w", err)
	}

	if caResp.State == privatecapb.CertificateAuthority_DELETED {
		return fmt.Errorf("unable to undelete Certificate Authority. Current state: %s", caResp.State.String())
	}

	fmt.Fprintf(w, "Successfully undeleted Certificate Authority: %s.", caId)
	return nil
}

// [END privateca_undelete_ca]
