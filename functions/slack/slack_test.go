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
	"encoding/hex"
	"fmt"
	"log"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

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
		Secret:    os.Getenv("GOLANG_SAMPLES_SLACK_SECRET"),
		Key:       os.Getenv("GOLANG_SAMPLES_KG_KEY"),
	}
	if config.Secret == "" {
		log.Print("GOLANG_SAMPLES_SLACK_SECRET is unset. Skipping.")
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
		"text": []string{"Google"},
	}

	secret := config.Secret
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	body := form.Encode()
	base := fmt.Sprintf("v0:%s:%s", ts, body)
	correctSHA2Signature := fmt.Sprintf("v0=%s", hex.EncodeToString(getSignature([]byte(base), []byte(secret))))

	req := httptest.NewRequest("POST", slackURL, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Slack-Request-Timestamp", ts)
	req.Header.Add("X-Slack-Signature", correctSHA2Signature)

	KGSearch(w, req)
	got := w.Body.String()
	if want := "Google"; !strings.Contains(got, want) {
		t.Errorf("KGSearch(%q) got %q, want %q", "Google", got, want)
	}
}
