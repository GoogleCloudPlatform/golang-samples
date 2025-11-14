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
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/genai"
)

type mockBatchesService struct{}

func (m *mockBatchesService) Create(
	ctx context.Context,
	model string,
	source *genai.BatchJobSource,
	config *genai.CreateBatchJobConfig,
) (*genai.BatchJob, error) {

	return &genai.BatchJob{
		Name:  "projects/test/locations/us-central1/batchPredictionJobs/1234567890",
		State: genai.JobStatePending,
	}, nil
}

func (m *mockBatchesService) Get(
	ctx context.Context,
	name string,
	_ interface{},
) (*genai.BatchJob, error) {

	return &genai.BatchJob{
		Name:  name,
		State: genai.JobStateSucceeded,
	}, nil
}

type mockGenAIClient struct {
	Batches *mockBatchesService
}

func generateBatchEmbeddingsMock(w io.Writer, outputURI string) error {
	ctx := context.Background()

	client := &mockGenAIClient{
		Batches: &mockBatchesService{},
	}

	job, err := client.Batches.Create(ctx,
		"text-embedding-005",
		&genai.BatchJobSource{
			Format: "jsonl",
			GCSURI: []string{"gs://cloud-samples-data/generative-ai/embeddings/embeddings_input.jsonl"},
		},
		&genai.CreateBatchJobConfig{
			Dest: &genai.BatchJobDestination{
				Format: "jsonl",
				GCSURI: outputURI,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create batch job: %w", err)
	}

	fmt.Fprintf(w, "Job name: %s\n", job.Name)
	fmt.Fprintf(w, "Job state: %s\n", job.State)

	completed := false
	for !completed {
		job, err = client.Batches.Get(ctx, job.Name, nil)
		if err != nil {
			return fmt.Errorf("failed to get batch job: %w", err)
		}

		fmt.Fprintf(w, "Job state: %s\n", job.State)

		if job.State == genai.JobStateSucceeded {
			completed = true
		}
	}

	return nil
}

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
		err := generateBatchEmbeddingsMock(buf, outputURI)
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
