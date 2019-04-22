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

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"github.com/golang/protobuf/ptypes"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1"
)

var orgID = ""
var sourceName = ""
var findingName = ""
var untouchedFindingName = ""

func createTestFinding(ctx context.Context, client *securitycenter.Client, findingID string, category string) (*securitycenterpb.Finding, error) {
	eventTime, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, fmt.Errorf("TimestampProto: %v", err)
	}

	req := &securitycenterpb.CreateFindingRequest{
		Parent:    sourceName,
		FindingId: findingID,
		Finding: &securitycenterpb.Finding{
			State: securitycenterpb.Finding_ACTIVE,
			// Resource the finding is associated with.  This is an
			// example any resource identifier can be used.
			ResourceName: "//cloudresourcemanager.googleapis.com/organizations/11232",
			// A free-form category.
			Category: category,
			// The time associated with discovering the issue.
			EventTime: eventTime,
		},
	}
	return client.CreateFinding(ctx, req)
}

// setupEntities initializes variables in this file with entityNames to
// use for testing.
func setupEntities() error {
	orgID = os.Getenv("GCLOUD_ORGANIZATION")
	if orgID == "" {
		// Each test checks for GCLOUD_ORGANIZATION. Return nil so we see every skip.
		return nil
	}
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %v", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	source, err := client.CreateSource(ctx, &securitycenterpb.CreateSourceRequest{
		Source: &securitycenterpb.Source{
			DisplayName: "Customized Display Name",
			Description: "A new custom source that does X",
		},
		Parent: fmt.Sprintf("organizations/%s", orgID),
	})

	if err != nil {
		return fmt.Errorf("CreateSource: %v", err)
	}
	sourceName = source.Name
	finding, err := createTestFinding(ctx, client, "updated", "MEDIUM_RISK_ONE")
	if err != nil {
		return fmt.Errorf("createTestFinding: %v", err)
	}
	findingName = finding.Name
	finding, err = createTestFinding(ctx, client, "untouched", "XSS")
	if err != nil {
		return fmt.Errorf("createTestFinding: %v", err)
	}
	untouchedFindingName = finding.Name
	return nil
}

func setup(t *testing.T) string {
	if orgID == "" {
		t.Skip("GCLOUD_ORGANIZATION not set")
	}
	if sourceName == "" {
		t.Fatalf("sourceName not set")
	}
	return orgID
}

func TestMain(m *testing.M) {
	if err := setupEntities(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialize findings test environment: %v\n", err)
		return
	}
	code := m.Run()
	os.Exit(code)
}

func TestCreateSource(t *testing.T) {
	orgID := setup(t)
	buf := new(bytes.Buffer)

	err := createSource(buf, orgID)

	if err != nil {
		t.Fatalf("createSource(%s) had error: %v", orgID, err)
	}

	got := buf.String()
	if want := "New source created"; !strings.Contains(got, want) {
		t.Errorf("createSource(%s) got: %s want %s", orgID, got, want)
	}
	if !strings.Contains(got, orgID) {
		t.Errorf("createSource(%s) got: %s want %s", orgID, got, orgID)
	}
}

func TestGetSource(t *testing.T) {
	setup(t)
	buf := new(bytes.Buffer)

	err := getSource(buf, sourceName)

	if err != nil {
		t.Fatalf("getSource(%s) had error: %v", sourceName, err)
	}

	got := buf.String()
	if want := "Description: A new custom source that does X"; !strings.Contains(got, want) {
		t.Errorf("getSource(%s) got: %s want %s", sourceName, got, want)
	}
	if !strings.Contains(got, orgID) {
		t.Errorf("getSource(%s) got: %s want %s", sourceName, got, sourceName)
	}
}

func TestListSources(t *testing.T) {
	setup(t)
	buf := new(bytes.Buffer)

	err := listSources(buf, orgID)

	if err != nil {
		t.Fatalf("listSource(%s) had error: %v", orgID, err)
	}

	got := buf.String()
	if want := "Description: A new custom source that does X"; !strings.Contains(got, want) {
		t.Errorf("listSource(%s) got: %s want %s", orgID, got, want)
	}
	if !strings.Contains(got, sourceName) {
		t.Errorf("listSource(%s) got: %s want %s", orgID, got, sourceName)
	}
}

