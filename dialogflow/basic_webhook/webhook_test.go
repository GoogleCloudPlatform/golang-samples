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

package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleWebhookRequest(t *testing.T) {
	var testCases = []struct {
		name           string
		requestBody    string
		wantBody       string
		wantStatusCode int
	}{
		{
			name: "get-agent-name",
			requestBody: `
				{
					"responseId": "response-id-123",
					"queryResult": {
						"queryText": "what is your name?",
						"parameters": {},
						"allRequiredParamsPresent": true,
						"fulfillmentText": "My name is Dialogflow!",
						"fulfillmentMessages": [
							{
								"text": {
									"text": [
										"My name is Dialogflow!"
									]
								}
							}
						],
						"outputContexts": [
							{
								"name": "projects/project-id-123/...",
								"parameters": {
									"no-input": 0,
									"no-match": 0
								}
							}
						],
						"intent": {
							"name": "projects/project-id-123/agent/intents/intent-id-123",
							"displayName": "get-agent-name"
						},
						"intentDetectionConfidence": 1,
						"languageCode": "en"
					},
					"originalDetectIntentRequest": {
						"payload": {}
					},
					"session": "projects/project-id-123/agent/sessions/session-id-123"
				}`,
			wantBody: `
			{
				"fulfillmentMessages": [
				  {
					  "text": {
						  "text": [
							  "My name is Dialogflow Go Webhook"
							]
						}
					}
				]
			}`,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "welcome",
			requestBody: `
				{
					"responseId": "response-id-123",
					"queryResult": {
						"intent": {
							"displayName": "Default Welcome Intent"
						}
					},
					"session": "projects/project-id-123/agent/sessions/session-id-123"
				}`,
			wantBody: `
			{
				"fulfillmentMessages": [
				  {
					  "text": {
						  "text": [
							  "Welcome from Dialogflow Go Webhook"
							]
						}
					}
				]
			}`,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "unknown",
			requestBody: `
				{
					"responseId": "response-id-123",
					"queryResult": {
						"intent": {
							"displayName": "unknown"
						}
					},
					"session": "projects/project-id-123/agent/sessions/session-id-123"
				}`,
			wantBody:       "",
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "bad-request",
			requestBody:    "===",
			wantBody:       "",
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tc := range testCases {
		req := httptest.NewRequest("POST", "/handler", strings.NewReader(tc.requestBody))
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleWebhookRequest)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != tc.wantStatusCode {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, tc.wantStatusCode)
		}

		// Check the response body is what we want.
		if tc.wantStatusCode == http.StatusOK {
			compactGot := new(bytes.Buffer)
			compactWant := new(bytes.Buffer)
			if err := json.Compact(compactGot, rr.Body.Bytes()); err != nil {
				t.Errorf("failed to compact: %v", err)
			}
			if err := json.Compact(compactWant, []byte(tc.wantBody)); err != nil {
				t.Errorf("failed to compact: %v", err)
			}
			if compactGot.String() != compactWant.String() {
				t.Errorf("%s: handler returned unexpected body, got:\n<%v>\nwant:\n<%v>\n",
					tc.name, compactGot.String(), compactWant.String())
			}
		}
	}
}
