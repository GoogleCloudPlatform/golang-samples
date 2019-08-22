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

// [START functions_slack_search]

// Package slack is a Cloud Function which recieves a query from
// a Slack command and responds with the KG API result.
package slack

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type attachment struct {
	Color     string `json:"color"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
	ImageURL  string `json:"image_url"`
}

// Message is the a Slack message event.
// see https://api.slack.com/docs/message-formatting
type Message struct {
	ResponseType string       `json:"response_type"`
	Text         string       `json:"text"`
	Attachments  []attachment `json:"attachments"`
}

// KGSearch uses the Knowledge Graph API to search for a query provided
// by a Slack command.
func KGSearch(w http.ResponseWriter, r *http.Request) {
	setup(r.Context())
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are accepted", 405)
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Couldn't parse form", 400)
		log.Fatalf("ParseForm: %v", err)
	}
	if err := verifyWebHook(r.Form); err != nil {
		log.Fatalf("verifyWebhook: %v", err)
	}
	if len(r.Form["text"]) == 0 {
		log.Fatalf("emtpy text in form")
	}
	kgSearchResponse, err := makeSearchRequest(r.Form["text"][0])
	if err != nil {
		log.Fatalf("makeSearchRequest: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(kgSearchResponse); err != nil {
		log.Fatalf("json.Marshal: %v", err)
	}
}

// [END functions_slack_search]

// [START functions_verify_webhook]
func verifyWebHook(form url.Values) error {
	t := form.Get("token")
	if len(t) == 0 {
		return fmt.Errorf("empty form token")
	}
	if t != config.Token {
		return fmt.Errorf("invalid request/credentials: %q", t[0])
	}
	return nil
}

// [END functions_verify_webhook]

// [START functions_slack_request]
func makeSearchRequest(query string) (*Message, error) {
	res, err := entitiesService.Search().Query(query).Limit(1).Do()
	if err != nil {
		return nil, fmt.Errorf("Do: %v", err)
	}
	return formatSlackMessage(query, res)
}

// [END functions_slack_request]
