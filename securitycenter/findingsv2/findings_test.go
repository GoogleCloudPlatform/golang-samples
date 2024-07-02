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

package findingsv2

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	"cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/golang/protobuf/ptypes"
)

var orgID = ""
var sourceName = ""
var findingName = ""
var untouchedFindingName = ""

func createTestFinding(ctx context.Context, client *securitycenter.Client, findingID string, category string) (*securitycenterpb.Finding, error) {
	eventTime, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, fmt.Errorf("TimestampProto: %w", err)
	}

	req := &securitycenterpb.CreateFindingRequest{
		Parent:    sourceName,
		FindingId: findingID,
		Finding: &securitycenterpb.Finding{
			State: securitycenterpb.Finding_ACTIVE,
			// Resource the finding is associated with.  This is an
			// example any resource identifier can be used.
			ResourceName: "//cloudresourcemanager.googleapis.com/organizations/11232/locations/global",
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
		return fmt.Errorf("securitycenter.NewClient: %w", err)
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
		return fmt.Errorf("CreateSource: %w", err)
	}
	sourceName = source.Name
	finding, err := createTestFinding(ctx, client, "updated", "MEDIUM_RISK_ONE")
	if err != nil {
		return fmt.Errorf("createTestFinding: %w", err)
	}
	findingName = finding.Name
	finding, err = createTestFinding(ctx, client, "untouched", "XSS")
	if err != nil {
		return fmt.Errorf("createTestFinding: %w", err)
	}
	untouchedFindingName = finding.Name
	return nil
}

func setup(t *testing.T) string {
	if orgID == "" {
		t.Skip("GCLOUD_ORGANIZATION not set")
	}
	if sourceName == "" {
		t.Errorf("sourceName not set")
		os.Exit(1)
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

func TestCreateFinding(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		err := createFinding(buf, sourceName)

		if err != nil {
			r.Errorf("createFinding(%s) had error: %v", sourceName, err)
			return
		}

		got := buf.String()
		if want := fmt.Sprintf("%s/findings/samplefindingid", sourceName); !strings.Contains(got, want) {
			r.Errorf("createFinding(%s) got: %s want %s", sourceName, got, want)
		}
	})
}

func TestUpdateFindingSourceProperties(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		err := updateFindingSourceProperties(buf, findingName)

		if err != nil {
			r.Errorf("updateFindingSourceProperties(%s) had error: %v", findingName, err)
			return
		}

		got := buf.String()
		if want := "s_value"; !strings.Contains(got, want) {
			r.Errorf("updateFindingSourceProperties(%s) got: %s want %s", findingName, got, want)
		}
		if !strings.Contains(got, findingName) {
			r.Errorf("updateFindingSourceProperties(%s) got: %s want %s", findingName, got, findingName)
		}
	})
}

func TestSetFindingState(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		err := setFindingState(buf, findingName)

		if err != nil {
			r.Errorf("setFindingState(%s) had error: %v", findingName, err)
			return
		}

		got := buf.String()
		if want := "INACTIVE"; !strings.Contains(got, want) {
			r.Errorf("setFindingState(%s) got: %s want %s", findingName, got, want)
		}
		if !strings.Contains(got, findingName) {
			r.Errorf("setFindingState(%s) got: %s want %s", findingName, got, findingName)
		}
	})
}