func TestUpdateSource(t *testing.T) {
	setup(t)
	buf := new(bytes.Buffer)

	err := updateSource(buf, sourceName)

	if err != nil {
		t.Fatalf("updateSource(%s) had error: %v", sourceName, err)
	}

	got := buf.String()
	if want := "Display name: New Display Name"; !strings.Contains(got, want) {
		t.Errorf("updateSource(%s) got: %s want %s", sourceName, got, want)
	}
	if !strings.Contains(got, sourceName) {
		t.Errorf("updateSource(%s) got: %s want %s", sourceName, got, sourceName)
	}
	if want := "does X"; !strings.Contains(got, want) {
		t.Errorf("updateSource(%s) got: %s want %s", sourceName, got, want)
	}

}

func TestCreateFinding(t *testing.T) {
	setup(t)
	buf := new(bytes.Buffer)

	err := createFinding(buf, sourceName)

	if err != nil {
		t.Fatalf("createFinding(%s) had error: %v", sourceName, err)
	}

	got := buf.String()
	if want := fmt.Sprintf("%s/findings/samplefindingid", sourceName); !strings.Contains(got, want) {
		t.Errorf("createFinding(%s) got: %s want %s", sourceName, got, want)
	}
}

func TestCreateFindingWithProperties(t *testing.T) {
	setup(t)
	buf := new(bytes.Buffer)

	err := createFindingWithProperties(buf, sourceName)

	if err != nil {
		t.Fatalf("createFindingWithProperties(%s) had error: %v", sourceName, err)
	}

	got := buf.String()
	if want := "s_value"; !strings.Contains(got, want) {
		t.Errorf("createFindingWithProperties(%s) got: %s want %s", sourceName, got, want)
	}
	if want := "n_value"; !strings.Contains(got, want) {
		t.Errorf("createFindingWithProperties(%s) got: %s want %s", sourceName, got, want)
	}

	if want := fmt.Sprintf("%s/findings/samplefindingprops", sourceName); !strings.Contains(got, want) {
		t.Errorf("createFindingWithProperties(%s) got: %s want %s", sourceName, got, want)
	}
}

func TestUpdateFindingSourceProperties(t *testing.T) {
	setup(t)
	buf := new(bytes.Buffer)

	err := updateFindingSourceProperties(buf, findingName)

	if err != nil {
		t.Fatalf("updateFindingSourceProperties(%s) had error: %v", findingName, err)
	}

	got := buf.String()
	if want := "s_value"; !strings.Contains(got, want) {
		t.Errorf("updateFindingSourceProperties(%s) got: %s want %s", findingName, got, want)
	}
	if !strings.Contains(got, findingName) {
		t.Errorf("updateFindingSourceProperties(%s) got: %s want %s", findingName, got, findingName)
	}
}

func TestSetFindingState(t *testing.T) {
	setup(t)
	buf := new(bytes.Buffer)

	err := setFindingState(buf, findingName)

	if err != nil {
		t.Fatalf("setFindingState(%s) had error: %v", findingName, err)
	}

	got := buf.String()
	if want := "INACTIVE"; !strings.Contains(got, want) {
		t.Errorf("setFindingState(%s) got: %s want %s", findingName, got, want)
	}
	if !strings.Contains(got, findingName) {
		t.Errorf("setFindingState(%s) got: %s want %s", findingName, got, findingName)
	}
}

func TestTestIam(t *testing.T) {
	setup(t)
	buf := new(bytes.Buffer)

	err := testIam(buf, sourceName)

	if err != nil {
		t.Fatalf("testIam(%s) had error: %v", sourceName, err)
	}

	got := buf.String()
	trueCount := strings.Count(got, "true")
	if want := 2; want != trueCount {
		t.Errorf("testIam(%s) got %q, true count: %d, want %d", sourceName, got, trueCount, want)
	}
}

