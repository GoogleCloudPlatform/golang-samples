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

package muteconfig

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"cloud.google.com/go/securitycenter/apiv1/securitycenterpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
)

var orgID = os.Getenv("SCC_PROJECT_ORG_ID")
var projectID = os.Getenv("SCC_PROJECT_ID")
var sourceName = ""
var finding1Name = ""
var finding2Name = ""

func createFinding(findingID string, category string, sourceName string) (*securitycenterpb.Finding, error) {
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("securitycenter.NewClient: %v", err)
	}

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
			ResourceName: fmt.Sprintf("//cloudresourcemanager.googleapis.com/organizations/%s", orgID),
			// A free-form category.
			Category: category,
			// The time associated with discovering the issue.
			EventTime: eventTime,
			Severity:  securitycenterpb.Finding_LOW,
		},
	}
	return client.CreateFinding(ctx, req)
}

func createSource(w io.Writer) error {
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %v", err)
	}
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
	fmt.Fprintf(w, source.Name)
	return nil
}

// setupResources initializes variables in this file with resources to
// use for testing.
func setupResources() error {
	orgID = os.Getenv("SCC_PROJECT_ORG_ID")
	projectID = os.Getenv("SCC_PROJECT_ID")
	if orgID == "" || projectID == "" {
		return nil
	}

	// Create source.
	buf := &bytes.Buffer{}
	if err := createSource(buf); err != nil {
		return fmt.Errorf("createSource: %v", err)
	}
	sourceName = buf.String()
	buf.Reset()

	// Create findings.
	finding, err := createFinding("updated", "MEDIUM_RISK_ONE", sourceName)
	if err != nil {
		return fmt.Errorf("createTestFinding: %v", err)
	}
	finding1Name = finding.Name
	finding, err = createFinding("untouched", "XSS", sourceName)
	if err != nil {
		return fmt.Errorf("createTestFinding: %v", err)
	}
	finding2Name = finding.Name
	return nil
}

func TestMuteConfigCRUD(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	if err := setupResources(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialize mute config test environment: %v\n", err)
		return
	}

	// Create mute rules.
	parent := fmt.Sprintf("projects/%s", projectID)
	muteConfigId1 := "random-mute-id-" + uuid.New().String()
	muteConfigId2 := "random-mute-id-" + uuid.New().String()
	err := createMuteRule(buf, parent, muteConfigId1)
	if err != nil {
		t.Errorf("createMuteRule failed: %v", err)
	}
	err = createMuteRule(buf, parent, muteConfigId2)
	if err != nil {
		t.Errorf("createMuteRule failed: %v", err)
	}

	buf.Reset()

	// List mute rules.
	if err := listMuteRules(buf, parent); err != nil {
		t.Errorf("listMuteRules had error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, muteConfigId1) {
		t.Errorf("listMuteRules got %q, expected %q", got, muteConfigId1)
	}

	buf.Reset()

	// Get mute rule.
	if err := getMuteRule(buf, parent, muteConfigId1); err != nil {
		t.Errorf("getMuteRule had error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, muteConfigId1) {
		t.Errorf("getMuteRule got %q, expected %q", got, muteConfigId1)
	}

	buf.Reset()

	// Mute an individual finding.
	if err := setMute(buf, finding1Name); err != nil {
		t.Errorf("setMute had error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, fmt.Sprintf("Mute value for the finding: %s is %s", finding1Name, "MUTE")) {
		t.Errorf("setMute got %q, expected %q", got, fmt.Sprintf("Mute value for the finding: %s is %s", finding1Name, "MUTE"))
	}

	buf.Reset()

	// Unmute an individual finding.
	if err := setUnmute(buf, finding1Name); err != nil {
		t.Errorf("setUnmute had error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, fmt.Sprintf("Mute value for the finding: %s is %s", finding1Name, "UNMUTE")) {
		t.Errorf("setUnmute got %q, expected %q", got, fmt.Sprintf("Mute value for the finding: %s is %s", finding1Name, "UNMUTE"))
	}
	if err := setUnmute(buf, finding2Name); err != nil {
		t.Errorf("setUnmute had error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, fmt.Sprintf("Mute value for the finding: %s is %s", finding1Name, "UNMUTE")) {
		t.Errorf("setUnmute got %q, expected %q", got, fmt.Sprintf("Mute value for the finding: %s is %s", finding1Name, "UNMUTE"))
	}

	buf.Reset()

	// Update a mute rule.
	muteConfigName := fmt.Sprintf("%s/muteConfigs/%s", parent, muteConfigId1)
	if err := updateMuteRule(buf, muteConfigName); err != nil {
		t.Errorf("updateMuteRule had error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "Mute rule updated") {
		t.Errorf("updateMuteRule got %q, expected %q", got, "Mute rule updated")
	}

	buf.Reset()

	// Bulk mute findings.
	if err := bulkMute(buf, parent, "severity=\"LOW\""); err != nil {
		t.Errorf("bulkMute had error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "Bulk mute findings completed successfully") {
		t.Errorf("bulkMute got %q, expected %q", got, "Bulk mute findings completed successfully")
	}

	buf.Reset()

	// Delete mute rules.
	if err := deleteMuteRule(buf, parent, muteConfigId1); err != nil {
		t.Errorf("deleteMuteRule had error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "Mute rule deleted successfully") {
		t.Errorf("deleteMuteRule got %q, expected %q", got, muteConfigId1)
	}
	if err := deleteMuteRule(buf, parent, muteConfigId2); err != nil {
		t.Errorf("deleteMuteRule had error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "Mute rule deleted successfully") {
		t.Errorf("deleteMuteRule got %q, expected %q", got, muteConfigId2)
	}

}
