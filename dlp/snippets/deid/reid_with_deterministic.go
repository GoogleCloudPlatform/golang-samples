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

// [START dlp_reidentify_deterministic]
import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// reidentifyWithDeterministic re-identifies content encrypted by deterministic encryption.
func reidentifyWithDeterministic(w io.Writer, projectID, inputStr, surrogateType, keyName, wrappedKey string) error {
	// projectId := "my-project-id"
	// inputStr := "EMAIL_ADDRESS_TOKEN(52):AVAx2eIEnIQP5jbNEr2j9wLOAd5m4kpSBR/0jjjGdAOmryzZbE/q"
	// surrogateType := "EMAIL_ADDRESS_TOKEN"
	/* keyName :=  "projects/YOUR_PROJECT/"
	   + "locations/YOUR_KEYRING_REGION/"
	   + "keyRings/YOUR_KEYRING_NAME/"
	   + "cryptoKeys/YOUR_KEY_NAME"
	*/
	// wrappedKey := "YOUR_ENCRYPTED_AES_256_KEY"

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
	keyBytes, err := base64.StdEncoding.DecodeString(wrappedKey)
	if err != nil {
		return err
	}

	// Create crypto deterministic config.
	cryptoDeterministicConfig := &dlppb.CryptoDeterministicConfig{
		CryptoKey: &dlppb.CryptoKey{
			Source: &dlppb.CryptoKey_KmsWrapped{
				KmsWrapped: &dlppb.KmsWrappedCryptoKey{
					WrappedKey:    keyBytes,
					CryptoKeyName: keyName,
				},
			},
		},
		SurrogateInfoType: &dlppb.InfoType{
			Name: surrogateType,
		},
	}

	// Create a config for primitive transformation.
	primitiveTransformation := &dlppb.PrimitiveTransformation{
		Transformation: &dlppb.PrimitiveTransformation_CryptoDeterministicConfig{
			CryptoDeterministicConfig: cryptoDeterministicConfig,
		},
	}

	transformation := &dlppb.DeidentifyConfig_InfoTypeTransformations{
		InfoTypeTransformations: &dlppb.InfoTypeTransformations{
			Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
				{
					PrimitiveTransformation: primitiveTransformation,
				},
			},
		},
	}

	// Construct config to re-identify the config.
	reIdentifyConfig := &dlppb.DeidentifyConfig{
		Transformation: transformation,
	}

	// Construct a config for inspection.
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

	// Item to be analyzed.
	item := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: inputStr,
		},
	}

	// Construct the Inspect request to be sent by the client.
	req := &dlppb.ReidentifyContentRequest{
		Parent:           fmt.Sprintf("projects/%s/locations/global", projectID),
		ReidentifyConfig: reIdentifyConfig,
		InspectConfig:    inspectConfig,
		Item:             item,
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

// [END dlp_reidentify_deterministic]
