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
	"strings"
	"testing"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	menuName   = "images/menu.png"
	bucketName = "golang-samples-ocr-test"
	topicName  = "ocr-test-topic"
)

func TestMain(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	client, err := pubsub.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Errorf("pubsub.NewClient: %v", err)
	}
	topic, err := client.CreateTopic(ctx, topicName)
	if err != nil {
		t.Errorf("CreateTopic: %v", err)
	}
	_ = topic
}

func TestSaveResult(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.Bucket(bucketName)
	err = detectText(buf, tc.ProjectID, bucketName, menuName)
	if err != nil {
		t.Errorf("TestInspectFile: %v", err)
	}
	err = translateText(buf, tc.ProjectID)
	if err != nil {
		t.Errorf("TestInspectFile: %v", err)
	}
	got := buf.String()
	if want := "Menu"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTranslateText(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.Bucket(bucketName)
	err = detectText(buf, tc.ProjectID, bucketName, menuName)
	if err != nil {
		t.Errorf("TestInspectFile: %v", err)
	}
	err = translateText(buf, tc.ProjectID)
	if err != nil {
		t.Errorf("TestInspectFile: %v", err)
	}
	got := buf.String()
	if want := "Menu"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDetectText(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.Bucket(bucketName)
	err = detectText(buf, tc.ProjectID, bucketName, menuName)
	if err != nil {
		t.Errorf("TestInspectFile: %v", err)
	}
	got := buf.String()
	if want := "Menu"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
