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

package tuning

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/genai"
)

type mockTunings struct{}
type mockModels struct{}

func (m *mockTunings) CreateMock(ctx context.Context) (*genai.TuningJob, error) {
	return &genai.TuningJob{
		Name: "projects/mock/locations/us-central1/tuningJobs/test-tuning-job",
	}, nil
}

func (m *mockTunings) GetMock(ctx context.Context, name string) (*genai.TuningJob, error) {
	if name != "test-tuning-job" {
		return nil, fmt.Errorf("mock: tuning job not found")
	}
	return &genai.TuningJob{
		Name: name,
		TunedModel: &genai.TunedModel{
			Endpoint: "projects/mock/locations/global/endpoints/tuned-model-123",
		},
		State: "JOB_STATE_SUCCEEDED",
	}, nil
}

func (m *mockModels) GenerateContentMock(ctx context.Context, model string, contents []*genai.Content, cfg *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	return &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []*genai.Part{{Text: "Mocked response: The sky is blue because of Rayleigh scattering."}},
				},
			},
		},
	}, nil
}

func createTuningJobMock(w io.Writer) error {
	job, err := (&mockTunings{}).CreateMock(context.Background())
	if err != nil {
		return fmt.Errorf("mock create failed: %w", err)
	}
	fmt.Fprintf(w, "Created tuning job: %s\n", job.Name)
	return nil
}

func predictWithTunedEndpointMock(w io.Writer, tuningJobName string) error {
	tunings := &mockTunings{}
	models := &mockModels{}
	ctx := context.Background()

	job, err := tunings.GetMock(ctx, tuningJobName)
	if err != nil {
		return fmt.Errorf("mock get tuning job failed: %w", err)
	}

	contents := []*genai.Content{
		{
			Role: genai.RoleUser,
			Parts: []*genai.Part{
				{Text: "Why is the sky blue?"},
			},
		},
	}

	resp, err := models.GenerateContentMock(ctx, job.TunedModel.Endpoint, contents, nil)
	if err != nil {
		return fmt.Errorf("mock generate content failed: %w", err)
	}

	fmt.Fprintln(w, resp.Text())
	return nil
}

func getTuningJobMock(w io.Writer, tuningJobName string) error {
	tunings := &mockTunings{}
	ctx := context.Background()

	job, err := tunings.GetMock(ctx, tuningJobName)
	if err != nil {
		return fmt.Errorf("mock get tuning job failed: %w", err)
	}
	fmt.Fprintf(w, "Job %s found, state: %s\n", job.Name, job.State)
	return nil
}

func TestTuningGeneration(t *testing.T) {
	tc := testutil.SystemTest(t)

	t.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "1")
	t.Setenv("GOOGLE_CLOUD_LOCATION", "us-central1")
	t.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	buf := new(bytes.Buffer)

	t.Run("create tuning job in project", func(t *testing.T) {
		buf.Reset()
		err := createTuningJobMock(buf)
		if err != nil {
			t.Fatalf("createTuningJob failed: %v", err)
		}
		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("predictWithTunedEndpoint in project", func(t *testing.T) {
		buf.Reset()
		err := predictWithTunedEndpointMock(buf, "test-tuning-job")
		if err != nil {
			t.Fatalf("predictWithTunedEndpoint failed: %v", err)
		}
		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("get Tuning Job in project", func(t *testing.T) {
		buf.Reset()
		err := getTuningJobMock(buf, "test-tuning-job")
		if err != nil {
			t.Fatalf("getTuningJob failed: %v", err)
		}
		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("list tuning jobs in project", func(t *testing.T) {
		buf.Reset()
		err := listTuningJobs(buf)
		if err != nil {
			t.Fatalf("listTuningJobs failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})
}
