// Copyright 2019 Google LLC
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

// findings contains example snippets for working with findings
// and there parent resource "sources".
package findings

// [START update_finding_source_properties]
import (
	"context"
	"fmt"
	"time"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"github.com/golang/protobuf/ptypes"
	structpb "github.com/golang/protobuf/ptypes/struct"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1"
	"google.golang.org/genproto/protobuf/field_mask"
)

// updateFindingSourceProperties demonstrates how to update a security finding
// in CSCC.  sourceName is the full resource name of the source the finding
// should be associated with.  Returns the updated finding.
func updateFindingSourceProperties(sourceName string) (*securitycenterpb.Finding, error) {
	// sourceName := "organizations/111122222444/sources/1234"
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error instantiating client %v\n", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.
	// Use now as the eventTime for the security finding.
	eventTime, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		fmt.Printf("Error converting now: %v", err)
		return nil, err
	}

	// Define key-value pair metadata to include with the finding.
	properties := map[string]*structpb.Value{
		"s_value": &structpb.Value{
			Kind: &structpb.Value_StringValue{"string_example"}},
	}

	// Findings are a child resource of sources.
	findingName := fmt.Sprintf("%s/findings/samplefindingprops", sourceName)
	req := &securitycenterpb.UpdateFindingRequest{
		Finding: &securitycenterpb.Finding{
			Name:             findingName,
			EventTime:        eventTime,
			SourceProperties: properties,
		},
		// Needed to only update the specific source property s_value
		// and EventTime.  EventTime is a required field.
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"event_time", "source_properties.s_value"},
		},
	}

	finding, err := client.UpdateFinding(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Error updating finding: %v", err)
	}
	return finding, nil
}

// [END update_finding_source_properties]
