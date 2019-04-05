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

func createFindingHelper(ctx context.Context, client *securitycenter.Client, findingID string) (*securitycenterpb.Finding, error) {
	eventTime, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, fmt.Errorf("Error converting now: %v", err)
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
			Category: "MEDIUM_RISK_ONE",
			// The time associated with discovering the issue.
			EventTime: eventTime,
		},
	}
	return client.CreateFinding(ctx, req)
}

// setupEntities initializes variables in this file with entityNames to
// use for testing.
func setupEntities() {
	orgID = os.Getenv("GCLOUD_ORGANIZATION")
	if orgID == "" {
		return
	}
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		fmt.Printf("Error instantiating client %v\n", err)
		return
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
		fmt.Printf("Error creating source: %v", err)
		return
	}
	sourceName = source.Name
	finding, err := createFindingHelper(ctx, client, "finding1")
	if err != nil {
		fmt.Printf("Error creating findings1 %v", err)
		return
	}
	findingName = finding.Name
	finding, err = createFindingHelper(ctx, client, "finding2")
	if err != nil {
		fmt.Printf("Error creating findings1 %v", err)
		return
	}
	untouchedFindingName = finding.Name

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
	setupEntities()
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
