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

package image_generation

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/genai"
)

type mockModelsServiceStyleRef struct{}

func (m *mockModelsServiceStyleRef) EditImage(
	ctx context.Context,
	model string,
	prompt string,
	referenceImages []genai.ReferenceImage,
	config *genai.EditImageConfig,
) (*genai.EditImageResponse, error) {

	return &genai.EditImageResponse{
		GeneratedImages: []*genai.GeneratedImage{
			{
				Image: &genai.Image{
					GCSURI: config.OutputGCSURI,
				},
			},
		},
	}, nil
}

type mockGenAIClientStyleRef struct {
	Models *mockModelsServiceStyleRef
}

func generateStyleRefWithTextMock(w io.Writer, outputGCSURI string) error {
	ctx := context.Background()

	client := &mockGenAIClientStyleRef{
		Models: &mockModelsServiceStyleRef{},
	}

	styleRefImg := &genai.StyleReferenceImage{
		ReferenceID: 1,
		ReferenceImage: &genai.Image{
			GCSURI: "gs://cloud-samples-data/generative-ai/image/neon.png",
		},
		Config: &genai.StyleReferenceConfig{
			StyleDescription: "neon sign",
		},
	}

	prompt := "generate an image of a neon sign [1] with the words: have a great day"
	modelName := "imagen-3.0-capability-001"

	resp, err := client.Models.EditImage(
		ctx,
		modelName,
		prompt,
		[]genai.ReferenceImage{styleRefImg},
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

	fmt.Fprintln(w, resp.GeneratedImages[0].Image.GCSURI)
	return nil
}

type mockModelsServiceCanny struct{}

func (m *mockModelsServiceCanny) EditImage(
	ctx context.Context,
	model string,
	prompt string,
	referenceImages []genai.ReferenceImage,
	config *genai.EditImageConfig,
) (*genai.EditImageResponse, error) {
	return &genai.EditImageResponse{
		GeneratedImages: []*genai.GeneratedImage{
			{
				Image: &genai.Image{
					GCSURI: config.OutputGCSURI,
				},
			},
		},
	}, nil
}

type mockGenAIClientCanny struct {
	Models *mockModelsServiceCanny
}

func generateCannyCtrlTypeWithTextMock(w io.Writer, outputGCSURI string) error {
	ctx := context.Background()

	client := &mockGenAIClientCanny{
		Models: &mockModelsServiceCanny{},
	}

	controlReference := &genai.ControlReferenceConfig{
		ControlType: genai.ControlReferenceTypeCanny,
	}

	referenceImage := &genai.ControlReferenceImage{
		ReferenceID: 1,
		ReferenceImage: &genai.Image{
			GCSURI: "gs://cloud-samples-data/generative-ai/image/car_canny.png",
		},
		Config: controlReference,
	}

	modelName := "imagen-3.0-capability-001"
	prompt := "a watercolor painting of a red car[1] driving on a road"

	resp, err := client.Models.EditImage(
		ctx,
		modelName,
		prompt,
		[]genai.ReferenceImage{referenceImage},
		&genai.EditImageConfig{
			EditMode:          genai.EditModeControlledEditing,
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
		return fmt.Errorf("no image generated")
	}

	fmt.Fprintln(w, resp.GeneratedImages[0].Image.GCSURI)
	return nil
}

type mockModelsServiceScribble struct{}

func (m *mockModelsServiceScribble) EditImage(
	ctx context.Context,
	model string,
	prompt string,
	referenceImages []genai.ReferenceImage,
	config *genai.EditImageConfig,
) (*genai.EditImageResponse, error) {

	return &genai.EditImageResponse{
		GeneratedImages: []*genai.GeneratedImage{
			{
				Image: &genai.Image{
					GCSURI: config.OutputGCSURI,
				},
			},
		},
	}, nil
}

type mockGenAIClientScribble struct {
	Models *mockModelsServiceScribble
}

func generateScribbleCtrlTypeWithTextMock(w io.Writer, outputGCSURI string) error {
	ctx := context.Background()

	client := &mockGenAIClientScribble{
		Models: &mockModelsServiceScribble{},
	}

	controlReference := &genai.ControlReferenceConfig{
		ControlType: genai.ControlReferenceTypeScribble,
	}

	referenceImage := &genai.ControlReferenceImage{
		ReferenceID: 1,
		ReferenceImage: &genai.Image{
			GCSURI: "gs://cloud-samples-data/generative-ai/image/car_scribble.png",
		},
		Config: controlReference,
	}

	modelName := "imagen-3.0-capability-001"
	prompt := "an oil painting showing the side of a red car[1]"

	resp, err := client.Models.EditImage(
		ctx,
		modelName,
		prompt,
		[]genai.ReferenceImage{referenceImage},
		&genai.EditImageConfig{
			EditMode:          genai.EditModeControlledEditing,
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
		return fmt.Errorf("no image generated")
	}

	fmt.Fprintln(w, resp.GeneratedImages[0].Image.GCSURI)
	return nil
}

type mockModelsServiceSubjCtrl struct{}

func (m *mockModelsServiceSubjCtrl) EditImage(
	ctx context.Context,
	model string,
	prompt string,
	referenceImages []genai.ReferenceImage,
	config *genai.EditImageConfig,
) (*genai.EditImageResponse, error) {

	return &genai.EditImageResponse{
		GeneratedImages: []*genai.GeneratedImage{
			{
				Image: &genai.Image{
					GCSURI: config.OutputGCSURI,
				},
			},
		},
	}, nil
}

type mockGenAIClientSubjCtrl struct {
	Models *mockModelsServiceSubjCtrl
}

func generateSubjRefCtrlReferWithTextMock(w io.Writer, outputGCSURI string) error {
	ctx := context.Background()

	client := &mockGenAIClientSubjCtrl{
		Models: &mockModelsServiceSubjCtrl{},
	}

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

	prompt := "a portrait of a woman[1] in the pose of the control image[2] in a watercolor style..."
	modelName := "imagen-3.0-capability-001"

	resp, err := client.Models.EditImage(
		ctx,
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

	fmt.Fprintln(w, resp.GeneratedImages[0].Image.GCSURI)
	return nil
}

type mockModelsServiceRawRef struct{}

func (m *mockModelsServiceRawRef) EditImage(
	ctx context.Context,
	model string,
	prompt string,
	referenceImages []genai.ReferenceImage,
	config *genai.EditImageConfig,
) (*genai.EditImageResponse, error) {

	return &genai.EditImageResponse{
		GeneratedImages: []*genai.GeneratedImage{
			{
				Image: &genai.Image{
					GCSURI: config.OutputGCSURI,
				},
			},
		},
	}, nil
}

func generateRawReferWithTextMock(w io.Writer, outputGCSURI string) error {
	ctx := context.Background()

	client := &mockGenAIClientRawRef{
		Models: &mockModelsServiceRawRef{},
	}

	rawRefImage := &genai.RawReferenceImage{
		ReferenceID: 1,
		ReferenceImage: &genai.Image{
			GCSURI: "gs://cloud-samples-data/generative-ai/image/teacup-1.png",
		},
	}

	prompt := "transform the subject in the image so that the teacup[1] is made entirely out of chocolate"
	modelName := "imagen-3.0-capability-001"

	resp, err := client.Models.EditImage(
		ctx,
		modelName,
		prompt,
		[]genai.ReferenceImage{rawRefImage},
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

	fmt.Fprintln(w, resp.GeneratedImages[0].Image.GCSURI)
	return nil
}

type mockGenAIClientRawRef struct {
	Models *mockModelsServiceRawRef
}

func TestImageGeneration(t *testing.T) {
	tc := testutil.SystemTest(t)

	t.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "1")
	t.Setenv("GOOGLE_CLOUD_LOCATION", "global")
	t.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	buf := new(bytes.Buffer)

	t.Run("generate multimodal flash content with text and image", func(t *testing.T) {
		buf.Reset()
		err := generateMMFlashWithText(buf)
		if err != nil {
			t.Fatalf("generateMMFlashWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate mmflash text and image recipe", func(t *testing.T) {
		buf.Reset()
		err := generateMMFlashTxtImgWithText(buf)
		if err != nil {
			t.Fatalf("generateMMFlashTxtImgWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("style customization with style reference", func(t *testing.T) {
		buf.Reset()
		// TODO(developer): update with your bucket
		outputGCSURI := "gs://your-bucket/your-prefix"

		err := generateStyleRefWithTextMock(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateStyleRefWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected printed output, got empty")
		}
	})

	t.Run("canny edge customization with text+image", func(t *testing.T) {
		buf.Reset()
		// TODO(developer): update with your bucket
		outputGCSURI := "gs://your-bucket/your-prefix"

		err := generateCannyCtrlTypeWithTextMock(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateCannyCtrlTypeWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate image with scribble control type", func(t *testing.T) {
		buf.Reset()
		// TODO(developer): update with your bucket
		outputGCSURI := "gs://your-bucket/your-prefix"

		err := generateScribbleCtrlTypeWithTextMock(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateScribbleCtrlTypeWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("subject customization with control reference", func(t *testing.T) {
		buf.Reset()
		// TODO(developer): update with your bucket
		outputGCSURI := "gs://your-bucket/your-prefix"

		err := generateSubjRefCtrlReferWithTextMock(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateSubjRefCtrlReferWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate style transfer customization with raw reference", func(t *testing.T) {
		buf.Reset()
		// TODO(developer): update with your bucket
		outputGCSURI := "gs://your-bucket/your-prefix"

		err := generateRawReferWithTextMock(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateRawReferWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate image content with text", func(t *testing.T) {
		buf.Reset()
		err := generateImageWithText(buf)
		if err != nil {
			t.Fatalf("generateImageWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate mmflash image content with text and image", func(t *testing.T) {
		buf.Reset()
		err := generateImageMMFlashEditWithTextImg(buf)
		if err != nil {
			t.Fatalf("generateImageMMFlashEditWithTextImg failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("style customization with style reference", func(t *testing.T) {
		buf.Reset()
		// TODO(developer): update with your bucket
		outputGCSURI := "gs://your-bucket/your-prefix"

		err := generateStyleRefWithTextMock(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateStyleRefWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected printed output, got empty")
		}
	})

	t.Run("canny edge customization with text+image", func(t *testing.T) {
		buf.Reset()
		// TODO(developer): update with your bucket
		outputGCSURI := "gs://your-bucket/your-prefix"

		err := generateCannyCtrlTypeWithTextMock(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateCannyCtrlTypeWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate image with scribble control type", func(t *testing.T) {
		buf.Reset()
		// TODO(developer): update with your bucket
		outputGCSURI := "gs://your-bucket/your-prefix"

		err := generateScribbleCtrlTypeWithTextMock(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateScribbleCtrlTypeWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("subject customization with control reference", func(t *testing.T) {
		buf.Reset()
		// TODO(developer): update with your bucket
		outputGCSURI := "gs://your-bucket/your-prefix"

		err := generateSubjRefCtrlReferWithTextMock(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateSubjRefCtrlReferWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("generate style transfer customization with raw reference", func(t *testing.T) {
		buf.Reset()
		// TODO(developer): update with your bucket
		outputGCSURI := "gs://your-bucket/your-prefix"

		err := generateRawReferWithTextMock(buf, outputGCSURI)
		if err != nil {
			t.Fatalf("generateRawReferWithText failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})
}
