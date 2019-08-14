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
	resultTopic    string
	resultBucket   string
	translateTopic string
	translate      bool
	toLang         []string
}

type ocrmessage struct {
	text     string
	fileName string
	lang     language.Tag
	srcLang  language.Tag
}

func setup() error {
	ctx := context.Background()
	visionClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return fmt.Errorf("vision.NewImageAnnotatorClient: %v", err)
	}
	_ = visionClient
	translateClient, err := translate.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("translate.NewClient: %v", err)
	}
	_ = translateClient
	publisher, err := pubsub.NewPublisherClient(ctx)
	if err != nil {
		return fmt.Errorf("translate.NewClient: %v", err)
	}
	_ = publisher
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	_ = storageClient
	projectID := "GCP_PROJECT"
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile: %v", err)
	}
	_ = projectID
	config, err := json.Marshal(data)
	_ = config
	return nil
}

// [END functions_ocr_setup]
