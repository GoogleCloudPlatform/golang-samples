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

// [START test_iam]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	iam "google.golang.org/genproto/googleapis/iam/v1"
)

// testIam demonstrates how to determine if your service user has appropriate
// access to create and update findings, it writes permissions to w.
// sourceName is the full resource name of the source to test for permissions.
func testIam(w io.Writer, sourceName string) error {
	// sourceName := "organizations/111122222444/sources/1234"
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %v", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.
	// Check for create/update Permissions.
	req := &iam.TestIamPermissionsRequest{
		Resource:    sourceName,
		Permissions: []string{"securitycenter.findings.update"},
	}

	policy, err := client.TestIamPermissions(ctx, req)
	if err != nil {
		return fmt.Errorf("Error getting IAM policy: %v", err)
	}
	fmt.Fprintf(w, "Permision to create/update findings? %t",
		len(policy.Permissions) > 0)

	// Check for updating state Permissions
	req = &iam.TestIamPermissionsRequest{
		Resource:    sourceName,
		Permissions: []string{"securitycenter.findings.setState"},
	}

	policy, err = client.TestIamPermissions(ctx, req)
	if err != nil {
		return fmt.Errorf("Error getting IAM policy: %v", err)
	}
	fmt.Fprintf(w, "Permision to update state? %t",
		len(policy.Permissions) > 0)

	return nil
}

// [END test_iam]
