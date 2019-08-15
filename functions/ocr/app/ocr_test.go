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
	"strings"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/text/language"
)

const (
	menuName = "menu.jpg"
	signName = "sign.png"
)

var (
	bucketName      string
	imageBucketName string
)

func TestMain(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName = fmt.Sprintf("%s-result", tc.ProjectID)
	imageBucketName = fmt.Sprintf("%s-image", tc.ProjectID)
}

func TestSaveResult(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)
	bucketName = fmt.Sprintf("%s-result", tc.ProjectID)
	imageBucketName = fmt.Sprintf("%s-image", tc.ProjectID)
	buf := new(bytes.Buffer)
	bkt := storageClient.Bucket(bucketName)
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
	err = saveResult(buf, pubsub.Message{
		Data: data,
	})
	if err != nil {
		t.Errorf("TestSaveResult: %v", err)
	}
	r, err := bkt.Object(fmt.Sprintf("%s_%s.txt", menuName, en)).NewReader(ctx)
	if err != nil {
		t.Errorf("NewReader: %v", err)
	}
	fbuf := make([]byte, 100, 100)
	_, err = r.Read(fbuf)
	if err != nil {
		t.Errorf("Reader: %v", err)
	}
	got := string(fbuf)
	if want := "Hello"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTranslateText(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
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
		Lang:     fr,
		SrcLang:  en,
	})
	if err != nil {
		t.Errorf("json.Marshal: %v", err)
	}
	err = translateText(buf, tc.ProjectID, pubsub.Message{
		Data: data,
	})
	if err != nil {
		t.Errorf("translateText: %v", err)
	}
	got := buf.String()
	if want := "Bonjour"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDetectText(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	storageClient.Bucket(imageBucketName)
	err := detectText(buf, tc.ProjectID, imageBucketName, menuName)
	if err != nil {
		t.Errorf("TestDetectText: %v", err)
	}
	got := buf.String()
	if want := "Filets de Boeuf"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
