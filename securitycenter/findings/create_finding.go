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

// [START create_finding]
import (
	"context"
	"fmt"
	"io"
	"time"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"github.com/golang/protobuf/ptypes"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1"
)

// createFinding demonstrates how to create a new security finding in CSCC.
// sourceName is the full resource name of the source the finding should
// be associated with.
func createFinding(w io.Writer, sourceName string) error {
	// sourceName := "organizations/111122222444/sources/1234"
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

	req := &securitycenterpb.CreateFindingRequest{
		Parent:    sourceName,
		FindingId: "samplefindingid",
		Finding: &securitycenterpb.Finding{
			State: securitycenterpb.Finding_ACTIVE,
			// Resource the finding is associated with. This is an
			// example any resource identifier can be used.
			ResourceName: "//cloudresourcemanager.googleapis.com/organizations/11232",
			// A free-form category.
			Category: "MEDIUM_RISK_ONE",
			// The time associated with discovering the issue.
			EventTime: eventTime,
		},
	}
	finding, err := client.CreateFinding(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateFinding: %v", err)
	}
	fmt.Fprintf(w, "New finding created: %s\n", finding.Name)
	fmt.Fprintf(w, "Event time (Epoch Seconds): %d\n", eventTime.Seconds)
	return nil
}

// [END create_finding]
