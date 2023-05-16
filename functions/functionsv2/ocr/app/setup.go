// Copyright 2023 Google LLC
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

// [START functions_ocr_setup]

// Package ocr contains Go samples for creating OCR
// (Optical Character Recognition) Cloud functions.
package ocr

import (
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"cloud.google.com/go/translate"
	vision "cloud.google.com/go/vision/apiv1"
	"golang.org/x/text/language"
)

type ocrMessage struct {
	Text     string       `json:"text"`
	FileName string       `json:"fileName"`
	Lang     language.Tag `json:"lang"`
	SrcLang  language.Tag `json:"srcLang"`
}

// Eventarc sends a MessagePublishedData object.
// See the documentation for additional fields and more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub_1
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for additional fields and more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

var (
	visionClient    *vision.ImageAnnotatorClient
	translateClient *translate.Client
	pubsubClient    *pubsub.Client
	storageClient   *storage.Client

	projectID      string
	resultBucket   string
	resultTopic    string
	toLang         []string
	translateTopic string
)

func setup(ctx context.Context) error {
	projectID = os.Getenv("GCP_PROJECT")
	resultBucket = os.Getenv("RESULT_BUCKET")
	resultTopic = os.Getenv("RESULT_TOPIC")
	toLang = strings.Split(os.Getenv("TO_LANG"), ",")
	translateTopic = os.Getenv("TRANSLATE_TOPIC")

	var err error // Prevent shadowing clients with :=.

	if visionClient == nil {
		visionClient, err = vision.NewImageAnnotatorClient(ctx)
		if err != nil {
			return fmt.Errorf("vision.NewImageAnnotatorClient: %w", err)
		}
	}

	if translateClient == nil {
		translateClient, err = translate.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("translate.NewClient: %w", err)
		}
	}

	if pubsubClient == nil {
		pubsubClient, err = pubsub.NewClient(ctx, projectID)
		if err != nil {
			return fmt.Errorf("translate.NewClient: %w", err)
		}
	}

	if storageClient == nil {
		storageClient, err = storage.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("storage.NewClient: %w", err)
		}
	}
	return nil
}

// [END functions_ocr_setup]
