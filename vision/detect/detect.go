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

package main

// [START imports]
import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// [END imports]

func init() {
	// Refer to these functions so that goimports is happy before boilerplate is inserted.
	_ = context.Background()
	_ = vision.ImageAnnotatorClient{}
	_ = os.Open
}

// [START vision_face_detection]

// detectFaces gets faces from the Vision API for an image at the given file path.
func detectFaces(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	annotations, err := client.DetectFaces(ctx, image, nil, 10)
	if err != nil {
		return err
	}
	if len(annotations) == 0 {
		fmt.Fprintln(w, "No faces found.")
	} else {
		fmt.Fprintln(w, "Faces:")
		for i, annotation := range annotations {
			fmt.Fprintln(w, "  Face", i)
			fmt.Fprintln(w, "    Anger:", annotation.AngerLikelihood)
			fmt.Fprintln(w, "    Joy:", annotation.JoyLikelihood)
			fmt.Fprintln(w, "    Surprise:", annotation.SurpriseLikelihood)
		}
	}
	return nil
}

// [END vision_face_detection]

// [START vision_label_detection]

// detectLabels gets labels from the Vision API for an image at the given file path.
func detectLabels(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	annotations, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No labels found.")
	} else {
		fmt.Fprintln(w, "Labels:")
		for _, annotation := range annotations {
			fmt.Fprintln(w, annotation.Description)
		}
	}

	return nil
}

// [END vision_label_detection]

// [START vision_landmark_detection]

// detectLandmarks gets landmarks from the Vision API for an image at the given file path.
func detectLandmarks(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	annotations, err := client.DetectLandmarks(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No landmarks found.")
	} else {
		fmt.Fprintln(w, "Landmarks:")
		for _, annotation := range annotations {
			fmt.Fprintln(w, annotation.Description)
		}
	}

	return nil
}

// [END vision_landmark_detection]

// [START vision_text_detection]

// detectText gets text from the Vision API for an image at the given file path.
func detectText(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No text found.")
	} else {
		fmt.Fprintln(w, "Text:")
		for _, annotation := range annotations {
			fmt.Fprintf(w, "%q\n", annotation.Description)
		}
	}

	return nil
}

// [END vision_text_detection]

// [START vision_fulltext_detection]

// detectDocumentText gets the full document text from the Vision API for an image at the given file path.
func detectDocumentText(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	annotation, err := client.DetectDocumentText(ctx, image, nil)
	if err != nil {
		return err
	}

	if annotation == nil {
		fmt.Fprintln(w, "No text found.")
	} else {
		fmt.Fprintln(w, "Document Text:")
		fmt.Fprintf(w, "%q\n", annotation.Text)

		fmt.Fprintln(w, "Pages:")
		for _, page := range annotation.Pages {
			fmt.Fprintf(w, "\tConfidence: %f, Width: %d, Height: %d\n", page.Confidence, page.Width, page.Height)
			fmt.Fprintln(w, "\tBlocks:")
			for _, block := range page.Blocks {
				fmt.Fprintf(w, "\t\tConfidence: %f, Block type: %v\n", block.Confidence, block.BlockType)
				fmt.Fprintln(w, "\t\tParagraphs:")
				for _, paragraph := range block.Paragraphs {
					fmt.Fprintf(w, "\t\t\tConfidence: %f", paragraph.Confidence)
					fmt.Fprintln(w, "\t\t\tWords:")
					for _, word := range paragraph.Words {
						symbols := make([]string, len(word.Symbols))
						for i, s := range word.Symbols {
							symbols[i] = s.Text
						}
						wordText := strings.Join(symbols, "")
						fmt.Fprintf(w, "\t\t\t\tConfidence: %f, Symbols: %s\n", word.Confidence, wordText)
					}
				}
			}
		}
	}

	return nil
}

// [END vision_fulltext_detection]

// [START vision_image_property_detection]

