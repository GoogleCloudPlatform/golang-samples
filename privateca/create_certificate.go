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

// [START privateca_create_certificate]
import (
	"context"
	"fmt"
	"io"

	privateca "cloud.google.com/go/security/privateca/apiv1"
	"cloud.google.com/go/security/privateca/apiv1/privatecapb"
	"google.golang.org/protobuf/types/known/durationpb"
)

// Create a Certificate which is issued by the Certificate Authority present in the CA Pool.
// The key used to sign the certificate is created by the Cloud KMS.
func createCertificate(
	w io.Writer,
	projectId string,
	location string,
	caPoolId string,
	caId string,
	certId string,
	commonName string,
	domainName string,
	certDuration int64,
	publicKeyBytes []byte) error {
	// projectId := "your_project_id"
	// location := "us-central1"		// For a list of locations, see: https://cloud.google.com/certificate-authority-service/docs/locations.
	// caPoolId := "ca-pool-id"			// The CA Pool id in which the certificate authority exists.
	// caId := "ca-id"					// The name of the certificate authority which issues the certificate.
	// certId := "certificate"			// A unique name for the certificate.
	// commonName := "cert-name"		// A common name for the certificate.
	// domainName := "cert.example.com"	// Fully qualified domain name for the certificate.
	// certDuration := int64(31536000)	// The validity of the certificate in seconds.
	// publicKeyBytes 					// The public key used in signing the certificates.

	ctx := context.Background()
	caClient, err := privateca.NewCertificateAuthorityClient(ctx)
	if err != nil {
		return fmt.Errorf("NewCertificateAuthorityClient creation failed: %w", err)
	}
	defer caClient.Close()

	// Set the Public Key and its format.
	publicKey := &privatecapb.PublicKey{
		Key:    publicKeyBytes,
		Format: privatecapb.PublicKey_PEM,
	}

	// Set Certificate subject config.
	subjectConfig := &privatecapb.CertificateConfig_SubjectConfig{
		Subject: &privatecapb.Subject{
			CommonName: commonName,
		},
		SubjectAltName: &privatecapb.SubjectAltNames{
			DnsNames: []string{domainName},
		},
	}

	// Set the X.509 fields required for the certificate.
	x509Parameters := &privatecapb.X509Parameters{
		KeyUsage: &privatecapb.KeyUsage{
			BaseKeyUsage: &privatecapb.KeyUsage_KeyUsageOptions{
				DigitalSignature: true,
				KeyEncipherment:  true,
			},
			ExtendedKeyUsage: &privatecapb.KeyUsage_ExtendedKeyUsageOptions{
				ServerAuth: true,
				ClientAuth: true,
			},
		},
	}

	// Set certificate settings.
	cert := &privatecapb.Certificate{
		CertificateConfig: &privatecapb.Certificate_Config{
			Config: &privatecapb.CertificateConfig{
				PublicKey:     publicKey,
				SubjectConfig: subjectConfig,
				X509Config:    x509Parameters,
			},
		},
		Lifetime: &durationpb.Duration{
			Seconds: certDuration,
		},
	}

	fullCaPoolName := fmt.Sprintf("projects/%s/locations/%s/caPools/%s", projectId, location, caPoolId)

	// Create the CreateCertificateRequest.
	// See https://pkg.go.dev/cloud.google.com/go/security/privateca/apiv1/privatecapb#CreateCertificateRequest.
	req := &privatecapb.CreateCertificateRequest{
		Parent:                        fullCaPoolName,
		CertificateId:                 certId,
		Certificate:                   cert,
		IssuingCertificateAuthorityId: caId,
	}

	_, err = caClient.CreateCertificate(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateCertificate failed: %w", err)
	}

	fmt.Fprintf(w, "Certificate %s created", certId)

	return nil
}

// [END privateca_create_certificate]
