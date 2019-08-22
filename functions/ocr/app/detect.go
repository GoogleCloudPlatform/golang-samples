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

// [START functions_ocr_detect]

package ocr

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"golang.org/x/text/language"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

// detectText detects the text in an image using the Google Vision API.
func detectText(ctx context.Context, bucketName, fileName string) error {
	log.Printf("Looking for text in image %v", fileName)
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
		log.Printf("No text detected in image %q. Returning early.", fileName)
		return nil
	}
	log.Printf("Extracted text %q from image (%d chars).", text, len(text))

	detectResponse, err := translateClient.DetectLanguage(ctx, []string{text})
	if err != nil {
		return fmt.Errorf("DetectLanguage: %v", err)
	}
	if len(detectResponse) == 0 || len(detectResponse[0]) == 0 {
		return fmt.Errorf("DetectLanguage gave empty response")
	}
	srcLang := detectResponse[0][0].Language.String()
	log.Printf("Detected language %q for text %q.", srcLang, text)

	// Submit a message to the bus for each target language
	for _, targetLang := range config.ToLang {
		topicName := config.TranslateTopic
		if srcLang == targetLang || srcLang == "und" { // detection returns "und" for undefined language
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
		topic := pubsubClient.Topic(topicName)
		ok, err := topic.Exists(ctx)
		if err != nil {
			return fmt.Errorf("Exists: %v", err)
		}
		if !ok {
			topic, err = pubsubClient.CreateTopic(ctx, topicName)
			if err != nil {
				return fmt.Errorf("CreateTopic: %v", err)
			}
		}
		msg := &pubsub.Message{
			Data: []byte(message),
		}
		if _, err = topic.Publish(ctx, msg).Get(ctx); err != nil {
			return fmt.Errorf("Get: %v", err)
		}
	}
	return nil
}

// [END functions_ocr_detect]
