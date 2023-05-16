// Copyright 2022 Google LLC
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

// [START functions_ocr_translate]

package ocr

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/translate"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	functions.CloudEvent("translate-text", TranslateText)
}

// TranslateText is executed when a message is published to the Cloud Pub/Sub
// topic specified by the TRANSLATE_TOPIC environment variable, and translates
// the text using the Google Translate API.
func TranslateText(ctx context.Context, cloudevent event.Event) error {
	var event MessagePublishedData
	if err := setup(ctx); err != nil {
		return fmt.Errorf("setup: %w", err)
	}
	if err := cloudevent.DataAs(&event); err != nil {
		return fmt.Errorf("Failed to parse CloudEvent data: %w", err)
	}
	if event.Message.Data == nil {
		log.Printf("event: %s", event)
		return fmt.Errorf("empty data")
	}
	var message ocrMessage
	if err := json.Unmarshal(event.Message.Data, &message); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	log.Printf("Translating text into %s.", message.Lang.String())
	opts := translate.Options{
		Source: message.SrcLang,
	}
	translateResponse, err := translateClient.Translate(ctx, []string{message.Text}, message.Lang, &opts)
	if err != nil {
		return fmt.Errorf("Translate: %w", err)
	}
	if len(translateResponse) == 0 {
		return fmt.Errorf("Empty Translate response")
	}
	translatedText := translateResponse[0]

	messageData, err := json.Marshal(ocrMessage{
		Text:     translatedText.Text,
		FileName: message.FileName,
		Lang:     message.Lang,
		SrcLang:  message.SrcLang,
	})
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	topic := pubsubClient.Topic(resultTopic)
	ok, err := topic.Exists(ctx)
	if err != nil {
		return fmt.Errorf("Exists: %w", err)
	}
	if !ok {
		topic, err = pubsubClient.CreateTopic(ctx, resultTopic)
		if err != nil {
			return fmt.Errorf("CreateTopic: %w", err)
		}
	}
	msg := &pubsub.Message{
		Data: messageData,
	}
	if _, err = topic.Publish(ctx, msg).Get(ctx); err != nil {
		return fmt.Errorf("Get: %w", err)
	}
	log.Printf("Sent translation: %q", translatedText.Text)
	return nil
}

// [END functions_ocr_translate]
