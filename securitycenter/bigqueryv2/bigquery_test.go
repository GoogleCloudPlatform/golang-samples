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

package bigqueryv2

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
	"github.com/google/uuid"
)

var projectID = ""

// setup initializes variables in this file with entityNames to
// use for testing.
func setup(t *testing.T) string {
	projectID = os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	return projectID
}

func createDemoDataset(t *testing.T) string {
	return "projects/project-a-id/datasets/sampledataset"
}

func addBigQueryExport(t *testing.T, bigQueryExportID string) error {
	projectID := setup(t)

	bigQueryDatasetName := createDemoDataset(t)

	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	bigQueryExport := &securitycenterpb.BigQueryExport{
		Description: "BigQueryExport that receives all HIGH severity Findings",
		Filter:      "severity=\"HIGH\"",
		Dataset:     bigQueryDatasetName,
	}

	req := &securitycenterpb.CreateBigQueryExportRequest{
		Parent:           fmt.Sprintf("projects/%s/locations/global", projectID),
		BigQueryExport:   bigQueryExport,
		BigQueryExportId: bigQueryExportID,
	}

	_, err0 := client.CreateBigQueryExport(ctx, req)
	if err0 != nil {
		return fmt.Errorf("Failed to create BigQueryConfig: %w", err0)
	}
	return nil
}

func cleanupBigQueryExport(t *testing.T, bigQueryExportID string) error {
	projectID := setup(t)

	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycenterpb.DeleteBigQueryExportRequest{
		Name: fmt.Sprintf("projects/%s/locations/global/bigQueryExports/%s", projectID, bigQueryExportID),
	}

	if err := client.DeleteBigQueryExport(ctx, req); err != nil {
		return fmt.Errorf("failed to delete BigQueryExport: %w", err)
	}

	return nil
}

func TestListBigQuery(t *testing.T) {
	projectID := setup(t)

	buf := new(bytes.Buffer)

	// Create Test BigQueryExport Config
	rand, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("Issue generating id.")
		return
	}
	configID := "random-bqexport-id-" + rand.String()

	if err := addBigQueryExport(t, configID); err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
		return
	}

	parent := fmt.Sprintf("projects/%s/locations/global", projectID)

	// Call List BigQueryExport
	err = listBigQueryExport(buf, parent)

	if err != nil {
		t.Fatalf("listBigQueryExport() had error: %v", err)
		return
	}

	got := buf.String()

	if !strings.Contains(got, configID) {
		t.Fatalf("listBigQueryConfigs() got: %s want %s", got, configID)
	}

	// Cleanup
	cleanupBigQueryExport(t, configID)
}

func TestGetBigQuery(t *testing.T) {
	projectID := setup(t)

	buf := new(bytes.Buffer)
	// Create Test BigQueryExport Config
	rand, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("Issue generating id.")
		return
	}
	configID := "random-bqexport-id-" + rand.String()

	if err := addBigQueryExport(t, configID); err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
		return
	}

	parent := fmt.Sprintf("projects/%s/locations/global", projectID)

	// Call GetBigQueryExport
	err = getBigQueryExport(buf, parent, configID)

	if err != nil {
		t.Fatalf("getBigQueryExport() had error: %v", err)
		return
	}

	got := buf.String()

	if !strings.Contains(got, configID) {
		t.Fatalf("getBigQueryExport() got: %s want %s", got, configID)
	}

	// Cleanup
	cleanupBigQueryExport(t, configID)
}

func TestDeleteBigQuery(t *testing.T) {
	projectID := setup(t)

	buf := new(bytes.Buffer)
	// Create Test BigQueryExport Config
	rand, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("Issue generating id.")
		return
	}
	configID := "random-bqexport-id-" + rand.String()

	if err := addBigQueryExport(t, configID); err != nil {
		t.Fatalf("Could not setup test environment: %v", err)
		return
	}

	parent := fmt.Sprintf("projects/%s/locations/global", projectID)

	// Call DeleteBigQueryExport
	err = deleteBigQueryExport(buf, parent, configID)

	if err != nil {
		t.Fatalf("deleteBigQueryExport() had error: %v", err)
		return
	}

	got := buf.String()

	if !strings.Contains(got, configID) {
		t.Fatalf("deleteBigQueryExport() got: %s want %s", got, configID)
	}
}

func TestUpdateBigQuery(t *testing.T) {
	projectID := setup(t)

	testutil.Retry(t, 3, 10*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		// Create Test BigQueryExport Config
		rand, err := uuid.NewUUID()
		if err != nil {
			t.Fatalf("Issue generating id.")
			return
		}
		configID := "random-bqexport-id-" + rand.String()

		if err := addBigQueryExport(t, configID); err != nil {
			t.Fatalf("Could not setup test environment: %v", err)
			return
		}

		parent := fmt.Sprintf("projects/%s/locations/global", projectID)

		// Call UpdateBigQueryExport
		err = updateBigQueryExport(buf, parent, configID)

		if err != nil {
			t.Fatalf("updateBigQueryExport() had error: %v", err)
			return
		}

		got := buf.String()

		if !strings.Contains(got, configID) {
			t.Fatalf("updateBigQueryExport() got: %s want %s", got, configID)
		}

		// Cleanup
		cleanupBigQueryExport(t, configID)
	})
}

func TestCreateBigQuery(t *testing.T) {
	projectID := setup(t)

	buf := new(bytes.Buffer)

	rand, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("Issue generating id.")
		return
	}
	configID := "random-bqexport-id-" + rand.String()

	parent := fmt.Sprintf("projects/%s/locations/global", projectID)

	err = createBigQueryExport(buf, parent, configID, projectID)

	if err != nil {
		t.Fatalf("createBigQueryExport() had error: %v", err)
		return
	}

	got := buf.String()

	if !strings.Contains(got, configID) {
		t.Fatalf("createBigQueryExport() got: %s want %s", got, configID)
	}

	// Cleanup
	cleanupBigQueryExport(t, configID)

}
