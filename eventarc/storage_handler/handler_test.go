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

package main

// [START eventarc_testing_cloudevent]

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestHelloStorage(t *testing.T) {
	so := storagedata.StorageObjectData{
		Bucket:  "example-bucket",
		Name:    "example-object",
		Updated: timestamppb.New(time.Now()),
	}
	jsondata, err := protojson.Marshal(&so)
	if err != nil {
		t.Fatalf("protojson.Marshal: %v", err)
	}

	ce := cloudevents.NewEvent()
	ce.SetID("sample-id")
	ce.SetSource("//sample/source")
	ce.SetType("google.cloud.storage.object.v1.finalized")
	ce.SetData(*cloudevents.StringOfApplicationJSON(), jsondata)

	w := httptest.NewRecorder()
	r, err := cloudevents.NewHTTPRequestFromEvent(context.Background(), "http://localhost", ce)
	if err != nil {
		t.Fatalf("cloudevents.NewHTTPRequestFromEvent: %v", err)
	}
	HelloStorage(w, r)

	if got := w.Result().StatusCode; got != 200 {
		t.Errorf("got %q, want contained %q", got, 200)
	}

	want := path.Join(so.Bucket, so.Name)
	if got := w.Body.String(); !strings.Contains(got, want) {
		t.Errorf("got %q, want contained %q", got, want)
	}
}

// [END eventarc_testing_cloudevent]

func TestHelloStorage_NotCloudEvent(t *testing.T) {
	so := storagedata.StorageObjectData{
		Bucket:  "example-bucket",
		Name:    "example-object",
		Updated: timestamppb.New(time.Now()),
	}
	jsondata, err := protojson.Marshal(&so)
	if err != nil {
		t.Fatalf("protojson.Marshal: %v", err)
	}

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "http://localhost", bytes.NewReader(jsondata))
	if err != nil {
		t.Fatalf("http.NewRequest: %v", err)
	}
	HelloStorage(w, r)

	wantStatus := http.StatusBadRequest
	if got := w.Result().StatusCode; got != wantStatus {
		t.Errorf("got %q, want %q", got, wantStatus)
	}
	// Ensure failure on malformed cloudevent.
	wantBody := "Bad Request: expected CloudEvent\n"
	if got := w.Body.String(); got != wantBody {
		t.Errorf("got %q, want %q", got, wantBody)
	}
}
