// Copyright 2019 Google LLC
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

// [START kms_sign_asymmetric]
import (
	"context"
	"crypto/sha256"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// signAsymmetric will sign a plaintext message using a saved asymmetric private key.
func signAsymmetric(name string, message []byte) ([]byte, error) {
	// name: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID/cryptoKeyVersions/1"
	// Note: some key algorithms will require a different hash function.
	// For example, EC_SIGN_P384_SHA384 requires SHA-384.
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cloudkms.NewKeyManagementClient: %v", err)
	}
	// Find the digest of the message.
	digest := sha256.New()
	digest.Write(message)
	// Build the signing request.
	req := &kmspb.AsymmetricSignRequest{
		Name: name,
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: digest.Sum(nil),
			},
		},
	}
	// Call the API.
	response, err := client.AsymmetricSign(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("AsymmetricSign: %v", err)
	}
	// Return the signature bytes.
	return response.Signature, nil
}

// [END kms_sign_asymmetric]
