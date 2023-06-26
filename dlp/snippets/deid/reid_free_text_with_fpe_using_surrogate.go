// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package deid

// [START dlp_reidentify_free_text_with_fpe_using_surrogate]
import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// reidentifyFreeTextWithFPEUsingSurrogate re-identifies free text with FPE using surrogate
func reidentifyFreeTextWithFPEUsingSurrogate(w io.Writer, projectID, inputStr, surrogateType, unWrappedKey string) error {
	// projectId := "your-project-id"
	// inputStr := "My phone number is PHONE_TOKEN(10):4068372189"
	// surrogateType := "PHONE_TOKEN"
	// unWrappedKey := "hXribiHR6kvVKmj5ptTZLRTj8+u0eBQzwcrEPwyGhe8="

	ctx := context.Background()

	// Initialize a client once and reuse it to send multiple requests. Clients
	// are safe to use across goroutines. When the client is no longer needed,
	// call the Close method to cleanup its resources.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		return err
	}

	// Closing the client safely cleans up background resources.
	defer client.Close()

	// The wrapped key is base64-encoded, but the library expects a binary
	// string, so decode it here.
	keyBytes, err := base64.StdEncoding.DecodeString(unWrappedKey)
	if err != nil {
		return err
	}

	// using the unwrapped key and surrogate info type construct cryptoReplaceConfig
	cryptoReplaceConfig := &dlppb.CryptoReplaceFfxFpeConfig{
		CryptoKey: &dlppb.CryptoKey{
			Source: &dlppb.CryptoKey_Unwrapped{
				Unwrapped: &dlppb.UnwrappedCryptoKey{
					Key: keyBytes,
				},
			},
		},
		Alphabet: &dlppb.CryptoReplaceFfxFpeConfig_CommonAlphabet{
			CommonAlphabet: dlppb.CryptoReplaceFfxFpeConfig_NUMERIC,
		},
		SurrogateInfoType: &dlppb.InfoType{
			Name: surrogateType,
		},
	}

	// construct infoType transformations
	infoTypeTransformations := &dlppb.InfoTypeTransformations{
		Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
			{
				PrimitiveTransformation: &dlppb.PrimitiveTransformation{
					Transformation: &dlppb.PrimitiveTransformation_CryptoReplaceFfxFpeConfig{
						CryptoReplaceFfxFpeConfig: cryptoReplaceConfig,
					},
				},
			},
		},
	}

	// construct re-identify config
	reIdentifyConfig := &dlppb.DeidentifyConfig{
		Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
			InfoTypeTransformations: infoTypeTransformations,
		},
	}

	// construct inspect config to pass in req
	inspectConfig := &dlppb.InspectConfig{
		CustomInfoTypes: []*dlppb.CustomInfoType{
			{
				InfoType: &dlppb.InfoType{
					Name: surrogateType,
				},
				Type: &dlppb.CustomInfoType_SurrogateType_{
					SurrogateType: &dlppb.CustomInfoType_SurrogateType{},
				},
			},
		},
	}

	// Construct the de-identification request to be sent by the client.
	req := &dlppb.ReidentifyContentRequest{
		Parent:           fmt.Sprintf("projects/%s/locations/global", projectID),
		ReidentifyConfig: reIdentifyConfig,
		InspectConfig:    inspectConfig,
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: inputStr,
			},
		},
	}

	// Send the request.
	r, err := client.ReidentifyContent(ctx, req)
	if err != nil {
		return err
	}

	// Print the result.
	fmt.Fprintf(w, "output: %v", r.GetItem().GetValue())
	return nil
}

// [END dlp_reidentify_free_text_with_fpe_using_surrogate]