func TestListAllFindings(t *testing.T) {
	orgID := setup(t)
	buf := new(bytes.Buffer)

	err := listFindings(buf, orgID)

	if err != nil {
		t.Fatalf("listFindings(%s) had error: %v", orgID, err)
	}

	got := buf.String()
	if !strings.Contains(got, findingName) {
		t.Errorf("listFindings(%s) got: %s want %s", orgID, got, findingName)
	}

	if !strings.Contains(got, untouchedFindingName) {
		t.Errorf("listFindings(%s) got: %s want %s", orgID, got, untouchedFindingName)
	}
}

func TestListFilteredFindings(t *testing.T) {
	setup(t)
	buf := new(bytes.Buffer)

	err := listFilteredFindings(buf, sourceName)

	if err != nil {
		t.Fatalf("listFilteredFindings(%s) had error: %v", sourceName, err)
	}

	got := buf.String()
	if !strings.Contains(got, findingName) {
		t.Errorf("listFilteredFindings(%s) got: %s want %s", sourceName, got, findingName)
	}

	if strings.Contains(got, untouchedFindingName) {
		t.Errorf("listFilteredFindings(%s) got: %s didn't want %s", sourceName, got, untouchedFindingName)
	}
}

func TestListFindingsAtTime(t *testing.T) {
	setup(t)
	buf := new(bytes.Buffer)

	err := listFindingsAtTime(buf, sourceName)

	if err != nil {
		t.Fatalf("listFindingsAtTime(%s) had error: %v", sourceName, err)
	}

	got := buf.String()
	// No findings 5 days ago.
	if got != "" {
		t.Errorf("listFindingsAtTime(%s) got: %s wanted empty", sourceName, got)
	}

}

func TestAddSecurityMarks(t *testing.T) {
	setup(t)
	buf := new(bytes.Buffer)

	err := addSecurityMarks(buf, findingName)

	if err != nil {
		t.Fatalf("addSecurityMarks(%s) adding marks had error: %v", findingName, err)
	}

	got := buf.String()
	if want := "key_a = value_a"; !strings.Contains(got, want) {
		t.Errorf("addSecurityMarks(%s) got: %s want %s", findingName, got, want)
	}

	if want := "key_b = value_b"; !strings.Contains(got, want) {
		t.Errorf("addSecurityMarks(%s) got: %s want %s", findingName, got, want)
	}
}

func TestListFindingsWithMarks(t *testing.T) {
	orgID := setup(t)
	buf := new(bytes.Buffer)
	// Ensure security marks have been added so filter is effective.
	err := addSecurityMarks(buf, findingName)
	buf.Truncate(0)
	if err != nil {
		t.Fatalf("listFindingsWithMark(%s) adding marks had error: %v", findingName, err)
	}

	err = listFindingsWithMarks(buf, sourceName)

	if err != nil {
		t.Fatalf("listFindingsWithMark(%s) had error: %v", sourceName, err)
	}

	got := buf.String()
	if !strings.Contains(got, untouchedFindingName) {
		t.Errorf("listFindingWithMarks(%s) got: %s want %s", sourceName, got, untouchedFindingName)
	}

	if strings.Contains(got, findingName) {
		t.Errorf("listFindingWithMarks(%s) got: %s didn't want %s", orgID, got, findingName)
	}

}

func TestGetSourceIamPolicy(t *testing.T) {
	setup(t)
	buf := new(bytes.Buffer)

	err := getSourceIamPolicy(buf, sourceName)

	if err != nil {
		t.Fatalf("getSourceIamPolicy(%s) had error: %v", sourceName, err)
	}

	got := buf.String()
	if want := "Policy"; !strings.Contains(got, want) {
		t.Errorf("getSourceIamPolicy(%s) got: %s want %s", sourceName, got, want)
	}
}

func TestSetSourceIamPolicy(t *testing.T) {
	setup(t)
	buf := new(bytes.Buffer)

	user := "csccclienttest@gmail.com"
	err := setSourceIamPolicy(buf, sourceName, user)

	if err != nil {
		t.Fatalf("setSourceIamPolicy(%s, %s) had error: %v", sourceName, user, err)
	}

	got := buf.String()
	if want := "Principal: user:csccclienttest@gmail.com Role: roles/securitycenter.findingsEditor"; !strings.Contains(got, want) {
		t.Errorf("setSourceIamPolicy(%s, %s) got: %s want %s", sourceName, user, got, want)
	}
}
