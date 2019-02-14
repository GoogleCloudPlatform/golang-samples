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

package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	dlp "cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// [START dlp_deidentify_masking]

// mask deidentifies the input by masking all provided info types with maskingCharacter
// and prints the result to w.
func mask(w io.Writer, client *dlp.Client, project, input string, infoTypes []string, maskingCharacter string, numberToMask int32) {
	// Convert the info type strings to a list of InfoTypes.
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}
	// Create a configured request.
	req := &dlppb.DeidentifyContentRequest{
		Parent: "projects/" + project,
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: i,
		},
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
				InfoTypeTransformations: &dlppb.InfoTypeTransformations{
					Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
						{
							InfoTypes: []*dlppb.InfoType{}, // Match all info types.
							PrimitiveTransformation: &dlppb.PrimitiveTransformation{
								Transformation: &dlppb.PrimitiveTransformation_CharacterMaskConfig{
									CharacterMaskConfig: &dlppb.CharacterMaskConfig{
										MaskingCharacter: maskingCharacter,
										NumberToMask:     numberToMask,
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
	r, err := client.DeidentifyContent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	// Print the result.
	fmt.Fprint(w, r.GetItem().GetValue())
}

// [END dlp_deidentify_masking]

// [START dlp_deidentify_date_shift]

// deidentifyDateShift shifts dates found in the input between lowerBoundDays and
// upperBoundDays.
func deidentifyDateShift(w io.Writer, client *dlp.Client, project string, lowerBoundDays, upperBoundDays int32, input string) {
	// Create a configured request.
	req := &dlppb.DeidentifyContentRequest{
		Parent: "projects/" + project,
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
				InfoTypeTransformations: &dlppb.InfoTypeTransformations{
					Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
						{
							InfoTypes: []*dlppb.InfoType{}, // Match all info types.
							PrimitiveTransformation: &dlppb.PrimitiveTransformation{
								Transformation: &dlppb.PrimitiveTransformation_DateShiftConfig{
									DateShiftConfig: &dlppb.DateShiftConfig{
										LowerBoundDays: lowerBoundDays,
										UpperBoundDays: upperBoundDays,
									},
								},
							},
						},
					},
				},
			},
		},
		// The InspectConfig is used to identify the DATE fields.
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: []*dlppb.InfoType{
				{
					Name: "DATE",
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
	r, err := client.DeidentifyContent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	// Print the result.
	fmt.Fprint(w, r.GetItem().GetValue())
}

// [END dlp_deidentify_date_shift]

// [START dlp_deidentify_fpe]

// deidentifyFPE deidentifies the input with FPE (Format Preserving Encryption).
// keyFileName is the file name with the KMS wrapped key and cryptoKeyName is the
// full KMS key resource name used to wrap the key. surrogateInfoType is an
// optional identifier needed for reidentification. surrogateInfoType can be any
// value not found in your input.
func deidentifyFPE(w io.Writer, client *dlp.Client, project, input string, infoTypes []string, keyFileName, cryptoKeyName, surrogateInfoType string) {
	// Convert the info type strings to a list of InfoTypes.
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}
	// Read the key file.
	keyBytes, err := ioutil.ReadFile(keyFileName)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	// Create a configured request.
	req := &dlppb.DeidentifyContentRequest{
		Parent: "projects/" + project,
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: i,
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
													WrappedKey:    keyBytes,
													CryptoKeyName: cryptoKeyName,
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
	r, err := client.DeidentifyContent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	// Print the result.
	fmt.Fprint(w, r.GetItem().GetValue())
}

// [END dlp_deidentify_fpe]

// [START dlp_reidentify_fpe]

// reidentifyFPE reidentifies the input with FPE (Format Preserving Encryption).
// keyFileName is the file name with the KMS wrapped key and cryptoKeyName is the
// full KMS key resource name used to wrap the key. surrogateInfoType is an
// the identifier used during deidentification.
func reidentifyFPE(w io.Writer, client *dlp.Client, project, s, keyFileName, cryptoKeyName, surrogateInfoType string) {
	// Read the key file.
	keyBytes, err := ioutil.ReadFile(keyFileName)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	// Create a configured request.
	req := &dlppb.ReidentifyContentRequest{
		Parent: "projects/" + project,
		ReidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
				InfoTypeTransformations: &dlppb.InfoTypeTransformations{
					Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
						{
							InfoTypes: []*dlppb.InfoType{},
							PrimitiveTransformation: &dlppb.PrimitiveTransformation{
								Transformation: &dlppb.PrimitiveTransformation_CryptoReplaceFfxFpeConfig{
									CryptoReplaceFfxFpeConfig: &dlppb.CryptoReplaceFfxFpeConfig{
										CryptoKey: &dlppb.CryptoKey{
											Source: &dlppb.CryptoKey_KmsWrapped{
												KmsWrapped: &dlppb.KmsWrappedCryptoKey{
													WrappedKey:    keyBytes,
													CryptoKeyName: cryptoKeyName,
												},
											},
										},
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
				Value: s,
			},
		},
	}
	// Send the request.
	r, err := client.ReidentifyContent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	// Print the result.
	fmt.Fprint(w, r.GetItem().GetValue())
}

// [END dlp_reidentify_fpe]
