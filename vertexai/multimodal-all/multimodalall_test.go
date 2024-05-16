// Copyright 2024 Google LLC
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

package multimodalall

import (
	"bytes"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func Test_generateContentFromVideoWithAudio(t *testing.T) {
	tc := testutil.SystemTest(t)

	buf := new(bytes.Buffer)

	prompt := multimodalPrompt{
		video: "gs://cloud-samples-data/generative-ai/video/behind_the_scenes_pixel.mp4",
		image: "gs://cloud-samples-data/generative-ai/image/a-man-and-a-dog.png",
		question: `
		Watch each frame in the video carefully and answer the questions.
		Only base your answers strictly on what information is available in the video attached.
		Do not make up any information that is not part of the video and do not be too
		verbose, be to the point.

		Questions:
		- When is the moment in the image happening in the video? Provide a timestamp.
		- What is the context of the moment and what does the narrator say about it?
	`,
	}
	location := "us-central1"
	modelName := "gemini-1.5-pro-preview-0409"

	err := generateContentFromVideoWithAudio(buf, prompt, tc.ProjectID, location, modelName)
	if err != nil {
		t.Errorf("Test_generateContentFromVideoWithAudio: %v", err.Error())
	}

	// The generated text may look like "The moment in the image happens at
	// approximately 00:49 in the video. The context of the moment is..."
}
