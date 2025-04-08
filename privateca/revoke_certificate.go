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

// [START privateca_revoke_certificate]
import (
	"context"
	"fmt"
	"io"

	privateca "cloud.google.com/go/security/privateca/apiv1"
	"cloud.google.com/go/security/privateca/apiv1/privatecapb"
)

// Revoke an issued certificate. Once revoked, the certificate will become invalid
// and will expire post its lifetime.
func revokeCertificate(
	w io.Writer,
	projectId string,
	location string,
	caPoolId string,
	certId string) error {
	// projectId := "your_project_id"
	// location := "us-central1"		// For a list of locations, see: https://cloud.google.com/certificate-authority-service/docs/locations.
	// caPoolId := "ca-pool-id"			// The CA Pool id in which the certificate exists.
	// certId := "certificate"			// A unique name for the certificate.

	ctx := context.Background()
	caClient, err := privateca.NewCertificateAuthorityClient(ctx)
	if err != nil {
		return fmt.Errorf("NewCertificateAuthorityClient creation failed: %w", err)
	}
	defer caClient.Close()

	fullCertName := fmt.Sprintf("projects/%s/locations/%s/caPools/%s/certificates/%s", projectId, location,
		caPoolId, certId)

	// Create the RevokeCertificateRequest and specify the appropriate revocation reason.
	// See https://pkg.go.dev/cloud.google.com/go/security/privateca/apiv1/privatecapb#RevokeCertificateRequest.
	req := &privatecapb.RevokeCertificateRequest{
		Name:   fullCertName,
		Reason: privatecapb.RevocationReason_PRIVILEGE_WITHDRAWN,
	}

	_, err = caClient.RevokeCertificate(ctx, req)
	if err != nil {
		return fmt.Errorf("RevokeCertificate failed: %w", err)
	}

	fmt.Fprintf(w, "Certificate %s revoked", certId)

	return nil
}

// [END privateca_revoke_certificate]
