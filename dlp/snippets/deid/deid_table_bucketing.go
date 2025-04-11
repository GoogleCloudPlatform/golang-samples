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

// [START dlp_deidentify_table_bucketing]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// deIdentifyTableBucketing de-identifies data using table bucketing
func deIdentifyTableBucketing(w io.Writer, projectID string) error {
	// projectId := "your-project-id"
	// table := "your-table-value"

	row1 := &dlppb.Table_Row{
		Values: []*dlppb.Value{
			{Type: &dlppb.Value_StringValue{StringValue: "22"}},
			{Type: &dlppb.Value_StringValue{StringValue: "Jane Austen"}},
			{Type: &dlppb.Value_StringValue{StringValue: "21"}},
		},
	}

	row2 := &dlppb.Table_Row{
		Values: []*dlppb.Value{
			{Type: &dlppb.Value_StringValue{StringValue: "55"}},
			{Type: &dlppb.Value_StringValue{StringValue: "Mark Twain"}},
			{Type: &dlppb.Value_StringValue{StringValue: "75"}},
		},
	}

	row3 := &dlppb.Table_Row{
		Values: []*dlppb.Value{
			{Type: &dlppb.Value_StringValue{StringValue: "101"}},
			{Type: &dlppb.Value_StringValue{StringValue: "Charles Dickens"}},
			{Type: &dlppb.Value_StringValue{StringValue: "95"}},
		},
	}

	table := &dlppb.Table{
		Headers: []*dlppb.FieldId{
			{Name: "AGE"},
			{Name: "PATIENT"},
			{Name: "HAPPINESS SCORE"},
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

	// Specify how the content should be de-identified.
	fixedSizeBucketingConfig := &dlppb.FixedSizeBucketingConfig{
		BucketSize: 10,
		LowerBound: &dlppb.Value{
			Type: &dlppb.Value_IntegerValue{
				IntegerValue: 0,
			},
		},
		UpperBound: &dlppb.Value{
			Type: &dlppb.Value_IntegerValue{
				IntegerValue: 100,
			},
		},
	}
	primitiveTransformation := &dlppb.PrimitiveTransformation_FixedSizeBucketingConfig{
		FixedSizeBucketingConfig: fixedSizeBucketingConfig,
	}

	// Specify field to be encrypted.
	fieldId := &dlppb.FieldId{
		Name: "HAPPINESS SCORE",
	}

	// Associate the encryption with the specified field.
	fieldTransformation := &dlppb.FieldTransformation{
		Transformation: &dlppb.FieldTransformation_PrimitiveTransformation{
			PrimitiveTransformation: &dlppb.PrimitiveTransformation{
				Transformation: primitiveTransformation,
			},
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

	// Construct the de-identification request to be sent by the client.
	req := &dlppb.DeidentifyContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_RecordTransformations{
				RecordTransformations: recordTransformations,
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
	fmt.Fprintf(w, "Table after de-identification : %v", resp.GetItem().GetTable())
	return nil
}

// [END dlp_deidentify_table_bucketing]
