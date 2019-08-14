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

// Package ocr contains Go samples for OCR functions.
package ocr

// [START functions_ocr_setup]
import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	pubsub "cloud.google.com/go/pubsub/apiv1"
	"cloud.google.com/go/storage"
	"cloud.google.com/go/translate"
	vision "cloud.google.com/go/vision/apiv1"
	"golang.org/x/text/language"
)

type config struct {
	ResultTopic    string   `json:"RESULT_TOPIC"`
	ResultBucket   string   `json:"RESULT_BUCKET"`
	TranslateTopic string   `json:"TRANSLATE_TOPIC"`
	Translate      bool     `json:"TRANSLATE"`
	ToLang         []string `json:"TO_LANG"`
}

type ocrmessage struct {
	Text     string       `json:"text"`
	FileName string       `json:"fileName"`
	Lang     language.Tag `json:"lang"`
	SrcLang  language.Tag `json:"srcLang"`
}

func setup() error {
	ctx := context.Background()
	projectID := "GCP_PROJECT"

	visionClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return fmt.Errorf("vision.NewImageAnnotatorClient: %v", err)
	}

	translateClient, err := translate.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("translate.NewClient: %v", err)
	}

	publisher, err := pubsub.NewPublisherClient(ctx)
	if err != nil {
		return fmt.Errorf("translate.NewClient: %v", err)
	}

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile: %v", err)
	}
	config, err := json.Marshal(data)

	// [END functions_ocr_setup]

	_ = visionClient
	_ = translateClient
	_ = publisher
	_ = storageClient
	_ = projectID
	_ = config
	return nil
}
