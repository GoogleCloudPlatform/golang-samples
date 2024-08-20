// Copyright 2024 Google LLC
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

package findingsv2

// [START securitycenter_update_finding_state_v2]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	"cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
)

// updateFindingState demonstrates how to update a security finding's state
// in CSCC.  findingName is the full resource name of the finding to update.
func setFindingState(w io.Writer, findingName string) error {
	// findingName := "organizations/111122222444/sources/1234/locations/global"
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	req := &securitycenterpb.SetFindingStateRequest{
		Name:  findingName,
		State: securitycenterpb.Finding_INACTIVE,
		// New state is effective immediately.
	}

	finding, err := client.SetFindingState(ctx, req)
	if err != nil {
		return fmt.Errorf("SetFindingState: %w", err)
	}

	fmt.Fprintf(w, "Finding updated: %s\n", finding.Name)
	fmt.Fprintf(w, "Finding state: %v\n", finding.State)
	fmt.Fprintf(w, "Event time (Epoch Seconds): %d\n", finding.EventTime.Seconds)

	return nil
}

// [END securitycenter_update_finding_state_v2]
