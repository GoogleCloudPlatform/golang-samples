// Copyright 2018, Google, LLC.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START functions_ocr_detect]

package ocr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
	"golang.org/x/text/language"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// detectText detects the text in an image using the Google Vision API.
func detectText(w io.Writer, projectID, bucketName, fileName string) error {
	fmt.Fprintf(w, "Looking for text in image %v", fileName)
	ctx := context.Background()
	maxResults := 1
	image := &visionpb.Image{
		Source: &visionpb.ImageSource{
			GcsImageUri: fmt.Sprintf("gs://%s/%s", bucketName, fileName),
		},
	}
	annotations, err := visionClient.DetectTexts(ctx, image, &visionpb.ImageContext{}, maxResults)
	if err != nil {
		return fmt.Errorf("DetectTexts: %v", err)
	}
	text := ""
	if len(annotations) > 0 {
		text = annotations[0].Description
	}
	if len(annotations) == 0 || len(text) == 0 {
		fmt.Fprintf(w, "No text detected in image %q. Returning early.", fileName)
		return nil
	}
	fmt.Fprintf(w, "Extracted text %q from image (%d chars).", text, len(text))

	detectResponse, err := translateClient.DetectLanguage(ctx, []string{text})
	if err != nil {
		return fmt.Errorf("DetectLanguage: %v", err)
	}
	if len(detectResponse) == 0 || len(detectResponse[0]) == 0 {
		return fmt.Errorf("DetectLanguage gave empty response")
	}
	srcLang := detectResponse[0][0].Language.String()
	fmt.Fprintf(w, "Detected language %q for text %q.", srcLang, text)

	// Submit a message to the bus for each target language
	for _, targetLang := range config.ToLang {
		topicName := config.TranslateTopic
		if srcLang == targetLang || srcLang == "und" {
			topicName = config.ResultTopic
		}
		targetTag, err := language.Parse(targetLang)
		if err != nil {
			return fmt.Errorf("language.Parse: %v", err)
		}
		srcTag, err := language.Parse(srcLang)
		if err != nil {
			return fmt.Errorf("language.Parse: %v", err)
		}
		message, err := json.Marshal(ocrMessage{
			Text:     text,
			FileName: fileName,
			Lang:     targetTag,
			SrcLang:  srcTag,
		})
		if err != nil {
			return fmt.Errorf("json.Marshal: %v", err)
		}
		topic := publisher.Topic(topicName)
		ok, err := topic.Exists(ctx)
		if err != nil {
			return fmt.Errorf("Exists: %v", err)
		}
		if !ok {
			return fmt.Errorf("topic %q does not exist", topicName)
		}
		r := topic.Publish(ctx,
			&pubsub.Message{
				Data: []byte(message),
			})
		_, err = r.Get(ctx)
		if err != nil {
			return fmt.Errorf("Get: %v", err)
		}
	}
	return nil
}

// [END functions_ocr_detect]
