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

// [START functions_ocr_setup]

// Package ocr contains Go samples for creating OCR
// (Optical Character Recognition) Cloud functions.
package ocr

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"cloud.google.com/go/translate"
	vision "cloud.google.com/go/vision/apiv1"
	"golang.org/x/text/language"
)

type configType struct {
	ProjectID      string   `json:"PROJECT_ID"`
	ResultTopic    string   `json:"RESULT_TOPIC"`
	ResultBucket   string   `json:"RESULT_BUCKET"`
	TranslateTopic string   `json:"TRANSLATE_TOPIC"`
	Translate      bool     `json:"TRANSLATE"`
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
	publisher       *pubsub.Client
	storageClient   *storage.Client
	config          *configType
)

func setup(ctx context.Context) {
	cfgFile, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}

	d := json.NewDecoder(cfgFile)
	config = &configType{}
	err = d.Decode(config)
	if err != nil {
		log.Fatalf("Decode: %v", err)
	}

	projectID := config.ProjectID

	visionClient, err = vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("vision.NewImageAnnotatorClient: %v", err)
	}

	translateClient, err = translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("translate.NewClient: %v", err)
	}

	publisher, err = pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("translate.NewClient: %v", err)
	}

	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
}

// [END functions_ocr_setup]
