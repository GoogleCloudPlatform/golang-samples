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

func orgID(t *testing.T) string {
	orgID := os.Getenv("GCLOUD_ORGANIZATION")
	if orgID == "" {
		t.Skip("GCLOUD_ORGANIZATION not set")
	}
	return orgID
}

func bigQueryDatasetName(t *testing.T) string {
	bigQueryDatasetName := os.Getenv("SCC_BIGQUERY_DATASET_NAME")
	if bigQueryDatasetName == "" {
		t.Skip("SCC_BIGQUERY_DATASET_NAME not set")
	}
	return bigQueryDatasetName
}

func addBigQueryExport(t *testing.T, bigQueryExportID string) error {
	orgID := orgID(t)
	bigQueryDatasetName := bigQueryDatasetName(t)

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
		Parent:           fmt.Sprintf("organizations/%s/locations/global", orgID),
		BigQueryExport:   bigQueryExport,
		BigQueryExportId: bigQueryExportID,
	}

	_, err0 := client.CreateBigQueryExport(ctx, req)
	if err0 != nil {
		return fmt.Errorf("Failed to create BigQueryConfig: %w", err)
	}
	return nil
}

func cleanupBigQueryExport(t *testing.T, bigQueryExportID string) error {
	orgID := orgID(t)

	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycenterpb.DeleteBigQueryExportRequest{
		Name: fmt.Sprintf("organizations/%s/locations/global/bigQueryExports/%s", orgID, bigQueryExportID),
	}

	if err := client.DeleteBigQueryExport(ctx, req); err != nil {
		return fmt.Errorf("failed to delete BigQueryExport: %w", err)
	}

	return nil
}

func TestListBigQuery(t *testing.T) {

	testutil.Retry(t, 5, 20*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		// Create Test BigQueryExport Config
		rand, err := uuid.NewUUID()
		if err != nil {
			r.Errorf("Issue generating id.")
			return
		}
		configID := "random-bqexport-id-" + rand.String()

		if err := addBigQueryExport(t, configID); err != nil {
			r.Errorf("Could not setup test environment: %v", err)
			return
		}

		parent := fmt.Sprintf("organizations/%s/locations/global", orgID(t))

		// Call List BigQueryExport
		err = listBigQueryExport(buf, parent)

		if err != nil {
			r.Errorf("listBigQueryExport() had error: %v", err)
			return
		}

		got := buf.String()

		fmt.Println(got)

		if !strings.Contains(got, configID) {
			r.Errorf("listBigQueryConfigs() got: %s want %s", got, configID)
		}

		// Cleanup
		cleanupBigQueryExport(t, configID)
	})
}

func TestGetBigQuery(t *testing.T) {

	testutil.Retry(t, 5, 20*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		// Create Test BigQueryExport Config
		rand, err := uuid.NewUUID()
		if err != nil {
			r.Errorf("Issue generating id.")
			return
		}
		configID := "random-bqexport-id-" + rand.String()

		if err := addBigQueryExport(t, configID); err != nil {
			r.Errorf("Could not setup test environment: %v", err)
			return
		}

		parent := fmt.Sprintf("organizations/%s/locations/global", orgID(t))

		// Call GetBigQueryExport
		err = getBigQueryExport(buf, parent, configID)

		if err != nil {
			r.Errorf("getBigQueryExport() had error: %v", err)
			return
		}

		got := buf.String()

		fmt.Println(got)

		if !strings.Contains(got, configID) {
			r.Errorf("getBigQueryExport() got: %s want %s", got, configID)
		}

		// Cleanup
		cleanupBigQueryExport(t, configID)
	})
}

func TestDeleteBigQuery(t *testing.T) {

	testutil.Retry(t, 5, 20*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)
		// Create Test BigQueryExport Config
		rand, err := uuid.NewUUID()
		if err != nil {
			r.Errorf("Issue generating id.")
			return
		}
		configID := "random-bqexport-id-" + rand.String()

		if err := addBigQueryExport(t, configID); err != nil {
			r.Errorf("Could not setup test environment: %v", err)
			return
		}

		parent := fmt.Sprintf("organizations/%s/locations/global", orgID(t))

		// Call DeleteBigQueryExport
		err = deleteBigQueryExport(buf, parent, configID)

		if err != nil {
			r.Errorf("getBigQueryExport() had error: %v", err)
			return
		}

		got := buf.String()

		fmt.Println(got)

		if !strings.Contains(got, configID) {
			r.Errorf("deleteBigQueryExport() got: %s want %s", got, configID)
		}
	})
}

func TestUpdateBigQuery(t *testing.T) {

	testutil.Retry(t, 5, 20*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		// Create Test BigQueryExport Config
		rand, err := uuid.NewUUID()
		if err != nil {
			r.Errorf("Issue generating id.")
			return
		}
		configID := "random-bqexport-id-" + rand.String()

		if err := addBigQueryExport(t, configID); err != nil {
			r.Errorf("Could not setup test environment: %v", err)
			return
		}

		parent := fmt.Sprintf("organizations/%s/locations/global", orgID(t))

		// Call UpdateBigQueryExport
		err = updateBigQueryExport(buf, parent, configID)

		if err != nil {
			r.Errorf("getBigQueryExport() had error: %v", err)
			return
		}

		got := buf.String()

		fmt.Println(got)

		if !strings.Contains(got, configID) {
			r.Errorf("getBigQueryExport() got: %s want %s", got, configID)
		}

		// Cleanup
		cleanupBigQueryExport(t, configID)
	})
}

func TestCreateBigQuery(t *testing.T) {

	testutil.Retry(t, 5, 20*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		rand, err := uuid.NewUUID()
		if err != nil {
			r.Errorf("Issue generating id.")
			return
		}
		configID := "random-bqexport-id-" + rand.String()

		parent := fmt.Sprintf("organizations/%s/locations/global", orgID(t))

		err = createBigQueryExport(buf, parent, configID, bigQueryDatasetName(t))

		if err != nil {
			r.Errorf("getBigQueryExport() had error: %v", err)
			return
		}

		got := buf.String()

		fmt.Println(got)

		if !strings.Contains(got, configID) {
			r.Errorf("createBigQueryExport() got: %s want %s", got, configID)
		}

		// Cleanup
		cleanupBigQueryExport(t, configID)

	})
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
