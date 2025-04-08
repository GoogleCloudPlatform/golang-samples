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

// [START dlp_deidentify_free_text_with_fpe_using_surrogate]
import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// deidentifyFreeTextWithFPEUsingSurrogate de-identifies free text with FPE by using a surrogate
func deidentifyFreeTextWithFPEUsingSurrogate(w io.Writer, projectID, inputStr, infoType, surrogateType, unwrappedKey string) error {
	// projectId := "your-project-id"
	// inputStr := "My phone number is 9824376677"
	// infoType := "PHONE_NUMBER"
	// surrogateType := "PHONE_TOKEN"
	// unwrappedKey := "The base64-encoded AES-256 key to use"

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

	// decoding the unwrapped key
	keyDecode, err := base64.StdEncoding.DecodeString(unwrappedKey)
	if err != nil {
		return err
	}

	inspectConfig := &dlppb.InspectConfig{
		InfoTypes: []*dlppb.InfoType{
			{Name: infoType},
		},
		MinLikelihood: dlppb.Likelihood_UNLIKELY,
	}

	// specify the content to be de-identified
	contentItem := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: inputStr,
		},
	}

	// Specify how the content should be encrypted.
	cryptoReplaceConfig := &dlppb.CryptoReplaceFfxFpeConfig{
		CryptoKey: &dlppb.CryptoKey{
			Source: &dlppb.CryptoKey_Unwrapped{
				Unwrapped: &dlppb.UnwrappedCryptoKey{
					Key: keyDecode,
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

	// Apply crypto replace config to primitive transformation.
	primitiveTransformation := &dlppb.PrimitiveTransformation{
		Transformation: &dlppb.PrimitiveTransformation_CryptoReplaceFfxFpeConfig{
			CryptoReplaceFfxFpeConfig: cryptoReplaceConfig,
		},
	}

	infoTypeTransformation := &dlppb.DeidentifyConfig_InfoTypeTransformations{
		InfoTypeTransformations: &dlppb.InfoTypeTransformations{
			Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
				{
					InfoTypes: []*dlppb.InfoType{
						{Name: infoType},
					},
					PrimitiveTransformation: primitiveTransformation,
				},
			},
		},
	}

	// Construct the de-identification request to be sent by the client.
	req := &dlppb.DeidentifyContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: infoTypeTransformation,
		},
		InspectConfig: inspectConfig,
		Item:          contentItem,
	}

	// Send the request.
	resp, err := client.DeidentifyContent(ctx, req)
	if err != nil {
		return err
	}

	// Print the results.
	fmt.Fprintf(w, "output: %v", resp.GetItem().GetValue())
	return nil
}

// [END dlp_deidentify_free_text_with_fpe_using_surrogate]
