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

// [START get_iam_policy_source]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	iam "google.golang.org/genproto/googleapis/iam/v1"
)

// getSourceIamPolicy prints the policy for sourceName to w and return it.
// sourceName is the full resource name of the source with the policy of interest.
func getSourceIamPolicy(w io.Writer, sourceName string) error {
	// sourceName := "organizations/111122222444/sources/1234"
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %v", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	req := &iam.GetIamPolicyRequest{
		Resource: sourceName,
	}

	policy, err := client.GetIamPolicy(ctx, req)
	if err != nil {
		return fmt.Errorf("GetIamPolicy(%s): %v", sourceName, err)
	}

	fmt.Fprintf(w, "Policy: %v", policy)
	return nil
}

// [END get_iam_policy_source]
