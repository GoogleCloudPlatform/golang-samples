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

package ocr

// [START functions_ocr_detect]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	pubsub "cloud.google.com/go/pubsub/apiv1"
	"cloud.google.com/go/translate"
	vision "cloud.google.com/go/vision/apiv1"
	"golang.org/x/text/language"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
	pubsubpb "google.golang.org/genproto/googleapis/pubsub/v1"
)

// detectText detects the text in an image using the Google Vision API.
func detectText(w io.Writer, projectID, bucketName, fileName string) error {
	// bucketName := "ocr-image-bucket123"
	// fileName := "menu.jpg"
	fmt.Fprintf(w, "Looking for text in image %v", fileName)
	ctx := context.Background()
	visionClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return fmt.Errorf("vision.NewImageAnnotatorClient: %v", err)
	}
	maxResults := 1
	annotations, err := visionClient.DetectTexts(ctx,
		&pb.Image{
			Source: &pb.ImageSource{
				GcsImageUri: fmt.Sprintf("gs://%s/%s", bucketName, fileName),
			},
		},
		&pb.ImageContext{}, maxResults,
	)
	if err != nil {
		return fmt.Errorf("DetectTexts: %v", err)
	}
	text := ""
	if len(annotations) > 0 {
		text = annotations[0].Description
	}
	fmt.Fprintf(w, "Extracted text %q from image (%d chars).", text, len(text))

	translateClient, err := translate.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("translate.NewClient: %v", err)
	}
	detectResponse, err := translateClient.DetectLanguage(ctx, []string{text})
	if err != nil {
		return fmt.Errorf("DetectLanguage: %v", err)
	}
	if len(detectResponse) == 0 || len(detectResponse[0]) == 0 {
		return fmt.Errorf("DetectLanguage gave empty response")
	}
	srcLang := detectResponse[0][0].Language.String()
	fmt.Fprintf(w, "Detected language %q for text %q.", srcLang, text)

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile: %v", err)
	}
	config := &config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	// Submit a message to the bus for each target language
	publisher, err := pubsub.NewPublisherClient(ctx)
	if err != nil {
		return fmt.Errorf("translate.NewClient: %v", err)
	}
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
		message, err := json.Marshal(ocrmessage{
			Text:     text,
			FileName: fileName,
			Lang:     targetTag,
			SrcLang:  srcTag,
		})
		if err != nil {
			return fmt.Errorf("json.Marshal: %v", err)
		}
		messageData := []byte(message)
		topicPath := pubsub.PublisherTopicPath(projectID, topicName)
		_, err = publisher.Publish(ctx,
			&pubsubpb.PublishRequest{
				Topic: topicPath,
				Messages: []*pubsubpb.PubsubMessage{
					&pubsubpb.PubsubMessage{
						Data: messageData,
					},
				},
			})
		if err != nil {
			return fmt.Errorf("Publish: %v", err)
		}
	}
	return nil
}

// [END functions_ocr_detect]
