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

// Package text_generation shows examples of generating text using the GenAI SDK.
package text_generation

// [START googlegenaisdk_textgen_with_local_video]
import (
	"context"
	"fmt"
	"io"
	"os"

	genai "google.golang.org/genai"
)

// generateWithLocalVideo shows how to generate text using a local video input.
func generateWithLocalVideo(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Read local video file content
	data, err := os.ReadFile("testdata/describe_video_content.mp4")
	if err != nil {
		return fmt.Errorf("failed to read local video: %w", err)
	}

	modelName := "gemini-2.5-flash"
	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: `Analyze the provided local video file, including its audio.
Summarize the main points of the video concisely.
Create a chapter breakdown with timestamps for key sections or topics discussed.`},
				{InlineData: &genai.Blob{
					MIMEType: "video/mp4",
					Data:     data,
				}},
			},
		},
	}

	resp, err := client.Models.GenerateContent(ctx, modelName, contents, nil)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	respText := resp.Text()
	fmt.Fprintln(w, respText)

	// Example response:
	// The video features a male climber engaged in lead rock climbing at an indoor gym. He is shown ascending a vertical wall adorned with various colored holds, demonstrating fluid movements and proper technique. During his ascent, he clips his rope into a quickdraw for safety and manages the rope for continued climbing. The video is silent.
	//
	// **Summary of Main Points:**
	//
	// The video shows an indoor rock climber, equipped with a harness and chalk bag, efficiently ascending a climbing wall. He skillfully clips...
	//
	// **Chapter Breakdown**
	//
	// *  **0:00 - Beginning of Ascent:** The climber starts his ascent, moving gracefully between holds.
	//*   **0:01 - Clipping the Rope into a Quickdraw:** The climber pauses to clip his safety rope into a quickdraw, securing his position.
	//*   **0:03 - Continued Climbing and Technique Display:** The climber resumes his upward movement, demonstrating his climbing technique and body control.
	//*   **0:08 - Rope Management / Preparing for Next Moves:** The climber adjusts the rope, taking up slack, and surveys the next section of the climb.
	//*   **0:13 - Video End:** The video concludes with the climber still in the process of ascending the wall.
	// ...

	return nil
}

// [END googlegenaisdk_textgen_with_local_video]
