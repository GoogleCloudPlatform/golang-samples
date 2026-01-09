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

// Package image_generation shows how to use the GenAI SDK to generate images from prompt.

package image_generation

// [START googlegenaisdk_imggen_subj_refer_ctrl_refer_with_txt_imgs]
import (
	"context"
	"fmt"
	"io"

	"google.golang.org/genai"
)

// generateSubjRefCtrlReferWithText demonstrates subject & control reference customization.
func generateSubjRefCtrlReferWithText(w io.Writer, outputGCSURI string) error {
	//outputGCSURI = "gs://your-bucket/your-prefix"
	ctx := context.Background()

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	// Create subject and control reference images of a photograph stored in Google Cloud Storage
	// using https://storage.googleapis.com/cloud-samples-data/generative-ai/image/person.png
	subjectReferenceImage := &genai.SubjectReferenceImage{
		ReferenceID: 1,
		ReferenceImage: &genai.Image{
			GCSURI: "gs://cloud-samples-data/generative-ai/image/person.png",
		},
		Config: &genai.SubjectReferenceConfig{
			SubjectType:        genai.SubjectReferenceTypeSubjectTypePerson,
			SubjectDescription: "a headshot of a woman",
		},
	}

	controlReferenceImage := &genai.ControlReferenceImage{
		ReferenceID: 2,
		ReferenceImage: &genai.Image{
			GCSURI: "gs://cloud-samples-data/generative-ai/image/person.png",
		},
		Config: &genai.ControlReferenceConfig{
			ControlType: genai.ControlReferenceTypeFaceMesh,
		},
	}

	// prompt that references the style image with [1]
	prompt := "a portrait of a woman[1] in the pose of the control image[2] in a watercolor style by a professional artist, light and low-contrast strokes, bright pastel colors, a warm atmosphere, clean background, grainy paper, bold visible brushstrokes, patchy details"
	modelName := "imagen-3.0-capability-001"

	resp, err := client.Models.EditImage(ctx,
		modelName,
		prompt,
		[]genai.ReferenceImage{
			subjectReferenceImage,
			controlReferenceImage,
		},
		&genai.EditImageConfig{
			EditMode:          genai.EditModeDefault,
			NumberOfImages:    1,
			SafetyFilterLevel: genai.SafetyFilterLevelBlockMediumAndAbove,
			PersonGeneration:  genai.PersonGenerationAllowAdult,
			OutputGCSURI:      outputGCSURI,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to edit image: %w", err)
	}

	if len(resp.GeneratedImages) == 0 || resp.GeneratedImages[0].Image == nil {
		return fmt.Errorf("no generated images returned")
	}

	uri := resp.GeneratedImages[0].Image.GCSURI
	fmt.Fprintln(w, uri)

	// Example response:
	// gs://your-bucket/your-prefix
	return nil
}

// [END googlegenaisdk_imggen_subj_refer_ctrl_refer_with_txt_imgs]
