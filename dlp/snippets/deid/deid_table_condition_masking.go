// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deid

// [START dlp_deidentify_table_condition_masking]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// deidentifyTableMaskingCondition de-identifies the table data using masking
// and conditional logic
func deidentifyTableMaskingCondition(w io.Writer, projectID string) error {
	// projectId := "your-project-id"

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
	charMaskConfig := &dlppb.CharacterMaskConfig{
		MaskingCharacter: "*",
	}
	primitiveTransformation := &dlppb.PrimitiveTransformation_CharacterMaskConfig{
		CharacterMaskConfig: charMaskConfig,
	}

	// Specify field to be de-identified.
	fieldId := &dlppb.FieldId{
		Name: "HAPPINESS SCORE",
	}

	// Apply the condition to the records present in table.
	recordCondition := &dlppb.RecordCondition_Condition{
		Field:    fieldId,
		Operator: dlppb.RelationalOperator_GREATER_THAN,
		Value: &dlppb.Value{
			Type: &dlppb.Value_IntegerValue{
				IntegerValue: 89,
			},
		},
	}

	expression := &dlppb.RecordCondition_Expressions{
		Type: &dlppb.RecordCondition_Expressions_Conditions{
			Conditions: &dlppb.RecordCondition_Conditions{
				Conditions: []*dlppb.RecordCondition_Condition{
					recordCondition,
				},
			},
		},
	}

	// Specify when the above field should be de-identified.
	condition := &dlppb.RecordCondition{
		// Apply the condition to records
		Expressions: expression,
	}

	// Associate the de-identification and conditions with the specified field.
	fieldTransformation := &dlppb.FieldTransformation{
		Fields: []*dlppb.FieldId{
			fieldId,
		},
		Transformation: &dlppb.FieldTransformation_PrimitiveTransformation{
			PrimitiveTransformation: &dlppb.PrimitiveTransformation{
				Transformation: primitiveTransformation,
			},
		},
		Condition: condition,
	}

	recordTransformations := &dlppb.RecordTransformations{
		FieldTransformations: []*dlppb.FieldTransformation{
			fieldTransformation,
		},
	}

	// Combine configurations into a request for the service.
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

// [END dlp_deidentify_table_condition_masking]
