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

package muteconfigv2

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"io"
	"log"
	"os"
	"strings"
	"testing"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	"cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
)

var fixture *muteconfigFixture

func TestMain(m *testing.M) {
	var err error
	fixture, err = newMuteConfigFixture()
	if err != nil {
		log.Fatalf("failed to create fixture: %s", err)
	}

	exitCode := m.Run()

	if err := fixture.Cleanup(); err != nil {
		log.Fatalf("failed to cleanup resources: %s", err)
	}
	fixture.client.Close()
	os.Exit(exitCode)
}

func createFinding(findingID string, category string, sourceName string, orgId string) (*securitycenterpb.Finding, error) {
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
			ResourceName: fmt.Sprintf("//cloudresourcemanager.googleapis.com/organizations/%s", orgId),
			// A free-form category.
			Category: category,
			// The time associated with discovering the issue.
			EventTime: eventTime,
			Severity:  securitycenterpb.Finding_LOW,
		},
	}
	return client.CreateFinding(ctx, req)
}

func createSource(w io.Writer, orgId string) error {
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	source, err := client.CreateSource(ctx, &securitycenterpb.CreateSourceRequest{
		Source: &securitycenterpb.Source{
			DisplayName: "Customized Display Name",
			Description: "A new custom source that does X",
		},
		Parent: fmt.Sprintf("organizations/%s", orgId),
	})

	if err != nil {
		return fmt.Errorf("CreateSource: %w", err)
	}
	fmt.Fprintf(w, source.Name)
	return nil
}

type muteconfigFixture struct {
	client        *securitycenter.Client
	orgId         string
	projectId     string
	locationId    string
	parent        string
	sourceName    string
	finding1Name  string
	finding2Name  string
	muteConfigId1 string
	muteConfigId2 string
}

// newMuteConfigFixture initializes variables in this file with resources to
// use for testing.
func newMuteConfigFixture() (*muteconfigFixture, error) {
	var mc muteconfigFixture
	var err error

	mc.client, err = securitycenter.NewClient(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create security center client: %w", err)
	}

	orgId := os.Getenv("SCC_PROJECT_ORG_ID")
	projectId := os.Getenv("SCC_PROJECT_ID")
	locationId := "global"
	if orgId == "" || projectId == "" {
		return nil, fmt.Errorf("env variables not set: %v", err)
	}
	mc.orgId = orgId
	mc.projectId = projectId
	mc.locationId = locationId
	// Create source.
	buf := &bytes.Buffer{}
	if err := createSource(buf, orgId); err != nil {
		return nil, fmt.Errorf("createSource: %w", err)
	}
	sourceName := buf.String()
	mc.sourceName = sourceName
	buf.Reset()

	// Create findings.
	finding, err := createFinding("updated", "MEDIUM_RISK_ONE", sourceName, orgId)
	if err != nil {
		return nil, fmt.Errorf("createTestFinding: %w", err)
	}
	mc.finding1Name = finding.Name
	finding, err = createFinding("untouched", "XSS", sourceName, orgId)
	if err != nil {
		return nil, fmt.Errorf("createTestFinding: %w", err)
	}
	mc.finding2Name = finding.Name

	// Create mute rules.
	parent := fmt.Sprintf("projects/%s/locations/%s", projectId, locationId)
	mc.parent = parent
	muteConfigId1 := "random-mute-id-" + uuid.New().String()
	muteConfigId2 := "random-mute-id-" + uuid.New().String()
	err = createMuteRule(buf, parent, muteConfigId1)
	if err != nil {
		return nil, fmt.Errorf("createMuteRule failed: %w", err)
	}
	err = createMuteRule(buf, parent, muteConfigId2)
	if err != nil {
		return nil, fmt.Errorf("createMuteRule failed: %w", err)
	}
	mc.muteConfigId1 = muteConfigId1
	mc.muteConfigId2 = muteConfigId2

	return &mc, nil
}

