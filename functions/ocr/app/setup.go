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

// [START functions_ocr_setup]

// Package ocr contains Go samples for creating OCR
// (Optical Character Recognition) Cloud functions.
package ocr

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"cloud.google.com/go/translate"
	vision "cloud.google.com/go/vision/apiv1"
	"golang.org/x/text/language"
)

type configuration struct {
	ProjectID      string   `json:"PROJECT_ID"`
	ResultTopic    string   `json:"RESULT_TOPIC"`
	ResultBucket   string   `json:"RESULT_BUCKET"`
	TranslateTopic string   `json:"TRANSLATE_TOPIC"`
	ToLang         []string `json:"TO_LANG"`
}

type ocrMessage struct {
	Text     string       `json:"text"`
	FileName string       `json:"fileName"`
	Lang     language.Tag `json:"lang"`
	SrcLang  language.Tag `json:"srcLang"`
}

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Bucket         string    `json:"bucket"`
	Name           string    `json:"name"`
	Metageneration string    `json:"metageneration"`
	ResourceState  string    `json:"resourceState"`
	TimeCreated    time.Time `json:"timeCreated"`
	Updated        time.Time `json:"updated"`
}

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

var (
	visionClient    *vision.ImageAnnotatorClient
	translateClient *translate.Client
	pubsubClient    *pubsub.Client
	storageClient   *storage.Client
	config          *configuration
)

func setup(ctx context.Context) error {
	if config == nil {
		cfgFile, err := os.Open("config.json")
		if err != nil {
			return fmt.Errorf("os.Open: %v", err)
		}

		d := json.NewDecoder(cfgFile)
		config = &configuration{}
		if err = d.Decode(config); err != nil {
			return fmt.Errorf("Decode: %v", err)
		}
	}

	var err error // Prevent shadowing clients with :=.

	if visionClient == nil {
		visionClient, err = vision.NewImageAnnotatorClient(ctx)
		if err != nil {
			return fmt.Errorf("vision.NewImageAnnotatorClient: %v", err)
		}
	}

	if translateClient == nil {
		translateClient, err = translate.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("translate.NewClient: %v", err)
		}
	}

	if pubsubClient == nil {
		pubsubClient, err = pubsub.NewClient(ctx, config.ProjectID)
		if err != nil {
			return fmt.Errorf("translate.NewClient: %v", err)
		}
	}

	if storageClient == nil {
		storageClient, err = storage.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("storage.NewClient: %v", err)
		}
	}
	return nil
}

// [END functions_ocr_setup]