// detectProperties gets image properties from the Vision API for an image at the given file path.
func detectProperties(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	props, err := client.DetectImageProperties(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Dominant colors:")
	for _, quantized := range props.DominantColors.Colors {
		color := quantized.Color
		r := int(color.Red) & 0xff
		g := int(color.Green) & 0xff
		b := int(color.Blue) & 0xff
		fmt.Fprintf(w, "%2.1f%% - #%02x%02x%02x\n", quantized.PixelFraction*100, r, g, b)
	}

	return nil
}

// [END vision_image_property_detection]

// [START vision_crop_hint_detection]

// detectCropHints gets suggested croppings the Vision API for an image at the given file path.
func detectCropHints(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	res, err := client.CropHints(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Crop hints:")
	for _, hint := range res.CropHints {
		for _, v := range hint.BoundingPoly.Vertices {
			fmt.Fprintf(w, "(%d,%d)\n", v.X, v.Y)
		}
	}

	return nil
}

// [END vision_crop_hint_detection]

// [START vision_safe_search_detection]

// detectSafeSearch gets image properties from the Vision API for an image at the given file path.
func detectSafeSearch(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	props, err := client.DetectSafeSearch(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Safe Search properties:")
	fmt.Fprintln(w, "Adult:", props.Adult)
	fmt.Fprintln(w, "Medical:", props.Medical)
	fmt.Fprintln(w, "Racy:", props.Racy)
	fmt.Fprintln(w, "Spoofed:", props.Spoof)
	fmt.Fprintln(w, "Violence:", props.Violence)

	return nil
}

// [END vision_safe_search_detection]

// [START vision_web_detection]

// detectWeb gets image properties from the Vision API for an image at the given file path.
func detectWeb(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	web, err := client.DetectWeb(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Web properties:")
	if len(web.FullMatchingImages) != 0 {
		fmt.Fprintln(w, "\tFull image matches:")
		for _, full := range web.FullMatchingImages {
			fmt.Fprintf(w, "\t\t%s\n", full.Url)
		}
	}
	if len(web.PagesWithMatchingImages) != 0 {
		fmt.Fprintln(w, "\tPages with this image:")
		for _, page := range web.PagesWithMatchingImages {
			fmt.Fprintf(w, "\t\t%s\n", page.Url)
		}
	}
	if len(web.WebEntities) != 0 {
		fmt.Fprintln(w, "\tEntities:")
		fmt.Fprintln(w, "\t\tEntity\t\tScore\tDescription")
		for _, entity := range web.WebEntities {
			fmt.Fprintf(w, "\t\t%-14s\t%-2.4f\t%s\n", entity.EntityId, entity.Score, entity.Description)
		}
	}
	if len(web.BestGuessLabels) != 0 {
		fmt.Fprintln(w, "\tBest guess labels:")
		for _, label := range web.BestGuessLabels {
			fmt.Fprintf(w, "\t\t%s\n", label.Label)
		}
	}

	return nil
}

// [END vision_web_detection]

// [START vision_web_detection_include_geo]

// detectWebGeo detects geographic metadata from the Vision API for an image at the given file path.
func detectWebGeo(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	imageContext := &visionpb.ImageContext{
		WebDetectionParams: &visionpb.WebDetectionParams{
			IncludeGeoResults: true,
		},
	}
	web, err := client.DetectWeb(ctx, image, imageContext)
	if err != nil {
		return err
	}

	if len(web.WebEntities) != 0 {
		fmt.Fprintln(w, "Entities:")
		fmt.Fprintln(w, "\tEntity\t\tScore\tDescription")
		for _, entity := range web.WebEntities {
			fmt.Fprintf(w, "\t%-14s\t%-2.4f\t%s\n", entity.EntityId, entity.Score, entity.Description)
		}
	}

	return nil
}

// [END vision_web_detection_include_geo]

// [START vision_logo_detection]

// detectLogos gets logos from the Vision API for an image at the given file path.
func detectLogos(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	annotations, err := client.DetectLogos(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No logos found.")
	} else {
		fmt.Fprintln(w, "Logos:")
		for _, annotation := range annotations {
			fmt.Fprintln(w, annotation.Description)
		}
	}

	return nil
}

// [END vision_logo_detection]

// [START vision_localize_objects]

// localizeObjects gets objects and bounding boxes from the Vision API for an image at the given file path.
func localizeObjects(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}
	annotations, err := client.LocalizeObjects(ctx, image, nil)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No objects found.")
		return nil
	}

	fmt.Fprintln(w, "Objects:")
	for _, annotation := range annotations {
		fmt.Fprintln(w, annotation.Name)
		fmt.Fprintln(w, annotation.Score)

		for _, v := range annotation.BoundingPoly.NormalizedVertices {
			fmt.Fprintf(w, "(%f,%f)\n", v.X, v.Y)
		}
	}

	return nil
}

// [END vision_localize_objects]

func init() {
	// Refer to these functions so that goimports is happy before boilerplate is inserted.
	_ = context.Background()
	_ = vision.ImageAnnotatorClient{}
	_ = os.Open
}

// [START vision_face_detection_gcs]

// detectFaces gets faces from the Vision API for an image at the given file path.
func detectFacesURI(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	image := vision.NewImageFromURI(file)
	annotations, err := client.DetectFaces(ctx, image, nil, 10)
	if err != nil {
		return err
	}
	if len(annotations) == 0 {
		fmt.Fprintln(w, "No faces found.")
	} else {
		fmt.Fprintln(w, "Faces:")
		for i, annotation := range annotations {
			fmt.Fprintln(w, "  Face", i)
			fmt.Fprintln(w, "    Anger:", annotation.AngerLikelihood)
			fmt.Fprintln(w, "    Joy:", annotation.JoyLikelihood)
			fmt.Fprintln(w, "    Surprise:", annotation.SurpriseLikelihood)
		}
	}
	return nil
}

// [END vision_face_detection_gcs]

// [START vision_label_detection_gcs]

// detectLabels gets labels from the Vision API for an image at the given file path.
func detectLabelsURI(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	image := vision.NewImageFromURI(file)
	annotations, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No labels found.")
	} else {
		fmt.Fprintln(w, "Labels:")
		for _, annotation := range annotations {
			fmt.Fprintln(w, annotation.Description)
		}
	}

	return nil
}

// [END vision_label_detection_gcs]

// [START vision_landmark_detection_gcs]

// detectLandmarks gets landmarks from the Vision API for an image at the given file path.
func detectLandmarksURI(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	image := vision.NewImageFromURI(file)
	annotations, err := client.DetectLandmarks(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No landmarks found.")
	} else {
		fmt.Fprintln(w, "Landmarks:")
		for _, annotation := range annotations {
			fmt.Fprintln(w, annotation.Description)
		}
	}

	return nil
}

// [END vision_landmark_detection_gcs]

// [START vision_text_detection_gcs]

// detectText gets text from the Vision API for an image at the given file path.
func detectTextURI(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	image := vision.NewImageFromURI(file)
	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No text found.")
	} else {
		fmt.Fprintln(w, "Text:")
		for _, annotation := range annotations {
			fmt.Fprintf(w, "%q\n", annotation.Description)
		}
	}

	return nil
}

