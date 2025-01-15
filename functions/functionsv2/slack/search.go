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

// [START functions_slack_search]

// Package slack is a Cloud Function which receives a query from
// a Slack command and responds with the KG API result.
package slack

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

const (
	version                     = "v0"
	slackRequestTimestampHeader = "X-Slack-Request-Timestamp"
	slackSignatureHeader        = "X-Slack-Signature"
)

type Attachment struct {
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
	Attachments  []Attachment `json:"attachments"`
}

func init() {
	functions.HTTP("KGSearch", kgSearch)
	setup(context.Background())
}

// kgSearch uses the Knowledge Graph API to search for a query provided
// by a Slack command.
func kgSearch(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		http.Error(w, "could not read request body", http.StatusBadRequest)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are accepted", http.StatusMethodNotAllowed)
		return
	}
	formData, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		log.Printf("Error: Failed to Parse Form: %v", err)
		http.Error(w, "Couldn't parse form", http.StatusBadRequest)
		return
	}

	result, err := verifyWebHook(r, bodyBytes, slackSecret)
	if err != nil || !result {
		log.Printf("verifyWebhook failed: %v", err)
		http.Error(w, "Failed to verify request signature", http.StatusBadRequest)
		return
	}

	if len(formData.Get("text")) == 0 {
		log.Printf("no search text found: %v", formData)
		http.Error(w, "search text was empty", http.StatusBadRequest)
		return
	}
	kgSearchResponse, err := makeSearchRequest(formData.Get("text"))
	if err != nil {
		log.Printf("makeSearchRequest failed: %v", err)
		http.Error(w, "makeSearchRequest failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(kgSearchResponse); err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal results: %v", err), 500)
		return
	}
}

// [END functions_slack_search]

// [START functions_slack_request]
func makeSearchRequest(query string) (*Message, error) {
	res, err := entitiesService.Search().Query(query).Limit(1).Do()
	if err != nil {
		return nil, err
	}
	return formatSlackMessage(query, res)
}

// [END functions_slack_request]

// [START functions_verify_webhook]

// verifyWebHook verifies the request signature.
// See https://api.slack.com/docs/verifying-requests-from-slack.
func verifyWebHook(r *http.Request, body []byte, slackSigningSecret string) (bool, error) {
	timeStamp := r.Header.Get(slackRequestTimestampHeader)
	slackSignature := r.Header.Get(slackSignatureHeader)

	t, err := strconv.ParseInt(timeStamp, 10, 64)
	if err != nil {
		return false, fmt.Errorf("strconv.ParseInt(%s): %w", timeStamp, err)
	}

	if ageOk, age := checkTimestamp(t); !ageOk {
		return false, fmt.Errorf("checkTimestamp(%v): %v %v", t, ageOk, age)
	}

	if timeStamp == "" || slackSignature == "" {
		return false, fmt.Errorf("timeStamp and/or signature headers were blank")
	}

	baseString := fmt.Sprintf("%s:%s:%s", version, timeStamp, body)

	signature := getSignature([]byte(baseString), []byte(slackSigningSecret))

	trimmed := strings.TrimPrefix(slackSignature, fmt.Sprintf("%s=", version))
	signatureInHeader, err := hex.DecodeString(trimmed)

	if err != nil {
		return false, fmt.Errorf("hex.DecodeString(%v): %w", trimmed, err)
	}

	return hmac.Equal(signature, signatureInHeader), nil
}

func getSignature(base []byte, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(base)

	return h.Sum(nil)
}

// checkTimestamp allows requests time stamped less than 5 minutes ago.
func checkTimestamp(timeStamp int64) (bool, time.Duration) {
	t := time.Since(time.Unix(timeStamp, 0))

	return t.Minutes() <= 5, t
}

// [END functions_verify_webhook]
