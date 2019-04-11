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

package findings

// [START update_finding_source_properties]
import (
	"context"
	"fmt"
	"io"
	"time"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"github.com/golang/protobuf/ptypes"
	structpb "github.com/golang/protobuf/ptypes/struct"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1"
	"google.golang.org/genproto/protobuf/field_mask"
)

// updateFindingSourceProperties demonstrates how to update a security finding
// in CSCC. findingName is the full resource name of the finding to update.
func updateFindingSourceProperties(w io.Writer, findingName string) error {
	// findingName := "organizations/111122222444/sources/1234/findings/findingid"
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %v", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.
	// Use now as the eventTime for the security finding.
	eventTime, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return fmt.Errorf("TimestampProto: %v", err)
	}

	req := &securitycenterpb.UpdateFindingRequest{
		Finding: &securitycenterpb.Finding{
			Name:      findingName,
			EventTime: eventTime,
			SourceProperties: map[string]*structpb.Value{
				"s_value": {
					Kind: &structpb.Value_StringValue{StringValue: "new_string_example"},
				},
			},
		},
		// Needed to only update the specific source property s_value
		// and EventTime. EventTime is a required field.
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"event_time", "source_properties.s_value"},
		},
	}

	finding, err := client.UpdateFinding(ctx, req)
	if err != nil {
		return fmt.Errorf("UpdateFinding: %v", err)
	}
	fmt.Fprintf(w, "Finding updated: %s\n", finding.Name)
	fmt.Fprintf(w, "Finding state: %v\n", finding.State)
	fmt.Fprintf(w, "Event time (Epoch Seconds): %d\n", eventTime.Seconds)
	fmt.Fprintf(w, "Source Properties:\n")
	for k, v := range finding.SourceProperties {
		fmt.Fprintf(w, "%s = %v\n", k, v)
	}
	return nil
}

// [END update_finding_source_properties]
