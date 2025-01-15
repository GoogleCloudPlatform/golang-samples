// Copyright 2020 Google LLC
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

package policytagmanager

// [START data_catalog_ptm_test_iam_permissions]
import (
	"context"
	"fmt"
	"io"
	"strings"

	datacatalog "cloud.google.com/go/datacatalog/apiv1beta1"
	"cloud.google.com/go/iam/apiv1/iampb"
)

// testIAMPermissions demonstrates a caller probing what permissions they hold on a given taxonomy
// or policy tag resource.
func testIAMPermissions(w io.Writer, resourceID string, permissions []string) error {
	// permissions := []string{"datacatalog.categories.fineGrainedGet"}
	ctx := context.Background()
	policyClient, err := datacatalog.NewPolicyTagManagerClient(ctx)
	if err != nil {
		return fmt.Errorf("datacatalog.NewPolicyTagManagerClient: %w", err)
	}
	defer policyClient.Close()

	req := &iampb.TestIamPermissionsRequest{
		Resource:    resourceID,
		Permissions: permissions,
	}
	resp, err := policyClient.TestIamPermissions(ctx, req)
	if err != nil {
		return fmt.Errorf("TestIamPermissions: %w", err)
	}
	fmt.Fprintf(w, "Testing the permissions on %s, of the %d permissions probed, caller has %d permissions", resourceID, len(permissions), len(resp.Permissions))
	if len(resp.Permissions) > 0 {
		fmt.Fprintf(w, ": %s", strings.Join(resp.Permissions, ", "))
	}
	fmt.Fprintln(w)
	return nil
}

// [END data_catalog_ptm_test_iam_permissions]
