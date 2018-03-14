/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"golang.org/x/net/context"

	dlp "cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

func mask(w io.Writer, client *dlp.Client, project, input, maskingCharacter string, numberToMask int32) {
	// [START dlp_deidentify_masking]
	req := &dlppb.DeidentifyContentRequest{
		Parent: "projects/" + project,
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
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: input,
			},
		},
	}
	r, err := client.DeidentifyContent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, r.GetItem().GetValue())
	// [END dlp_deidentify_masking]
}

func deidentifyDateShift(w io.Writer, client *dlp.Client, project string, lowerBoundDays, upperBoundDays int32, s string) {
	// [START dlp_deidentify_date_shift]
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
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: []*dlppb.InfoType{
				{
					Name: "DATE",
				},
			},
		},
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: s,
			},
		},
	}
	r, err := client.DeidentifyContent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, r.GetItem().GetValue())
	// [END dlp_deidentify_date_shift]
}

func deidentifyFPE(w io.Writer, client *dlp.Client, project, s, keyFileName, cryptoKeyName, surrogateInfoType string) {
	// [START dlp_deidentify_fpe]
	b, err := ioutil.ReadFile(keyFileName)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	req := &dlppb.DeidentifyContentRequest{
		Parent: "projects/" + project,
		DeidentifyConfig: &dlppb.DeidentifyConfig{
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
													WrappedKey:    b,
													CryptoKeyName: cryptoKeyName,
												},
											},
										},
										Alphabet: &dlppb.CryptoReplaceFfxFpeConfig_CommonAlphabet{
											CommonAlphabet: dlppb.CryptoReplaceFfxFpeConfig_ALPHA_NUMERIC,
										},
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
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: s,
			},
		},
	}
	r, err := client.DeidentifyContent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, r.GetItem().GetValue())
	// [END dlp_deidentify_fpe]
}

func reidentifyFPE(w io.Writer, client *dlp.Client, project, s, keyFileName, cryptoKeyName, surrogateInfoType string) {
	// [START reidentify_fpe]
	b, err := ioutil.ReadFile(keyFileName)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}
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
													WrappedKey:    b,
													CryptoKeyName: cryptoKeyName,
												},
											},
										},
										Alphabet: &dlppb.CryptoReplaceFfxFpeConfig_CommonAlphabet{
											CommonAlphabet: dlppb.CryptoReplaceFfxFpeConfig_ALPHA_NUMERIC,
										},
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
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: s,
			},
		},
	}
	r, err := client.ReidentifyContent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, r.GetItem().GetValue())
	// [END reidentify_fpe]
}
