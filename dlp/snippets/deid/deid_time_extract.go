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

// [START dlp_deidentify_time_extract]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// deIdentifyTimeExtract De-identifies a table by extracting specific parts
// of the time (year in this case) from designated fields.
func deIdentifyTimeExtract(w io.Writer, projectID string) error {

	row1 := &dlppb.Table_Row{
		Values: []*dlppb.Value{
			{Type: &dlppb.Value_StringValue{StringValue: "Ann"}},
			{Type: &dlppb.Value_StringValue{StringValue: "01/01/1970"}},
			{Type: &dlppb.Value_StringValue{StringValue: "4532908762519852"}},
			{Type: &dlppb.Value_StringValue{StringValue: "07/21/1996"}},
		},
	}

	row2 := &dlppb.Table_Row{
		Values: []*dlppb.Value{
			{Type: &dlppb.Value_StringValue{StringValue: "James"}},
			{Type: &dlppb.Value_StringValue{StringValue: "03/06/1988"}},
			{Type: &dlppb.Value_StringValue{StringValue: "4301261899725540"}},
			{Type: &dlppb.Value_StringValue{StringValue: "04/09/2001"}},
		},
	}

	// Specify the table to de-identify.
	tableToDeidentify := &dlppb.Table{
		Headers: []*dlppb.FieldId{
			{Name: "Name"},
			{Name: "Birth Date"},
			{Name: "Credit Card"},
			{Name: "Register Date"},
		},
		Rows: []*dlppb.Table_Row{
			{Values: row1.Values},
			{Values: row2.Values},
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

	// Specify the time part to extract.
	timePartConfig := &dlppb.PrimitiveTransformation_TimePartConfig{
		TimePartConfig: &dlppb.TimePartConfig{
			PartToExtract: dlppb.TimePartConfig_YEAR,
		},
	}

	transformation := &dlppb.PrimitiveTransformation{
		Transformation: timePartConfig,
	}

	// Specify which fields the TimePart should apply too.
	dateFields := []*dlppb.FieldId{
		{Name: "Birth Date"},
		{Name: "Register Date"},
	}

	fieldTransformation := &dlppb.FieldTransformation{
		Fields: dateFields,
		Transformation: &dlppb.FieldTransformation_PrimitiveTransformation{
			PrimitiveTransformation: transformation,
		},
	}

	recordTransformations := []*dlppb.RecordTransformations{
		{
			FieldTransformations: []*dlppb.FieldTransformation{
				fieldTransformation,
			},
		},
	}

	// Specify the config for the de-identify request.
	deidentifyConfig := &dlppb.DeidentifyConfig{
		Transformation: &dlppb.DeidentifyConfig_RecordTransformations{
			RecordTransformations: recordTransformations[0],
		},
	}

	// Construct the de-identification request to be sent by the client.
	req := &dlppb.DeidentifyContentRequest{
		Parent:           fmt.Sprintf("projects/%s/locations/global", projectID),
		DeidentifyConfig: deidentifyConfig,
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

// [END dlp_deidentify_time_extract]
