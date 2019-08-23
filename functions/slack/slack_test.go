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
var slackURL string

// TestMain sets up the config rather than using the config file
// which contains placeholder values.
func TestMain(m *testing.M) {
	ctx := context.Background()
	projectID = os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectID == "" {
		log.Print("GOLANG_SAMPLES_PROJECT_ID is unset. Skipping.")
		return
	}
	slackURL = os.Getenv("GOLANG_SAMPLES_SLACK_URL")
	if projectID == "" {
		log.Print("GOLANG_SAMPLES_SLACK_URL is unset. Skipping.")
		return
	}
	config = &configuration{
		ProjectID: projectID,
		Token:     os.Getenv("GOLANG_SAMPLES_SLACK_TOKEN"),
		Key:       os.Getenv("GOLANG_SAMPLES_KG_KEY"),
	}
	if config.Token == "" {
		log.Print("GOLANG_SAMPLES_SLACK_TOKEN is unset. Skipping.")
		return
	}
	if config.Key == "" {
		log.Print("GOLANG_SAMPLES_KG_KEY is unset. Skipping.")
		return
	}
	kgService, err := kgsearch.NewService(ctx, option.WithAPIKey(config.Key))
	if err != nil {
		log.Fatalf("kgsearch.NewClient: %v", err)
	}
	entitiesService = kgsearch.NewEntitiesService(kgService)

	os.Exit(m.Run())
}

func TestVerifyWebHook(t *testing.T) {
	tests := []struct {
		token   string
		wantErr bool
	}{
		{
			token:   config.Token,
			wantErr: false,
		},
		{
			token:   "this is not the token",
			wantErr: true,
		},
		{
			token:   "",
			wantErr: true,
		},
	}
	for _, test := range tests {
		v := make(url.Values)
		v.Set("token", test.token)
		err := verifyWebHook(v)
		if test.wantErr && err == nil {
			t.Errorf("verifyWebHook(%v) got no error, expected error", test.token)
		}
		if !test.wantErr && err != nil {
			t.Errorf("verifyWebHook(%v) got %v, want no error", test.token, err)
		}
	}
}

func TestFormatSlackMessage(t *testing.T) {
	tests := []struct {
		query string
		want  string
	}{
		{
			query: "Google",
			want:  "Google",
		},
		{
			query: "qwoiuqblaksdfj",
			want:  "No results",
		},
	}
	for _, test := range tests {
		res, err := entitiesService.Search().Query(test.query).Limit(1).Do()
		if err != nil {
			t.Errorf("Do: %v", err)
		}
		msg, err := formatSlackMessage(test.query, res)
		if err != nil {
			t.Errorf("formatSlackMessage: %v", err)
		}
		got := msg.Attachments[0].Text
		if !strings.Contains(got, test.want) {
			t.Errorf("formatSlackMessage(%q) got %q, want %q", test.query, got, test.want)
		}
	}
}

func TestMakeSearchRequest(t *testing.T) {
	query := "Google"
	want := "Google"
	msg, err := makeSearchRequest(query)
	if err != nil {
		t.Errorf("makeSearchRequest: %v", err)
	}
	if msg == nil {
		t.Errorf("empty message from query %q", query)
	}
	got := msg.Text
	if !strings.Contains(got, want) {
		t.Errorf("makeSearchRequest(%q) got %q, want %q", query, got, want)
	}
	if len(msg.Attachments) == 0 {
		t.Errorf("makeSearchRequest(%q) returned no attachments", query)
	}
	got = msg.Attachments[0].Text
	if !strings.Contains(got, want) {
		t.Errorf("makeSearchRequest(%q) got %q, want %q", query, got, want)
	}
}

func TestKGSearch(t *testing.T) {
	w := httptest.NewRecorder()
	form := url.Values{
		"token": []string{config.Token},
		"text":  []string{"Google"},
	}
	req := httptest.NewRequest("POST", slackURL, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	KGSearch(w, req)
	got := w.Body.String()
	if want := "Google"; !strings.Contains(got, want) {
		t.Errorf("KGSearch(%q) got %q, want %q", "Google", got, want)
	}
}
