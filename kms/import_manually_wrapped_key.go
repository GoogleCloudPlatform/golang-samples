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

// [START kms_import_manually_wrapped_key]
import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/google/tink/go/kwp/subtle"
)

// importManuallyWrappedKey wraps key material and imports it into KMS.
func importManuallyWrappedKey(w io.Writer, importJobName, cryptoKeyName string) error {
	// importJobName := "projects/PROJECT_ID/locations/global/keyRings/my-key-ring/importJobs/my-import-job"
	// cryptoKeyName := "projects/PROJECT_ID/locations/global/keyRings/my-key-ring/cryptoKeys/my-imported-key"

	// Generate a ECDSA keypair, and format the private key as PKCS #8 DER.
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate keypair: %w", err)
	}
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to format private key: %w", err)
	}

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create kms client: %w", err)
	}
	defer client.Close()

	// Generate a temporary 32-byte key for AES-KWP and wrap the key material.
	kwpKey := make([]byte, 32)
	if _, err := rand.Read(kwpKey); err != nil {
		return fmt.Errorf("failed to generate AES-KWP key: %w", err)
	}
	kwp, err := subtle.NewKWP(kwpKey)
	if err != nil {
		return fmt.Errorf("failed to create KWP cipher: %w", err)
	}
	wrappedTarget, err := kwp.Wrap(keyBytes)
	if err != nil {
		return fmt.Errorf("failed to wrap target key with KWP: %w", err)
	}

	// Retrieve the public key from the import job.
	importJob, err := client.GetImportJob(ctx, &kmspb.GetImportJobRequest{
		Name: importJobName,
	})
	if err != nil {
		return fmt.Errorf("failed to retrieve import job: %w", err)
	}
	pubBlock, _ := pem.Decode([]byte(importJob.PublicKey.Pem))
	pubAny, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse import job public key: %w", err)
	}
	pub, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("unexpected public key type %T, want *rsa.PublicKey", pubAny)
	}

	// Wrap the KWP key using the import job key.
	wrappedWrappingKey, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pub, kwpKey, nil)
	if err != nil {
		return fmt.Errorf("failed to wrap KWP key: %w", err)
	}

	// Concatenate the wrapped KWP key and the wrapped target key.
	combined := append(wrappedWrappingKey, wrappedTarget...)

	// Build the request.
	req := &kmspb.ImportCryptoKeyVersionRequest{
		Parent:     cryptoKeyName,
		ImportJob:  importJobName,
		Algorithm:  kmspb.CryptoKeyVersion_EC_SIGN_P256_SHA256,
		WrappedKey: combined,
	}

	// Call the API.
	result, err := client.ImportCryptoKeyVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to import crypto key version: %w", err)
	}
	fmt.Fprintf(w, "Created crypto key version: %s\n", result.Name)
	return nil
}

// [END kms_import_manually_wrapped_key]
