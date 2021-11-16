// Copyright 2021 Google LLC
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

// [START functions_log_cloudevent]

// Package helloworld provides a set of Cloud Functions samples.
package helloworld

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/cloudevents/sdk-go/v2/event"
	auditevents "github.com/googleapis/google-cloudevents-go/cloud/audit/v1"
	"google.golang.org/protobuf/proto"
)

func makeAuditLog(subject string, payload auditevents.ProtoPayload) (event.Event, error) {
	logevent := auditevents.LogEntryData{
		ProtoPayload: &payload,
	}
	e := event.New()
	e.SetSubject(subject)
	e.SetType("google.cloud.audit.log.v1.written")
	eventdata, err := json.Marshal(logevent)
	if err != nil {
		return event.New(), err
	}
	e.SetDataContentType("application/json")
	e.SetData("application/json", eventdata)
	return e, nil
}

func TestHelloAuditLog(t *testing.T) {

	tests := []struct {
		name         string
		subject      string
		payload      auditevents.ProtoPayload
		expectedLogs []string
	}{
		{"sample-output",
			"storage.googleapis.com/projects/_/buckets/my-bucket/objects/test.txt",
			auditevents.ProtoPayload{
				ResourceName: proto.String("my-resource"),
				Request: map[string]interface{}{
					"@type": "type.googleapis.com/storage.objects.write",
				},
				RequestMetadata: &auditevents.RequestMetadata{
					CallerIP:                proto.String("1.2.3.4"),
					CallerSuppliedUserAgent: proto.String("example-user-agent"),
				},
			},
			[]string{
				"Event Type: google.cloud.audit.log.v1.written",
				"Subject: storage.googleapis.com/projects/_/buckets/my-bucket/objects/test.txt",
				"Resource Name: my-resource",
				"Request Type: type.googleapis.com/storage.objects.write",
				"Caller IP: 1.2.3.4",
				"User Agent: example-user-agent",
			},
		},
	}
	for _, tt := range tests {

		// Capture log output
		r, w, _ := os.Pipe()
		log.SetOutput(w)
		originalFlags := log.Flags()
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

		t.Run(tt.name, func(t *testing.T) {
			event, err := makeAuditLog(tt.subject, tt.payload)
			if err != nil {
				t.Errorf("HelloAuditLog() failed to create audit.LogEntryData: %v", err)
			}
			if err := HelloAuditLog(context.Background(), event); err != nil {
				t.Errorf("HelloAuditLog() unexpected error: %v", err)
			}

			w.Close()
			log.SetOutput(os.Stderr)
			log.SetFlags(originalFlags)

			// check output sent to the logging pipe.
			output, err := ioutil.ReadAll(r)
			if err != nil {
				t.Errorf("Failed reading output from HelloAuditLog(): %v", err)
			}
			for _, l := range tt.expectedLogs {
				if !strings.Contains(string(output), l) {
					t.Errorf("HelloAuditlog() expected log not found: expected '%s', got '%s'", l, output)
				}
			}

		})

	}
}
