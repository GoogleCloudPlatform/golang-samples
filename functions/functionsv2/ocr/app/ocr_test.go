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

package ocr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"cloud.google.com/go/translate"
	vision "cloud.google.com/go/vision/apiv1"
	"github.com/cloudevents/sdk-go/v2/event"
	"golang.org/x/text/language"
)

const (
	menuName = "menu.jpg"
	signName = "sign.png"
)

var (
	imageBucketName string
)

func setupTests(t *testing.T) {
	ctx := context.Background()
	projectID = os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectID == "" {
		t.Skip("GOLANG_SAMPLES_PROJECT_ID is unset")
	}
	resultBucket = fmt.Sprintf("%s-result", projectID)
	os.Setenv("GOOGLE_CLOUD_PROJECT", projectID)
	os.Setenv("RESULT_BUCKET", resultBucket)
	os.Setenv("RESULT_TOPIC", "test-result-topic")
	os.Setenv("TO_LANG", "en,fr,es,ja,ru")
	os.Setenv("TRANSLATE_TOPIC", "test-translate-topic")

	imageBucketName = "cloud-samples-data/functions"

	var err error // Prevent shadowing clients with :=.
	visionClient, err = vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		t.Fatalf("vision.NewImageAnnotatorClient: %v", err)
	}

	translateClient, err = translate.NewClient(ctx)
	if err != nil {
		t.Fatalf("translate.NewClient: %v", err)
	}

	pubsubClient, err = pubsub.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("translate.NewClient: %v", err)
	}

	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}

	if _, err := storageClient.Bucket(resultBucket).Attrs(ctx); err != nil {
		t.Skipf("Could not get bucket %v: %v", resultBucket, err)
	}
}

func TestSaveResult(t *testing.T) {
	setupTests(t)
	ctx := context.Background()

	// Create sample data.
	en, err := language.Parse("en")
	if err != nil {
		t.Errorf("language.Parse: %v", err)
	}
	fr, err := language.Parse("fr")
	if err != nil {
		t.Errorf("language.Parse: %v", err)
	}
	data, err := json.Marshal(ocrMessage{
		Text:     "Hello",
		FileName: menuName,
		Lang:     en,
		SrcLang:  fr,
	})

	// Save data.
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	msg := MessagePublishedData{
		Message: PubSubMessage{
			Data: data,
		},
	}
	ce := event.New()
	ce.SetData(*event.StringOfApplicationJSON(), msg)
	if err = SaveResult(ctx, ce); err != nil {
		t.Errorf("SaveResult: %v", err)
	}

	// Check for saved object.
	r, err := storageClient.Bucket(os.Getenv("RESULT_BUCKET")).Object(fmt.Sprintf("%s_%s.txt", menuName, en)).NewReader(ctx)
	if err != nil {
		t.Errorf("NewReader: %v", err)
	}
	resp, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Reader: %v", err)
	}
	got := string(resp)
	if want := "Hello"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTranslateText(t *testing.T) {
	setupTests(t)
	ctx := context.Background()

	// Create data.
	en, err := language.Parse("en")
	if err != nil {
		t.Errorf("language.Parse: %v", err)
	}
	fr, err := language.Parse("fr")
	if err != nil {
		t.Errorf("language.Parse: %v", err)
	}
	data, err := json.Marshal(ocrMessage{
		Text:     "Thanks",
		FileName: menuName,
		Lang:     fr,
		SrcLang:  en,
	})
	if err != nil {
		t.Errorf("json.Marshal: %v", err)
	}

	// Translate data.
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	ce := event.New()
	ce.SetData(*event.StringOfApplicationJSON(), &MessagePublishedData{
		Message: PubSubMessage{
			Data: data,
		},
	})
	err = TranslateText(ctx, ce)
	if err != nil {
		t.Errorf("translateText: %v", err)
	}
	got := buf.String()
	if want := "Merci"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDetectText(t *testing.T) {
	setupTests(t)
	ctx := context.Background()

	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	if err := detectText(ctx, imageBucketName, menuName); err != nil {
		t.Errorf("TestDetectText: %v", err)
	}
	got := buf.String()
	if want := "Filets"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
