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

// [START dlp_deidentify_deterministic]

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// deIdentifyDeterministicEncryption de-identifies through deterministic encryption
func deIdentifyDeterministicEncryption(w io.Writer, projectID, inputStr string, infoTypeNames []string, keyFileName, cryptoKeyName, surrogateInfoType string) error {
	// projectId := "your-project-id"
	// inputStr := "My SSN is 111111111"
	// infoTypeNames := []string{"US_SOCIAL_SECURITY_NUMBER"}
	/* keyFileName :=  "projects/YOUR_PROJECT/"
	   + "locations/YOUR_KEYRING_REGION/"
	   + "keyRings/YOUR_KEYRING_NAME/"
	   + "cryptoKeys/YOUR_KEY_NAME"
	*/
	// cryptoKeyName := "YOUR_ENCRYPTED_AES_256_KEY"
	// surrogateInfoType := "SSN_TOKEN"

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

	// Specify an encrypted AES-256 key and the name of the Cloud KMS key that encrypted it.
	wrappedKey, err := base64.StdEncoding.DecodeString(cryptoKeyName)
	if err != nil {
		return err
	}

	// Specify the type of info the inspection will look for.
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types
	var infoTypes []*dlppb.InfoType
	for _, it := range infoTypeNames {
		infoTypes = append(infoTypes, &dlppb.InfoType{Name: it})
	}

	// Specify the key used by the encryption function for deterministic encryption.
	cryptoReplaceDeterministicConfig := &dlppb.CryptoDeterministicConfig{
		CryptoKey: &dlppb.CryptoKey{
			Source: &dlppb.CryptoKey_KmsWrapped{
				KmsWrapped: &dlppb.KmsWrappedCryptoKey{
					WrappedKey:    wrappedKey,
					CryptoKeyName: keyFileName,
				},
			},
		},
		SurrogateInfoType: &dlppb.InfoType{
			Name: surrogateInfoType,
		},
	}

	// Specifying the info-types to look for.
	inspectConfig := &dlppb.InspectConfig{
		InfoTypes: infoTypes,
	}

	// Specify what content you want the service to de-identify.
	contentItem := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: inputStr,
		},
	}

	// Specifying the deterministic crypto.
	primitiveTransformation := &dlppb.PrimitiveTransformation{
		Transformation: &dlppb.PrimitiveTransformation_CryptoDeterministicConfig{
			CryptoDeterministicConfig: cryptoReplaceDeterministicConfig,
		},
	}

	// Construct a de-identification config for de-identify deterministic request.
	deIdentifyConfig := &dlppb.DeidentifyConfig{
		Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
			InfoTypeTransformations: &dlppb.InfoTypeTransformations{
				Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
					{
						PrimitiveTransformation: primitiveTransformation,
					},
				},
			},
		},
	}

	// Construct the de-identification request to be sent by the client.
	req := &dlppb.DeidentifyContentRequest{
		Parent:           fmt.Sprintf("projects/%s/locations/global", projectID),
		DeidentifyConfig: deIdentifyConfig,
		InspectConfig:    inspectConfig,
		Item:             contentItem,
	}

	// Send the request.
	resp, err := client.DeidentifyContent(ctx, req)
	if err != nil {
		return err
	}

	// Print the results.
	fmt.Fprintf(w, "output : %v", resp.GetItem().GetValue())
	return nil
}

// [END dlp_deidentify_deterministic]