// [END vision_text_detection_gcs]

// [START vision_fulltext_detection_gcs]

// detectDocumentText gets the full document text from the Vision API for an image at the given file path.
func detectDocumentTextURI(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	image := vision.NewImageFromURI(file)
	annotation, err := client.DetectDocumentText(ctx, image, nil)
	if err != nil {
		return err
	}

	if annotation == nil {
		fmt.Fprintln(w, "No text found.")
	} else {
		fmt.Fprintln(w, "Document Text:")
		fmt.Fprintf(w, "%q\n", annotation.Text)

		fmt.Fprintln(w, "Pages:")
		for _, page := range annotation.Pages {
			fmt.Fprintf(w, "\tConfidence: %f, Width: %d, Height: %d\n", page.Confidence, page.Width, page.Height)
			fmt.Fprintln(w, "\tBlocks:")
			for _, block := range page.Blocks {
				fmt.Fprintf(w, "\t\tConfidence: %f, Block type: %v\n", block.Confidence, block.BlockType)
				fmt.Fprintln(w, "\t\tParagraphs:")
				for _, paragraph := range block.Paragraphs {
					fmt.Fprintf(w, "\t\t\tConfidence: %f", paragraph.Confidence)
					fmt.Fprintln(w, "\t\t\tWords:")
					for _, word := range paragraph.Words {
						symbols := make([]string, len(word.Symbols))
						for i, s := range word.Symbols {
							symbols[i] = s.Text
						}
						wordText := strings.Join(symbols, "")
						fmt.Fprintf(w, "\t\t\t\tConfidence: %f, Symbols: %s\n", word.Confidence, wordText)
					}
				}
			}
		}
	}

	return nil
}

