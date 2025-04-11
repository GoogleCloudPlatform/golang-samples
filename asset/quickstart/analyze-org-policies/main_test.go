// Copyright 2023 Google LLC
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

package create

import (
	"bytes"
	"context"
	"strings"
	"testing"

	asset "cloud.google.com/go/asset/apiv1"
)

var (
	projectID     string
	savedQueryID  string
	projectNumber string

	ctx    context.Context
	client *asset.Client
)

const (
	// organizations/474566717491 is ipa1.joonix.net, a test organization owned by mdb.cloud-asset-application-team@google.com.
	organization    = "organizations/474566717491"
	constraint      = "constraints/iam.allowServiceAccountCredentialLifetimeExtension"
	expected_result = "consolidated_policy"
)

func TestAnalyzeOrgPolicies(t *testing.T) {
	buf := new(bytes.Buffer)
	err := analyzeOrgPolicies(buf, organization, constraint)
	if err != nil {
		t.Errorf("analyzeOrgPolicies: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, expected_result) {
		t.Errorf("analyzeOrgPolicies got%q, want%q", got, expected_result)
	}
}
