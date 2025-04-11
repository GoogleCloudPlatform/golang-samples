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

// [START dlp_deidentify_fpe]
import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// deidentifyFPE deidentifies the input with FPE (Format Preserving Encryption).
// keyFileName is the file name with the KMS wrapped key and cryptoKeyName is the
// full KMS key resource name used to wrap the key. surrogateInfoType is an
// optional identifier needed for reidentification. surrogateInfoType can be any
// value not found in your input.
// Info types can be found with the infoTypes.list method or on https://cloud.google.com/dlp/docs/infotypes-reference
func deidentifyFPE(w io.Writer, projectID, input string, infoTypeNames []string, kmsKeyName, wrappedAESKey, surrogateInfoType string) error {
	// projectID := "my-project-id"
	// input := "My SSN is 123456789"
	// infoTypeNames := []string{"US_SOCIAL_SECURITY_NUMBER"}
	// kmsKeyName := "projects/YOUR_GCLOUD_PROJECT/locations/YOUR_LOCATION/keyRings/YOUR_KEYRING_NAME/cryptoKeys/YOUR_KEY_NAME"
	// wrappedAESKey := "YOUR_ENCRYPTED_AES_256_KEY"
	// surrogateInfoType := "AGE"
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %w", err)
	}
	defer client.Close()
	// Convert the info type strings to a list of InfoTypes.
	var infoTypes []*dlppb.InfoType
	for _, it := range infoTypeNames {
		infoTypes = append(infoTypes, &dlppb.InfoType{Name: it})
	}

	// Specify an encrypted AES-256 key and the name of the Cloud KMS key that encrypted it.
	kmsWrappedCryptoKey, err := base64.StdEncoding.DecodeString(wrappedAESKey)
	if err != nil {
		return err
	}

	// Create a configured request.
	req := &dlppb.DeidentifyContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: infoTypes,
		},
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
				InfoTypeTransformations: &dlppb.InfoTypeTransformations{
					Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
						{
							InfoTypes: []*dlppb.InfoType{}, // Match all info types.
							PrimitiveTransformation: &dlppb.PrimitiveTransformation{
								Transformation: &dlppb.PrimitiveTransformation_CryptoReplaceFfxFpeConfig{
									CryptoReplaceFfxFpeConfig: &dlppb.CryptoReplaceFfxFpeConfig{
										CryptoKey: &dlppb.CryptoKey{
											Source: &dlppb.CryptoKey_KmsWrapped{
												KmsWrapped: &dlppb.KmsWrappedCryptoKey{
													WrappedKey:    kmsWrappedCryptoKey,
													CryptoKeyName: kmsKeyName,
												},
											},
										},
										// Set the alphabet used for the output.
										Alphabet: &dlppb.CryptoReplaceFfxFpeConfig_CommonAlphabet{
											CommonAlphabet: dlppb.CryptoReplaceFfxFpeConfig_ALPHA_NUMERIC,
										},
										// Set the surrogate info type, used for reidentification.
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
		// The item to analyze.
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: input,
			},
		},
	}
	// Send the request.
	r, err := client.DeidentifyContent(ctx, req)
	if err != nil {
		return fmt.Errorf("DeidentifyContent: %w", err)
	}
	// Print the result.
	fmt.Fprint(w, r.GetItem().GetValue())
	return nil
}

// [END dlp_deidentify_fpe]
