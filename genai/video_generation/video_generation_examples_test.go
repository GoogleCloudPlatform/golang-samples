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

package video_generation

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

type mockModelsServiceVideo struct{}

func (m *mockModelsServiceVideo) GenerateVideos(
	ctx context.Context,
	model string,
	prompt string,
	img *genai.Image,
	config *genai.GenerateVideosConfig,
) (*genai.GenerateVideosOperation, error) {
	return &genai.GenerateVideosOperation{
		Done: false,
	}, nil
}

type mockOperationsServiceVideo struct{}

func (m *mockOperationsServiceVideo) GetVideosOperation(
	ctx context.Context,
	op *genai.GenerateVideosOperation,
	_ *genai.GetOperationConfig,
) (*genai.GenerateVideosOperation, error) {

	return &genai.GenerateVideosOperation{
		Done: true,
		Response: &genai.GenerateVideosResponse{
			GeneratedVideos: []*genai.GeneratedVideo{
				{
					Video: &genai.Video{
						URI: "gs://mock-video/output.mp4",
					},
				},
			},
		},
	}, nil
}

type mockGenAIClientVideo struct {
	Models     *mockModelsServiceVideo
	Operations *mockOperationsServiceVideo
}

func generateVideoFromImageMock(w io.Writer, outputGCSURI string) error {
	ctx := context.Background()

	client := &mockGenAIClientVideo{
		Models:     &mockModelsServiceVideo{},
		Operations: &mockOperationsServiceVideo{},
	}

	image := &genai.Image{
		GCSURI:   "gs://cloud-samples-data/generative-ai/image/flowers.png",
		MIMEType: "image/png",
	}

	config := &genai.GenerateVideosConfig{
		AspectRatio:  "16:9",
		OutputGCSURI: outputGCSURI,
	}

	modelName := "veo-3.0-generate-preview"
	prompt := "Extreme close-up of a cluster of vibrant wildflowers swaying gently."

	operation, err := client.Models.GenerateVideos(ctx, modelName, prompt, image, config)
	if err != nil {
		return fmt.Errorf("failed to start video generation: %w", err)
	}

	for !operation.Done {
		operation, err = client.Operations.GetVideosOperation(ctx, operation, nil)
		if err != nil {
			return fmt.Errorf("failed to get operation status: %w", err)
		}
	}

	if operation.Response != nil &&
		len(operation.Response.GeneratedVideos) > 0 &&
		operation.Response.GeneratedVideos[0].Video != nil {

		uri := operation.Response.GeneratedVideos[0].Video.URI
		fmt.Fprintln(w, uri)
		return nil
	}

	return fmt.Errorf("video generation failed or returned no results")
}

type mockModelsServiceVideoText struct{}

func (m *mockModelsServiceVideoText) GenerateVideos(
	ctx context.Context,
	model string,
	prompt string,
	image *genai.Image,
	config *genai.GenerateVideosConfig,
) (*genai.GenerateVideosOperation, error) {
	return &genai.GenerateVideosOperation{
		Done: false,
	}, nil
}

type mockOperationsServiceVideoText struct{}

func (m *mockOperationsServiceVideoText) GetVideosOperation(
	ctx context.Context,
	op *genai.GenerateVideosOperation,
	_ *genai.GetOperationConfig,
) (*genai.GenerateVideosOperation, error) {
	return &genai.GenerateVideosOperation{
		Done: true,
		Response: &genai.GenerateVideosResponse{
			GeneratedVideos: []*genai.GeneratedVideo{
				{
					Video: &genai.Video{
						URI: "gs://mock-bucket/video-from-text.mp4",
					},
				},
			},
		},
	}, nil
}

type mockGenAIClientVideoText struct {
	Models     *mockModelsServiceVideoText
	Operations *mockOperationsServiceVideoText
}

func generateVideoWithTextMock(w io.Writer, outputGCSURI string) error {
	ctx := context.Background()

	client := &mockGenAIClientVideoText{
		Models:     &mockModelsServiceVideoText{},
		Operations: &mockOperationsServiceVideoText{},
	}

	config := &genai.GenerateVideosConfig{
		AspectRatio:  "16:9",
		OutputGCSURI: outputGCSURI,
	}

	modelName := "veo-3.0-generate-preview"
	prompt := "a cat reading a book"

	operation, err := client.Models.GenerateVideos(ctx, modelName, prompt, nil, config)
	if err != nil {
		return fmt.Errorf("failed to start video generation: %w", err)
	}

	for !operation.Done {
		operation, err = client.Operations.GetVideosOperation(ctx, operation, nil)
		if err != nil {
			return fmt.Errorf("failed to get operation status: %w", err)
		}
	}

	if operation.Response != nil &&
		len(operation.Response.GeneratedVideos) > 0 &&
		operation.Response.GeneratedVideos[0].Video != nil {

		fmt.Fprintln(w, operation.Response.GeneratedVideos[0].Video.URI)
		return nil
	}

	return fmt.Errorf("video generation failed or returned no results")
}

func TestVideoGeneration(t *testing.T) {
	tc := testutil.SystemTest(t)

	t.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "1")
	t.Setenv("GOOGLE_CLOUD_LOCATION", "us-central1")
	t.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	buf := new(bytes.Buffer)

	gcsOutputBucket := "HERE-go-bucket-samples-tests"
	prefix := "go_videogen_test/" + time.Now().Format("20060102-150405")
	outputGCSURI := "gs://" + gcsOutputBucket + "/" + prefix

	t.Run("generate video content with img", func(t *testing.T) {
		buf.Reset()
		err := generateVideoFromImageMock(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateVideoFromImage failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate video content with text", func(t *testing.T) {
		buf.Reset()
		err := generateVideoWithTextMock(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateVideoWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

}
