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

// [START privateca_create_ca]
import (
	"context"
	"fmt"
	"io"

	privateca "cloud.google.com/go/security/privateca/apiv1"
	"cloud.google.com/go/security/privateca/apiv1/privatecapb"
	"google.golang.org/protobuf/types/known/durationpb"
)

// Create Certificate Authority which is the root CA in the given CA Pool. This CA will be
// responsible for signing certificates within this pool.
func createCa(
	w io.Writer,
	projectId string,
	location string,
	caPoolId string,
	caId string,
	caCommonName string,
	org string,
	caDuration int64) error {
	// projectId := "your_project_id"
	// location := "us-central1"		// For a list of locations, see: https://cloud.google.com/certificate-authority-service/docs/locations.
	// caPoolId := "ca-pool-id"			// The CA Pool id under which the CA should be created.
	// caId := "ca-id"					// A unique id/name for the ca.
	// caCommonName := "ca-name"		// A common name for your certificate authority.
	// org := "ca-org"					// The name of your company for your certificate authority.
	// ca_duration := int64(31536000)	// The validity of the certificate authority in seconds.

	ctx := context.Background()
	caClient, err := privateca.NewCertificateAuthorityClient(ctx)
	if err != nil {
		return fmt.Errorf("NewCertificateAuthorityClient creation failed: %w", err)
	}
	defer caClient.Close()

	// Set the types of Algorithm used to create a cloud KMS key.
	keySpec := &privatecapb.CertificateAuthority_KeyVersionSpec{
		KeyVersion: &privatecapb.CertificateAuthority_KeyVersionSpec_Algorithm{
			Algorithm: privatecapb.CertificateAuthority_RSA_PKCS1_2048_SHA256,
		},
	}

	// Set CA subject config.
	subjectConfig := &privatecapb.CertificateConfig_SubjectConfig{
		Subject: &privatecapb.Subject{
			CommonName:   caCommonName,
			Organization: org,
		},
	}

	// Set the key usage options for X.509 fields.
	isCa := true
	x509Parameters := &privatecapb.X509Parameters{
		KeyUsage: &privatecapb.KeyUsage{
			BaseKeyUsage: &privatecapb.KeyUsage_KeyUsageOptions{
				CrlSign:  true,
				CertSign: true,
			},
		},
		CaOptions: &privatecapb.X509Parameters_CaOptions{
			IsCa: &isCa,
		},
	}

	// Set certificate authority settings.
	// Type: SELF_SIGNED denotes that this CA is a root CA.
	ca := &privatecapb.CertificateAuthority{
		Type:    privatecapb.CertificateAuthority_SELF_SIGNED,
		KeySpec: keySpec,
		Config: &privatecapb.CertificateConfig{
			SubjectConfig: subjectConfig,
			X509Config:    x509Parameters,
		},
		Lifetime: &durationpb.Duration{
			Seconds: caDuration,
		},
	}

	fullCaPoolName := fmt.Sprintf("projects/%s/locations/%s/caPools/%s", projectId, location, caPoolId)

	// Create the CreateCertificateAuthorityRequest.
	// See https://pkg.go.dev/cloud.google.com/go/security/privateca/apiv1/privatecapb#CreateCertificateAuthorityRequest.
	req := &privatecapb.CreateCertificateAuthorityRequest{
		Parent:                 fullCaPoolName,
		CertificateAuthorityId: caId,
		CertificateAuthority:   ca,
	}

	op, err := caClient.CreateCertificateAuthority(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateCertificateAuthority failed: %w", err)
	}

	if _, err = op.Wait(ctx); err != nil {
		return fmt.Errorf("CreateCertificateAuthority failed during wait: %w", err)
	}

	fmt.Fprintf(w, "CA %s created", caId)

	return nil
}

// [END privateca_create_ca]
