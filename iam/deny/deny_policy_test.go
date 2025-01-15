// Copyright 2022 Google LLC
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

package snippets

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	iam "cloud.google.com/go/iam/apiv2"
	"cloud.google.com/go/iam/apiv2/iampb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDenyPolicySnippets(t *testing.T) {
	t.Skip("Skipped while investigating https://github.com/GoogleCloudPlatform/golang-samples/issues/2811")
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	policiesClient, err := iam.NewPoliciesClient(ctx)
	if err != nil {
		t.Fatalf("NewPoliciesClient: %v", err)
	}
	defer policiesClient.Close()

	buf := &bytes.Buffer{}
	policyID := fmt.Sprintf("test-policy-%v-%v", time.Now().Format("01-02-2006"), r.Int())
	policyName := fmt.Sprintf(
		"policies/cloudresourcemanager.googleapis.com%%2Fprojects%%2F949737848314/denypolicies/%s",
		policyID,
	)

	if err := createDenyPolicy(buf, tc.ProjectID, policyID); err != nil {
		t.Fatalf("createDenyPolicy: %v", err)
	}

	expectedResult := fmt.Sprintf("Policy %s created", policyName)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Fatalf("createDenyPolicy got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	if err := listDenyPolicies(buf, tc.ProjectID); err != nil {
		t.Fatalf("listDenyPolicies: %v", err)
	}

	expectedResult = fmt.Sprintf("- %s", policyName)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Fatalf("listDenyPolicies got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	if err := getDenyPolicy(buf, tc.ProjectID, policyID); err != nil {
		t.Fatalf("getDenyPolicy: %v", err)
	}

	expectedResult = fmt.Sprintf("Policy %s retrieved", policyName)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Fatalf("getDenyPolicy got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	req := &iampb.GetPolicyRequest{
		Name: policyName,
	}
	policy, err := policiesClient.GetPolicy(ctx, req)
	if err != nil {
		t.Fatalf("unable to get policy: %v", err)
	}

	if err := updateDenyPolicy(buf, tc.ProjectID, policyID, policy.GetEtag()); err != nil {
		t.Fatalf("updateDenyPolicy: %v", err)
	}

	expectedResult = fmt.Sprintf("Policy %s updated", policyName)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Fatalf("getDenyPolicy got %q, want %q", got, expectedResult)
	}

	buf.Reset()

	if err := deleteDenyPolicy(buf, tc.ProjectID, policyID); err != nil {
		t.Fatalf("deleteDenyPolicy: %v", err)
	}

	expectedResult = fmt.Sprintf("Policy %s deleted", policyName)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Fatalf("deleteDenyPolicy got %q, want %q", got, expectedResult)
	}
}
