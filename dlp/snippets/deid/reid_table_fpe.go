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

// [START dlp_reidentify_table_fpe]
import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// reidTableDataWithFPE re-identifies table data with FPE
func reidTableDataWithFPE(w io.Writer, projectID, kmsKeyName, wrappedAesKey string) error {
	// projectId := "my-project-id"
	/* kmsKeyName :=  "projects/YOUR_PROJECT/"
	   + "locations/YOUR_KEYRING_REGION/"
	   + "keyRings/YOUR_KEYRING_NAME/"
	   + "cryptoKeys/YOUR_KEY_NAME"
	*/
	// wrappedAesKey := "YOUR_ENCRYPTED_AES_256_KEY"

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

	// Specify the table data that needs to be re-identified.
	tableToReIdentify := &dlppb.Table{
		Headers: []*dlppb.FieldId{
			{Name: "Employee ID"},
		},
		Rows: []*dlppb.Table_Row{
			{
				Values: []*dlppb.Value{
					{
						Type: &dlppb.Value_StringValue{
							StringValue: "90511",
						},
					},
				},
			},
		},
	}

	// Specify the content you want to re-identify.
	contentItem := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Table{
			Table: tableToReIdentify,
		},
	}

	// Specify an encrypted AES-256 key and the name of the Cloud KMS key that encrypted it.
	kmsWrappedCryptoKey, err := base64.StdEncoding.DecodeString(wrappedAesKey)
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
			CommonAlphabet: dlppb.CryptoReplaceFfxFpeConfig_NUMERIC,
		},
	}

	primitiveTransformation := &dlppb.PrimitiveTransformation{
		Transformation: &dlppb.PrimitiveTransformation_CryptoReplaceFfxFpeConfig{
			CryptoReplaceFfxFpeConfig: cryptoReplaceFfxFpeConfig,
		},
	}

	// Specify field to be decrypted.
	fieldId := &dlppb.FieldId{
		Name: "Employee ID",
	}

	// Associate the decryption with the specified field.
	fieldTransformation := &dlppb.FieldTransformation{
		Transformation: &dlppb.FieldTransformation_PrimitiveTransformation{
			PrimitiveTransformation: primitiveTransformation,
		},
		Fields: []*dlppb.FieldId{
			fieldId,
		},
	}

	recordTransformations := &dlppb.RecordTransformations{
		FieldTransformations: []*dlppb.FieldTransformation{
			fieldTransformation,
		},
	}

	deidentifyConfig := &dlppb.DeidentifyConfig{
		Transformation: &dlppb.DeidentifyConfig_RecordTransformations{
			RecordTransformations: recordTransformations,
		},
	}

	// Combine configurations into a request for the service.
	req := &dlppb.ReidentifyContentRequest{
		Parent:           fmt.Sprintf("projects/%s/locations/global", projectID),
		ReidentifyConfig: deidentifyConfig,
		Item:             contentItem,
	}

	// Send the request and receive response from the service.
	resp, err := client.ReidentifyContent(ctx, req)
	if err != nil {
		return err
	}

	// Print the result.
	fmt.Fprintf(w, "Table after re-identification : %v", resp.GetItem().GetTable())
	return nil
}

// [END dlp_reidentify_table_fpe]
