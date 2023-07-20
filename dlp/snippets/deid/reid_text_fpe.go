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

// [START dlp_reidentify_text_fpe]
import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// ReidTextDataWithFPE re-identifies text data with FPE
func reidTextDataWithFPE(w io.Writer, projectID, textToReidentify, kmsKeyName, wrappedAESKey, surrogateInfoType string) error {
	// projectId := "my-project-id"
	// textToReidentify := "My SSN is AGE(9):FHKByqA13"
	/* kmsKeyName :=  "projects/YOUR_PROJECT/"
	   + "locations/YOUR_KEYRING_REGION/"
	   + "keyRings/YOUR_KEYRING_NAME/"
	   + "cryptoKeys/YOUR_KEY_NAME"
	*/
	// wrappedAesKey := "YOUR_ENCRYPTED_AES_256_KEY"
	// surrogateInfoType := "AGE"

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

	// Specify what content you want the service to re-identify.
	contentItem := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: textToReidentify,
		},
	}

	// Specify the type of info the inspection will re-identify. This must use the same custom
	// into type that was used as a surrogate during the initial encryption.
	infotype := &dlppb.InfoType{
		Name: surrogateInfoType,
	}

	customInfotype := &dlppb.CustomInfoType{
		InfoType: infotype,
		Type: &dlppb.CustomInfoType_SurrogateType_{
			SurrogateType: &dlppb.CustomInfoType_SurrogateType{},
		},
	}

	inspectConfig := &dlppb.InspectConfig{
		CustomInfoTypes: []*dlppb.CustomInfoType{
			customInfotype,
		},
	}

	// Specify an encrypted AES-256 key and the name of the Cloud KMS key that encrypted it.
	kmsWrappedCryptoKey, err := base64.StdEncoding.DecodeString(wrappedAESKey)
	if err != nil {
		return err
	}

	cryptoKey := &dlppb.CryptoKey{
		Source: &dlppb.CryptoKey_KmsWrapped{
			KmsWrapped: &dlppb.KmsWrappedCryptoKey{
				WrappedKey:    kmsWrappedCryptoKey,
				CryptoKeyName: kmsKeyName,
			},
		},
	}

	// Specify how to un-encrypt the previously de-identified information.
	cryptoReplaceFfxFpeConfig := &dlppb.CryptoReplaceFfxFpeConfig{
		CryptoKey: cryptoKey,
		Alphabet: &dlppb.CryptoReplaceFfxFpeConfig_CommonAlphabet{
			CommonAlphabet: dlppb.CryptoReplaceFfxFpeConfig_ALPHA_NUMERIC,
		},
		SurrogateInfoType: infotype,
	}

	primitiveTransformation := &dlppb.PrimitiveTransformation_CryptoReplaceFfxFpeConfig{
		CryptoReplaceFfxFpeConfig: cryptoReplaceFfxFpeConfig,
	}

	InfoTypeTransformations := &dlppb.InfoTypeTransformations{
		Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
			{
				PrimitiveTransformation: &dlppb.PrimitiveTransformation{
					Transformation: primitiveTransformation,
				},
				InfoTypes: []*dlppb.InfoType{
					infotype,
				},
			},
		},
	}

	// Combine configurations into a request for the service.
	reidentifyConfig := &dlppb.DeidentifyConfig{
		Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
			InfoTypeTransformations: InfoTypeTransformations,
		},
	}

	// Construct the Inspect request to be sent by the client.
	req := &dlppb.ReidentifyContentRequest{
		Parent:           fmt.Sprintf("projects/%s/locations/global", projectID),
		ReidentifyConfig: reidentifyConfig,
		Item:             contentItem,
		InspectConfig:    inspectConfig,
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

// [END dlp_reidentify_text_fpe]
