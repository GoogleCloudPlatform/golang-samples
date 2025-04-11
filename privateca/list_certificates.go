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

// [START privateca_list_certificate]
import (
	"context"
	"fmt"
	"io"

	privateca "cloud.google.com/go/security/privateca/apiv1"
	"cloud.google.com/go/security/privateca/apiv1/privatecapb"
	"google.golang.org/api/iterator"
)

// List Certificates present in the given CA pool.
func listCertificates(
	w io.Writer,
	projectId string,
	location string,
	caPoolId string) error {
	// projectId := "your_project_id"
	// location := "us-central1"		// For a list of locations, see: https://cloud.google.com/certificate-authority-service/docs/locations.
	// caPoolId := "ca-pool-id"			// The CA Pool id in which the certificate exists.

	ctx := context.Background()
	caClient, err := privateca.NewCertificateAuthorityClient(ctx)
	if err != nil {
		return fmt.Errorf("NewCertificateAuthorityClient creation failed: %w", err)
	}
	defer caClient.Close()

	fullCaName := fmt.Sprintf("projects/%s/locations/%s/caPools/%s", projectId, location, caPoolId)

	// Create the ListCertificatesRequest.
	// See https://pkg.go.dev/cloud.google.com/go/security/privateca/apiv1/privatecapb#ListCertificatesRequest.
	req := &privatecapb.ListCertificatesRequest{Parent: fullCaName}

	it := caClient.ListCertificates(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("unable to get the list of cerficates: %w", err)
		}

		fmt.Fprintf(w, " - %s (common name: %s)", resp.Name,
			resp.CertificateDescription.SubjectDescription.Subject.CommonName)
	}

	return nil
}

// [END privateca_list_certificate]
