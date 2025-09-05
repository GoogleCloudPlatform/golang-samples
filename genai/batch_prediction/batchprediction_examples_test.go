// Copyright 2025 Google LLC
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

package batch_prediction

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const gcsOutputBucket = "golang-docs-samples-tests"

func TestBatchPrediction(t *testing.T) {
	tc := testutil.SystemTest(t)

	t.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "1")
	t.Setenv("GOOGLE_CLOUD_LOCATION", "us-central1")
	t.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	prefix := fmt.Sprintf("embeddings_output/%d", time.Now().UnixNano())
	outputURI := fmt.Sprintf("gs://%s/%s", gcsOutputBucket, prefix)
	buf := new(bytes.Buffer)

	t.Run("generate batch embeddings with GCS", func(t *testing.T) {
		buf.Reset()
		err := generateBatchEmbeddings(buf, outputURI)
		if err != nil {
			t.Fatalf("generateBatchEmbeddings failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate batch predict with gcs input/output", func(t *testing.T) {
		buf.Reset()
		err := generateBatchPredict(buf, outputURI)
		if err != nil {
			t.Fatalf("generateBatchPredict failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}

	})

	t.Run("generate batch predict with BigQuery", func(t *testing.T) {
		buf.Reset()
		outputURIBQ := "bq://your-project.your_dataset.your_table"

		err := generateBatchPredictWithBQ(buf, outputURIBQ)
		if err != nil {
			t.Fatalf("generateBatchPredictWithBQ failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})
}
