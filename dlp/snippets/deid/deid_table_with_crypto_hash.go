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

// [START dlp_deidentify_table_with_crypto_hash]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// deIdentifyTableWithCryptoHash transforms findings using a cryptographic hash transformation.
func deIdentifyTableWithCryptoHash(w io.Writer, projectID, transientKeyName string) error {
	// projectId := "your-project-id"
	// transientKeyName := "YOUR_TRANSIENT_CRYPTO_KEY_NAME"

	row1 := &dlppb.Table_Row{
		Values: []*dlppb.Value{
			{Type: &dlppb.Value_StringValue{StringValue: "user1@example.org"}},
			{Type: &dlppb.Value_StringValue{StringValue: "my email is user1@example.org and phone is 858-555-0222"}},
		},
	}

	row2 := &dlppb.Table_Row{
		Values: []*dlppb.Value{
			{Type: &dlppb.Value_StringValue{StringValue: "user2@example.org"}},
			{Type: &dlppb.Value_StringValue{StringValue: "my email is user2@example.org and phone is 858-555-0232"}},
		},
	}

	row3 := &dlppb.Table_Row{
		Values: []*dlppb.Value{
			{Type: &dlppb.Value_StringValue{StringValue: "user3@example.org"}},
			{Type: &dlppb.Value_StringValue{StringValue: "my email is user3@example.org and phone is 858-555-0224"}},
		},
	}

	tableToDeidentify := &dlppb.Table{
		Headers: []*dlppb.FieldId{
			{Name: "userid"},
			{Name: "comments"},
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
			Table: tableToDeidentify,
		},
	}

	// Specify the type of info the inspection will look for.
	// See https://cloud.google.com/dlp/docs/infotypes-reference for complete list of info types
	infoTypes := []*dlppb.InfoType{
		{Name: "PHONE_NUMBER"},
		{Name: "EMAIL_ADDRESS"},
	}

	inspectConfig := &dlppb.InspectConfig{
		InfoTypes: infoTypes,
	}

	// Specify the transient key which will encrypt the data.
	if transientKeyName == "" {
		transientKeyName = "YOUR_TRANSIENT_CRYPTO_KEY_NAME"
	}

	// Specify the transient key which will encrypt the data.
	cryptoKey := &dlppb.CryptoKey{
		Source: &dlppb.CryptoKey_Transient{
			Transient: &dlppb.TransientCryptoKey{
				Name: transientKeyName,
			},
		},
	}

	// Specify how the info from the inspection should be encrypted.
	cryptoHashConfig := &dlppb.CryptoHashConfig{
		CryptoKey: cryptoKey,
	}

	// Define type of de-identification as cryptographic hash transformation.
	primitiveTransformation := &dlppb.PrimitiveTransformation_CryptoHashConfig{
		CryptoHashConfig: cryptoHashConfig,
	}

	infoTypeTransformation := &dlppb.InfoTypeTransformations_InfoTypeTransformation{
		InfoTypes: infoTypes,
		PrimitiveTransformation: &dlppb.PrimitiveTransformation{
			Transformation: primitiveTransformation,
		},
	}

	transformations := &dlppb.InfoTypeTransformations{
		Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
			infoTypeTransformation,
		},
	}

	// Specify the config for the de-identify request.
	deidentifyConfig := &dlppb.DeidentifyConfig{
		Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
			InfoTypeTransformations: transformations,
		},
	}

	// Construct the de-identification request to be sent by the client.
	req := &dlppb.DeidentifyContentRequest{
		Parent:           fmt.Sprintf("projects/%s/locations/global", projectID),
		DeidentifyConfig: deidentifyConfig,
		InspectConfig:    inspectConfig,
		Item:             contentItem,
	}

	// Send the request.
	resp, err := client.DeidentifyContent(ctx, req)
	if err != nil {
		return err
	}

	// Print the results.
	fmt.Fprintf(w, "Table after de-identification : %v", resp.GetItem().GetTable())
	return nil
}

// [END dlp_deidentify_table_with_crypto_hash]
