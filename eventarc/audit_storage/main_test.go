// Copyright 2020 Google LLC
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

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	cloudevent "github.com/cloudevents/sdk-go/v2"
)

func TestHelloEventsStorage(t *testing.T) {
	event := cloudevent.NewEvent("1.0")
	event.SetID("1")
	event.SetSource("test")
	event.SetSubject("storage.googleapis.com/projects/_/buckets/my-bucket")
	event.SetType("test")

	req, err := cloudevent.NewHTTPRequestFromEvent(context.Background(), "http://example.com", event)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	HelloEventsStorage(rr, req)

	want := "buckets/my-bucket"
	if !strings.Contains(rr.Body.String(), want) {
		t.Errorf("want Body to contain %s, got %s", want, rr.Body)
	}
}
