// Copyright 2019 Google LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     https://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package slack

import (
	"context"
	"log"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"google.golang.org/api/kgsearch/v1"
	"google.golang.org/api/option"
)

var projectID string

// TestMain sets up the config rather than using the config file
// which contains placeholder values.
func TestMain(m *testing.M) {
	ctx := context.Background()
	projectID = os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectID == "" {
		log.Print("GOLANG_SAMPLES_PROJECT_ID is unset. Skipping.")
		return
	}
	config = &configuration{
		ProjectID: projectID,
		Token:     "cbdkUvSjpoiSfaHDhVNiXOs4",
		Key:       "AIzaSyCsNvMEF7iGUM-Nu6m4ARpQw39Txzl4hJU",
	}
	kgService, err := kgsearch.NewService(ctx, option.WithAPIKey(config.Key))
	if err != nil {
		log.Fatalf("kgsearch.NewClient: %v", err)
	}
	entitiesService = kgsearch.NewEntitiesService(kgService)

	os.Exit(m.Run())
}

func TestVerifyWebHook(t *testing.T) {
	v := make(url.Values)
	v.Set("token", config.Token)
	if err := verifyWebHook(v); err != nil {
		t.Errorf("verifyWebHook: %v", err)
	}
	v = make(url.Values)
	v.Set("token", "this is not the token")
	if err := verifyWebHook(v); err == nil {
		t.Errorf("got %q, want %q", "nil", "invalid request/credentials")
	}
	v = make(url.Values)
	v.Set("token", "")
	if err := verifyWebHook(v); err == nil {
		t.Errorf("got %q, want %q", "nil", "empty form token")
	}
}

func TestFormatSlackMessage(t *testing.T) {
	query := "Google"
	req := entitiesService.Search().Query(query).Limit(1)
	res, err := req.Do()
	if err != nil {
		t.Errorf("Do: %v", err)
	}
	msg, err := formatSlackMessage(query, res)
	if err != nil {
		t.Errorf("formatSlackMessage: %v", err)
	}
	got := msg.Attachments[0].Text
	if want := "Google"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMakeSearchRequest(t *testing.T) {
	query := "Google"
	msg, err := makeSearchRequest(query)
	if err != nil {
		t.Errorf("makeSearchRequest: %v", err)
	}
	if msg == nil {
		t.Errorf("empty message from query %q", query)
	}
	got := msg.Text
	if want := "Google"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
	if len(msg.Attachments) == 0 {
		t.Errorf("no attachments from query %q", query)
	}
	got = msg.Attachments[0].Text
	if want := "Google"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}

	query = "qjwzl"
	msg, err = makeSearchRequest(query)
	if msg == nil {
		t.Errorf("empty message from query %q", query)
	}
	if len(msg.Attachments) == 0 {
		t.Errorf("no attachments from query %q", query)
	}
	got = msg.Attachments[0].Text
	if want := "No results"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestKGSearch(t *testing.T) {
	w := httptest.NewRecorder()
	r := strings.NewReader("Google")
	KGSearch(w, httptest.NewRequest("POST", "", r))
	got := w.Body.String()
	if want := "Google"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
