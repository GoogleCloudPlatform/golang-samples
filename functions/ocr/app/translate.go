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

// [START functions_ocr_translate]

package ocr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/translate"
)

// translateText is executed when a message is published to the Cloud Pub/Sub topic specified
// by TRANSLATE_TOPIC in config.json, and translates the text using the Google Translate API.
func translateText(w io.Writer, projectID string, event pubsub.Message) error {
	ctx := context.Background()
	if event.Data == nil {
		return fmt.Errorf("Empty data")
	}
	var message ocrMessage
	if event.Data != nil {
		messageData := event.Data
		err := json.Unmarshal(messageData, &message)
		if err != nil {
			return fmt.Errorf("json.Unmarshal: %v", err)
		}
	} else {
		return fmt.Errorf("Empty data")
	}

	text := message.Text
	fileName := message.FileName
	targetTag := message.Lang
	srcTag := message.SrcLang

	fmt.Fprintf(w, "Translating text into %s.", targetTag.String())
	translateResponse, err := translateClient.Translate(ctx, []string{text}, targetTag,
		&translate.Options{
			Source: srcTag,
		})
	if err != nil {
		return fmt.Errorf("Translate: %v", err)
	}
	if len(translateResponse) == 0 {
		return fmt.Errorf("Empty Translate response")
	}
	translatedText := translateResponse[0]

	topicName := config.ResultTopic
	if err != nil {
		return fmt.Errorf("language.Parse: %v", err)
	}
	messageData, err := json.Marshal(ocrMessage{
		Text:     translatedText.Text,
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
			Data: messageData,
		})
	_, err = r.Get(ctx)
	if err != nil {
		return fmt.Errorf("Get: %v", err)
	}
	fmt.Fprintf(w, "Sent translation: %q", translatedText.Text)
	return nil
}

// [END functions_ocr_translate]
