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
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type oldTimeStampError struct {
	s string
}

func (e *oldTimeStampError) Error() string {
	return e.s
}

const (
	version                     = "v0"
	slackRequestTimestampHeader = "X-Slack-Request-Timestamp"
	slackSignatureHeader        = "X-Slack-Signature"
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

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Couldn't read request body: %v", err)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if r.Method != "POST" {
		http.Error(w, "Only POST requests are accepted", http.StatusMethodNotAllowed)

	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Couldn't parse form", http.StatusBadRequest)
		log.Fatalf("ParseForm: %v", err)
	}

	// Reset r.Body as ParseForm depletes it by reading the io.ReadCloser.
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	result, err := verifyWebHook(r, slackSecret)
	if err != nil {
		log.Fatalf("verifyWebhook: %v", err)
	}
	if !result {
		log.Fatalf("signatures did not match.")
	}

	if len(r.Form["text"]) == 0 {
		log.Fatalf("empty text in form")
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

// [START functions_slack_request]
func makeSearchRequest(query string) (*Message, error) {
	res, err := entitiesService.Search().Query(query).Limit(1).Do()
	if err != nil {
		return nil, fmt.Errorf("Do: %w", err)
	}
	return formatSlackMessage(query, res)
}

// [END functions_slack_request]

// [START functions_verify_webhook]

// verifyWebHook verifies the request signature.
// See https://api.slack.com/docs/verifying-requests-from-slack.
func verifyWebHook(r *http.Request, slackSigningSecret string) (bool, error) {
	timeStamp := r.Header.Get(slackRequestTimestampHeader)
	slackSignature := r.Header.Get(slackSignatureHeader)

	t, err := strconv.ParseInt(timeStamp, 10, 64)
	if err != nil {
		return false, fmt.Errorf("strconv.ParseInt(%s): %w", timeStamp, err)
	}

	if ageOk, age := checkTimestamp(t); !ageOk {
		return false, &oldTimeStampError{fmt.Sprintf("checkTimestamp(%v): %v %v", t, ageOk, age)}
		// return false, fmt.Errorf("checkTimestamp(%v): %v %v", t, ageOk, age)
	}

	if timeStamp == "" || slackSignature == "" {
		return false, fmt.Errorf("either timeStamp or signature headers were blank")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, fmt.Errorf("ioutil.ReadAll(%v): %w", r.Body, err)
	}

	// Reset the body so other calls won't fail.
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

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

// Arbitrarily trusting requests time stamped less than 5 minutes ago.
func checkTimestamp(timeStamp int64) (bool, time.Duration) {
	t := time.Since(time.Unix(timeStamp, 0))

	return t.Minutes() <= 5, t
}

// [END functions_verify_webhook]
