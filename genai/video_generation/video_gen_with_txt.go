// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package video_generation shows how to use the GenAI SDK to generate video.
package video_generation

// [START googlegenaisdk_videogen_with_txt]
import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/genai"
)

// generateVideoWithText shows how to gen video from text.
func generateVideoWithText(w io.Writer, outputGCSURI string) error {
	//outputGCSURI = "gs://your-bucket/your-prefix"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
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

	// Polling until the operation is done
	for !operation.Done {
		time.Sleep(15 * time.Second)
		operation, err = client.Operations.GetVideosOperation(ctx, operation, nil)
		if err != nil {
			return fmt.Errorf("failed to get operation status: %w", err)
		}
	}

	if operation.Response != nil && len(operation.Response.GeneratedVideos) > 0 {
		videoURI := operation.Response.GeneratedVideos[0].Video.URI
		fmt.Fprintln(w, videoURI)
		return nil
	}

	// Example response:
	// gs://your-bucket/your-prefix/videoURI

	return fmt.Errorf("video generation failed or returned no results")
}

// [END googlegenaisdk_videogen_with_txt]
