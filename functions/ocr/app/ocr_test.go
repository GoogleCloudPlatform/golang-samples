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
	"golang.org/x/text/language"
)

const (
	menuName = "menu.jpg"
	signName = "sign.png"
)

var (
	projectID       string
	bucketName      string
	imageBucketName string
)

// TestMain sets up the config rather than using the config file
// which contains placeholder values.
func setupTests(t *testing.T) {
	ctx := context.Background()
	projectID = os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectID == "" {
		t.Skip("GOLANG_SAMPLES_PROJECT_ID is unset")
	}
	bucketName = fmt.Sprintf("%s-result", projectID)
	imageBucketName = "cloud-samples-data/functions"
	config = &configuration{
		ProjectID:      projectID,
		ResultTopic:    "test-result-topic",
		ResultBucket:   bucketName,
		TranslateTopic: "test-translate-topic",
		ToLang:         []string{"en", "fr", "es", "ja", "ru"},
	}

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

	if _, err := storageClient.Bucket(bucketName).Attrs(ctx); err != nil {
		t.Skipf("Could not get bucket %v: %v", bucketName, err)
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
	msg := PubSubMessage{
		Data: data,
	}
	if err = SaveResult(ctx, msg); err != nil {
		t.Errorf("SaveResult: %v", err)
	}

	// Check for saved object.
	r, err := storageClient.Bucket(bucketName).Object(fmt.Sprintf("%s_%s.txt", menuName, en)).NewReader(ctx)
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
	err = TranslateText(ctx, PubSubMessage{
		Data: data,
	})
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
	if want := "Filets de Bœuf"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
