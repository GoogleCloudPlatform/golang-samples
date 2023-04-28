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

func TestAnalyzeOrgPolicies(t *testing.T) {
	buf := new(bytes.Buffer)
	err := analyzeOrgPolicies(buf, "organizations/474566717491", "constraints/iam.allowServiceAccountCredentialLifetimeExtension")
	if err != nil {
		t.Errorf("analyzeOrgPolicies: %v", err)
	}
	got := buf.String()
	if want := "consolidated_policy"; !strings.Contains(got, want) {
		t.Errorf("analyzeOrgPolicies got%q, want%q", got, want)
	}
}
