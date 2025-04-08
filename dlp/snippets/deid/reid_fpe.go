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

package deid

// [START dlp_reidentify_fpe]
import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// reidentifyFPE reidentifies the input with FPE (Format Preserving Encryption).
// keyFileName is the file name with the KMS wrapped key and cryptoKeyName is the
// full KMS key resource name used to wrap the key. surrogateInfoType is an
// the identifier used during deidentification.
// Info types can be found with the infoTypes.list method or on https://cloud.google.com/dlp/docs/infotypes-reference
func reidentifyFPE(w io.Writer, projectID, input, kmsKeyName, wrappedAesKey, surrogateInfoType string) error {
	// projectID := "my-project-id"
	// input := "My SSN is 123456789"
	// keyFileName := "projects/YOUR_GCLOUD_PROJECT/locations/YOUR_LOCATION/keyRings/YOUR_KEYRING_NAME/cryptoKeys/YOUR_KEY_NAME"
	// cryptoKeyName := "YOUR_ENCRYPTED_AES_256_KEY"
	// surrogateInfoType := "AGE"
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %w", err)
	}
	defer client.Close()

	// Specify an encrypted AES-256 key and the name of the Cloud KMS key that encrypted it.
	kmsWrappedCryptoKey, err := base64.StdEncoding.DecodeString(wrappedAesKey)
	if err != nil {
		fmt.Fprintf(w, "error %v", err)
		return err
	}

	// Specify the crypto key configuration that will used for encryption.
	cryptoKey := &dlppb.CryptoKey{
		Source: &dlppb.CryptoKey_KmsWrapped{
			KmsWrapped: &dlppb.KmsWrappedCryptoKey{
				WrappedKey:    kmsWrappedCryptoKey,
				CryptoKeyName: kmsKeyName,
			},
		},
	}

	// Create a configured request.
	req := &dlppb.ReidentifyContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		ReidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
				InfoTypeTransformations: &dlppb.InfoTypeTransformations{
					Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
						{
							InfoTypes: []*dlppb.InfoType{},
							PrimitiveTransformation: &dlppb.PrimitiveTransformation{
								Transformation: &dlppb.PrimitiveTransformation_CryptoReplaceFfxFpeConfig{
									CryptoReplaceFfxFpeConfig: &dlppb.CryptoReplaceFfxFpeConfig{
										CryptoKey: cryptoKey,
										// Set the alphabet used for the encrypted fields.
										Alphabet: &dlppb.CryptoReplaceFfxFpeConfig_CommonAlphabet{
											CommonAlphabet: dlppb.CryptoReplaceFfxFpeConfig_ALPHA_NUMERIC,
										},
										// Set the surrogate info type used during deidentification.
										SurrogateInfoType: &dlppb.InfoType{
											Name: surrogateInfoType,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		// The InspectConfig must identify the surrogate info type to reidentify.
		InspectConfig: &dlppb.InspectConfig{
			CustomInfoTypes: []*dlppb.CustomInfoType{
				{
					InfoType: &dlppb.InfoType{
						Name: surrogateInfoType,
					},
					Type: &dlppb.CustomInfoType_SurrogateType_{
						SurrogateType: &dlppb.CustomInfoType_SurrogateType{},
					},
				},
			},
		},
		// The item to analyze.
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: input,
			},
		},
	}
	// Send the request.
	r, err := client.ReidentifyContent(ctx, req)
	if err != nil {
		return fmt.Errorf("ReidentifyContent: %w", err)
	}
	// Print the result.
	fmt.Fprint(w, r.GetItem().GetValue())
	return nil
}

// [END dlp_reidentify_fpe]
