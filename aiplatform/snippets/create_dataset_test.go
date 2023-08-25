// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Prompt: Write a test for createDataset in createDataset.go

package snippets

import (
	"context"
	"log"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

func TestCreateDataset(t *testing.T) {
	tc := testutil.SystemTest(t)
	datasetID := uuid.New().String()
	defer deleteDataset(tc.ProjectID, datasetID)

	if err := createDataset(tc.ProjectID, datasetID); err != nil {
		t.Fatalf("createDataset: %v", err)
	}
}

func deleteDataset(projectID, datasetID string) {
	client, err := bigquery.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	if err := client.Dataset(datasetID).Delete(context.Background()); err != nil {
		log.Fatalf("Dataset(%q).Delete: %v", datasetID, err)
	}
}
