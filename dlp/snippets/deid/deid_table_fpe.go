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

// [START dlp_deidentify_table_fpe]
import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// deidentifyTableFPE de-identifies table data with format preserving encryption.
func deidentifyTableFPE(w io.Writer, projectID string, kmsKeyName, wrappedAESKey string) error {
	// projectId := "your-project-id"
	/* keyFileName :=  "projects/YOUR_PROJECT/"
	   + "locations/YOUR_KEYRING_REGION/"
	   + "keyRings/YOUR_KEYRING_NAME/"
	   + "cryptoKeys/YOUR_KEY_NAME"
	*/
	// wrappedAESKey := "YOUR_ENCRYPTED_AES_256_KEY"

	// define your table.
	row1 := &dlppb.Table_Row{
		Values: []*dlppb.Value{
			{Type: &dlppb.Value_StringValue{StringValue: "11111"}},
			{Type: &dlppb.Value_StringValue{StringValue: "2015"}},
			{Type: &dlppb.Value_StringValue{StringValue: "$10"}},
		},
	}

	row2 := &dlppb.Table_Row{
		Values: []*dlppb.Value{
			{Type: &dlppb.Value_StringValue{StringValue: "22222"}},
			{Type: &dlppb.Value_StringValue{StringValue: "2016"}},
			{Type: &dlppb.Value_StringValue{StringValue: "$20"}},
		},
	}

	row3 := &dlppb.Table_Row{
		Values: []*dlppb.Value{
			{Type: &dlppb.Value_StringValue{StringValue: "33333"}},
			{Type: &dlppb.Value_StringValue{StringValue: "2016"}},
			{Type: &dlppb.Value_StringValue{StringValue: "$15"}},
		},
	}

	table := &dlppb.Table{
		Headers: []*dlppb.FieldId{
			{Name: "Employee ID"},
			{Name: "Date"},
			{Name: "Compensation"},
		},
		Rows: []*dlppb.Table_Row{
			{Values: row1.Values},
			{Values: row2.Values},
			{Values: row3.Values},
		},
	}

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

	// Specify what content you want the service to de-identify.
	contentItem := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Table{
			Table: table,
		},
	}

	// Specify an encrypted AES-256 key and the name of the Cloud KMS key that encrypted it.
	kmsKeyDecode, err := base64.StdEncoding.DecodeString(wrappedAESKey)
	if err != nil {
		return fmt.Errorf("error in decoding key: %w", err)
	}

	kmsWrappedCryptoKey := &dlppb.KmsWrappedCryptoKey{
		WrappedKey:    kmsKeyDecode,
		CryptoKeyName: kmsKeyName,
	}

	cryptoKey := &dlppb.CryptoKey_KmsWrapped{
		KmsWrapped: kmsWrappedCryptoKey,
	}

	// Specify how the content should be encrypted.
	cryptoReplaceFfxFpeConfig := &dlppb.CryptoReplaceFfxFpeConfig{
		CryptoKey: &dlppb.CryptoKey{
			Source: cryptoKey,
		},
		Alphabet: &dlppb.CryptoReplaceFfxFpeConfig_CommonAlphabet{
			CommonAlphabet: dlppb.CryptoReplaceFfxFpeConfig_NUMERIC,
		},
	}

	primitiveTransformation := &dlppb.PrimitiveTransformation{
		Transformation: &dlppb.PrimitiveTransformation_CryptoReplaceFfxFpeConfig{
			CryptoReplaceFfxFpeConfig: cryptoReplaceFfxFpeConfig,
		},
	}

	// Specify field to be encrypted.
	fieldId := &dlppb.FieldId{
		Name: "Employee ID",
	}

	// Associate the encryption with the specified field.
	fieldTransformation := &dlppb.FieldTransformation{
		Fields: []*dlppb.FieldId{
			fieldId,
		},
		Transformation: &dlppb.FieldTransformation_PrimitiveTransformation{
			PrimitiveTransformation: primitiveTransformation,
		},
	}

	transformations := &dlppb.RecordTransformations{
		FieldTransformations: []*dlppb.FieldTransformation{
			fieldTransformation,
		},
	}

	// Construct the de-identification request to be sent by the client.
	req := &dlppb.DeidentifyContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_RecordTransformations{
				RecordTransformations: transformations,
			},
		},
		Item: contentItem,
	}

	// Send the request.
	resp, err := client.DeidentifyContent(ctx, req)
	if err != nil {
		return err
	}

	// Print the results.
	fmt.Fprintf(w, "De-identify Table after format-preserving encryption : %+v", resp.GetItem().GetTable())
	return nil
}

// [END dlp_deidentify_table_fpe]