// Cleanup deletes any resources.
func (mc *muteconfigFixture) Cleanup() error {
	// Delete mute rules.
	buf := &bytes.Buffer{}
	if err := deleteMuteRule(buf, mc.parent, mc.muteConfigId1); err != nil {
		return fmt.Errorf("deleteMuteRule had error: %w", err)
	}
	if got := buf.String(); !strings.Contains(got, "Mute rule deleted successfully") {
		return fmt.Errorf("deleteMuteRule got %q, expected %q", got, mc.muteConfigId1)
	}
	buf.Reset()
	if err := deleteMuteRule(buf, mc.parent, mc.muteConfigId2); err != nil {
		return fmt.Errorf("deleteMuteRule had error: %w", err)
	}
	if got := buf.String(); !strings.Contains(got, "Mute rule deleted successfully") {
		return fmt.Errorf("deleteMuteRule got %q, expected %q", got, mc.muteConfigId2)
	}
	return nil
}

func TestListMuteRules(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	// List mute rules.
	if err := listMuteRules(&buf, fixture.parent); err != nil {
		t.Errorf("listMuteRules had error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, fixture.muteConfigId1) {
		t.Errorf("listMuteRules got %q, expected %q", got, fixture.muteConfigId1)
	}
	if !strings.Contains(got, fixture.muteConfigId2) {
		t.Errorf("listMuteRules got %q, expected %q", got, fixture.muteConfigId2)
	}
}

func TestGetMuteRule(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	// Get mute rule.
	if err := getMuteRule(&buf, fixture.parent, fixture.muteConfigId1); err != nil {
		t.Errorf("getMuteRule had error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, fixture.muteConfigId1) {
		t.Errorf("getMuteRule got %q, expected %q", got, fixture.muteConfigId1)
	}
}

func TestUpdateMuteRule(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	// Update a mute rule.
	muteConfigName := fmt.Sprintf("%s/muteConfigs/%s", fixture.parent, fixture.muteConfigId1)
	if err := updateMuteRule(&buf, muteConfigName); err != nil {
		t.Errorf("updateMuteRule had error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "Mute rule updated") {
		t.Errorf("updateMuteRule got %q, expected %q", got, "Mute rule updated")
	}
}

func TestSetMuteFinding(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	// Get mute rule.
	if err := getMuteRule(&buf, fixture.parent, fixture.muteConfigId1); err != nil {
		t.Errorf("getMuteRule had error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, fixture.muteConfigId1) {
		t.Errorf("getMuteRule got %q, expected %q", got, fixture.muteConfigId1)
	}
}

func TestSetUnmuteFinding(t *testing.T) {
	t.Skip("see https://github.com/GoogleCloudPlatform/golang-samples/issues/3793")
	testutil.SystemTest(t)

	var buf bytes.Buffer
	// Unmute an individual finding.
	if err := setUnmute(&buf, fixture.finding1Name); err != nil {
		t.Errorf("setUnmute had error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, fmt.Sprintf("Mute value for the finding: %s is %s", fixture.finding1Name, "UNMUTE")) {
		t.Errorf("setUnmute got %q, expected %q", got, fmt.Sprintf("Mute value for the finding: %s is %s", fixture.finding1Name, "UNMUTE"))
	}
	if err := setUnmute(&buf, fixture.finding2Name); err != nil {
		t.Errorf("setUnmute had error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, fmt.Sprintf("Mute value for the finding: %s is %s", fixture.finding1Name, "UNMUTE")) {
		t.Errorf("setUnmute got %q, expected %q", got, fmt.Sprintf("Mute value for the finding: %s is %s", fixture.finding1Name, "UNMUTE"))
	}
}

func TestBulkMuteFinding(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	// Bulk mute findings.
	if err := bulkMute(&buf, fixture.parent, "severity=\"LOW\""); err != nil {
		t.Errorf("bulkMute had error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "Bulk mute findings completed successfully") {
		t.Errorf("bulkMute got %q, expected %q", got, "Bulk mute findings completed successfully")
	}
}
