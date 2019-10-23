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

package v3

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

func TestBatchTranslateTextWithModel(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer

	bucketName := fmt.Sprintf("%s-translate_model-%v", tc.ProjectID, uuid.New().ID())
	location := "us-central1"
	inputURI := "gs://cloud-samples-data/translation/custom_model_text.txt"
	outputURI := fmt.Sprintf("gs://%s/translation/output/", bucketName)
	sourceLang := "en"
	targetLang := "ja"
	modelID := "TRL3128559826197068699"

	// Create a temporary bucket to store annotation output.
	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer c.Close()

	bucket := c.Bucket(bucketName)
	if err := bucket.Create(ctx, tc.ProjectID, nil); err != nil {
		t.Fatalf("bucket.Create: %v", err)
	}
	defer bucket.Delete(ctx)

	// Translate a sample text and check the number of translated characters.
	if err := batchTranslateTextWithModel(
		&buf,
		tc.ProjectID,
		location,
		inputURI,
		outputURI,
		sourceLang,
		targetLang,
		modelID,
	); err != nil {
		t.Fatalf("batchTranslateTextWithModel: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "Total characters: 15") {
		t.Fatalf("Got '%s', expected to contain 'Total characters: 15'", got)
	}
}