// [END vision_fulltext_detection_gcs]

// [START vision_image_property_detection_gcs]

// detectProperties gets image properties from the Vision API for an image at the given file path.
func detectPropertiesURI(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	image := vision.NewImageFromURI(file)
	props, err := client.DetectImageProperties(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Dominant colors:")
	for _, quantized := range props.DominantColors.Colors {
		color := quantized.Color
		r := int(color.Red) & 0xff
		g := int(color.Green) & 0xff
		b := int(color.Blue) & 0xff
		fmt.Fprintf(w, "%2.1f%% - #%02x%02x%02x\n", quantized.PixelFraction*100, r, g, b)
	}

	return nil
}

// [END vision_image_property_detection_gcs]

// [START vision_crop_hint_detection_gcs]

// detectCropHints gets suggested croppings the Vision API for an image at the given file path.
func detectCropHintsURI(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	image := vision.NewImageFromURI(file)
	res, err := client.CropHints(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Crop hints:")
	for _, hint := range res.CropHints {
		for _, v := range hint.BoundingPoly.Vertices {
			fmt.Fprintf(w, "(%d,%d)\n", v.X, v.Y)
		}
	}

	return nil
}

// [END vision_crop_hint_detection_gcs]

// [START vision_safe_search_detection_gcs]

// detectSafeSearch gets image properties from the Vision API for an image at the given file path.
func detectSafeSearchURI(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	image := vision.NewImageFromURI(file)
	props, err := client.DetectSafeSearch(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Safe Search properties:")
	fmt.Fprintln(w, "Adult:", props.Adult)
	fmt.Fprintln(w, "Medical:", props.Medical)
	fmt.Fprintln(w, "Racy:", props.Racy)
	fmt.Fprintln(w, "Spoofed:", props.Spoof)
	fmt.Fprintln(w, "Violence:", props.Violence)

	return nil
}

// [END vision_safe_search_detection_gcs]

// [START vision_web_detection_gcs]

// detectWeb gets image properties from the Vision API for an image at the given file path.
func detectWebURI(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	image := vision.NewImageFromURI(file)
	web, err := client.DetectWeb(ctx, image, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Web properties:")
	if len(web.FullMatchingImages) != 0 {
		fmt.Fprintln(w, "\tFull image matches:")
		for _, full := range web.FullMatchingImages {
			fmt.Fprintf(w, "\t\t%s\n", full.Url)
		}
	}
	if len(web.PagesWithMatchingImages) != 0 {
		fmt.Fprintln(w, "\tPages with this image:")
		for _, page := range web.PagesWithMatchingImages {
			fmt.Fprintf(w, "\t\t%s\n", page.Url)
		}
	}
	if len(web.WebEntities) != 0 {
		fmt.Fprintln(w, "\tEntities:")
		fmt.Fprintln(w, "\t\tEntity\t\tScore\tDescription")
		for _, entity := range web.WebEntities {
			fmt.Fprintf(w, "\t\t%-14s\t%-2.4f\t%s\n", entity.EntityId, entity.Score, entity.Description)
		}
	}
	if len(web.BestGuessLabels) != 0 {
		fmt.Fprintln(w, "\tBest guess labels:")
		for _, label := range web.BestGuessLabels {
			fmt.Fprintf(w, "\t\t%s\n", label.Label)
		}
	}

	return nil
}

// [END vision_web_detection_gcs]

// [START vision_web_detection_include_geo_gcs]

// detectWebGeo detects geographic metadata from the Vision API for an image at the given file path.
func detectWebGeoURI(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	image := vision.NewImageFromURI(file)
	imageContext := &visionpb.ImageContext{
		WebDetectionParams: &visionpb.WebDetectionParams{
			IncludeGeoResults: true,
		},
	}
	web, err := client.DetectWeb(ctx, image, imageContext)
	if err != nil {
		return err
	}

	if len(web.WebEntities) != 0 {
		fmt.Fprintln(w, "Entities:")
		fmt.Fprintln(w, "\tEntity\t\tScore\tDescription")
		for _, entity := range web.WebEntities {
			fmt.Fprintf(w, "\t%-14s\t%-2.4f\t%s\n", entity.EntityId, entity.Score, entity.Description)
		}
	}

	return nil
}

// [END vision_web_detection_include_geo_gcs]

// [START vision_logo_detection_gcs]

// detectLogos gets logos from the Vision API for an image at the given file path.
func detectLogosURI(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	image := vision.NewImageFromURI(file)
	annotations, err := client.DetectLogos(ctx, image, nil, 10)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No logos found.")
	} else {
		fmt.Fprintln(w, "Logos:")
		for _, annotation := range annotations {
			fmt.Fprintln(w, annotation.Description)
		}
	}

	return nil
}

// [END vision_logo_detection_gcs]

// [START vision_text_detection_pdf_gcs]

// detectAsyncDocumentURI performs Optical Character Recognition (OCR) on a
// PDF file stored in GCS.
func detectAsyncDocumentURI(w io.Writer, gcsSourceURI, gcsDestinationURI string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	request := &visionpb.AsyncBatchAnnotateFilesRequest{
		Requests: []*visionpb.AsyncAnnotateFileRequest{
			{
				Features: []*visionpb.Feature{
					{
						Type: visionpb.Feature_DOCUMENT_TEXT_DETECTION,
					},
				},
				InputConfig: &visionpb.InputConfig{
					GcsSource: &visionpb.GcsSource{Uri: gcsSourceURI},
					// Supported MimeTypes are: "application/pdf" and "image/tiff".
					MimeType: "application/pdf",
				},
				OutputConfig: &visionpb.OutputConfig{
					GcsDestination: &visionpb.GcsDestination{Uri: gcsDestinationURI},
					// How many pages should be grouped into each json output file.
					BatchSize: 2,
				},
			},
		},
	}

	operation, err := client.AsyncBatchAnnotateFiles(ctx, request)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Waiting for the operation to finish.")

	resp, err := operation.Wait(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "%v", resp)

	return nil
}

// [END vision_text_detection_pdf_gcs]

// [START vision_localize_objects_gcs]

// localizeObjects gets objects and bounding boxes from the Vision API for an image at the given file path.
func localizeObjectsURI(w io.Writer, file string) error {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}

	image := vision.NewImageFromURI(file)
	annotations, err := client.LocalizeObjects(ctx, image, nil)
	if err != nil {
		return err
	}

	if len(annotations) == 0 {
		fmt.Fprintln(w, "No objects found.")
		return nil
	}

	fmt.Fprintln(w, "Objects:")
	for _, annotation := range annotations {
		fmt.Fprintln(w, annotation.Name)
		fmt.Fprintln(w, annotation.Score)

		for _, v := range annotation.BoundingPoly.NormalizedVertices {
			fmt.Fprintf(w, "(%f,%f)\n", v.X, v.Y)
		}
	}

	return nil
}

// [END vision_localize_objects_gcs]
