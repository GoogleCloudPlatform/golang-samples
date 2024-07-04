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
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	aiplatformpb "cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/option"
)

func TestCreateDataset(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := createDataset(&buf, tc.ProjectID, region); err != nil {
		t.Fatalf("createDataset: %v", err)
	}

	got := buf.String()
	log.Println(got)
	if !strings.Contains(got, "Created dataset ") {
		t.Errorf("createDataset: wanted 'Created dataset ', got '%s'", got)
	}

	output := got
	teardownCreateDataset(output, t)
}

func teardownCreateDataset(output string, t *testing.T) {
	t.Helper()

	// parse dataset name--we cannot predict the dataset ID at creation time.
	_, tmp, ok := strings.Cut(output, "\n")
	if !ok {
		log.Println("couldn't parse dataset resource name")
		return
	}

	datasetName := tmp
	log.Println(datasetName)

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", region)
	clientOption := option.WithEndpoint(apiEndpoint)

	ctx := context.Background()

	client, err := aiplatform.NewDatasetClient(ctx, clientOption)
	if err != nil {
		log.Fatalf("aiplatform.NewDatasetClient: %v", err)
	}
	defer client.Close()

	log.Println(datasetName)

	req := &aiplatformpb.DeleteDatasetRequest{
		Name: datasetName,
	}

	_, err = client.DeleteDataset(ctx, req)
	if err != nil {
		log.Fatalf("Dataset(%s).Delete: %v", datasetName, err)
	}
}
