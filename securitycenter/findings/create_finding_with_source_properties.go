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

// [START securitycenter_create_finding_with_source_properties]
import (
	"context"
	"fmt"
	"io"
	"time"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"cloud.google.com/go/securitycenter/apiv1/securitycenterpb"
	"github.com/golang/protobuf/ptypes"
	structpb "github.com/golang/protobuf/ptypes/struct"
)

// createFindingWithProperties demonstrates how to create a new security
// finding in CSCC that includes additional metadata via sourceProperties.
// sourceName is the full resource name of the source the finding should be
// associated with.
func createFindingWithProperties(w io.Writer, sourceName string) error {
	// sourceName := "organizations/111122222444/sources/1234"
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.
	// Use now as the eventTime for the security finding.
	eventTime, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return fmt.Errorf("TimestampProto: %w", err)
	}

	req := &securitycenterpb.CreateFindingRequest{
		Parent:    sourceName,
		FindingId: "samplefindingprops",
		Finding: &securitycenterpb.Finding{
			State: securitycenterpb.Finding_ACTIVE,
			// Resource the finding is associated with.  This is an
			// example any resource identifier can be used.
			ResourceName: "//cloudresourcemanager.googleapis.com/organizations/11232",
			// A free-form category.Error converting now
			Category: "MEDIUM_RISK_ONE",
			// The time associated with discovering the issue.
			EventTime: eventTime,
			// Define key-value pair metadata to include with the finding.
			SourceProperties: map[string]*structpb.Value{
				"s_value": {
					Kind: &structpb.Value_StringValue{StringValue: "string_example"},
				},
				"n_value": {
					Kind: &structpb.Value_NumberValue{NumberValue: 1234},
				},
			},
		},
	}

	finding, err := client.CreateFinding(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateFinding: %w", err)
	}
	fmt.Fprintf(w, "New finding created: %s\n", finding.Name)
	fmt.Fprintf(w, "Event time (Epoch Seconds): %d\n", eventTime.Seconds)
	fmt.Fprintf(w, "Source Properties:\n")
	for k, v := range finding.SourceProperties {
		fmt.Fprintf(w, "%s = %v\n", k, v)
	}

	return nil
}

// [END securitycenter_create_finding_with_source_properties]
