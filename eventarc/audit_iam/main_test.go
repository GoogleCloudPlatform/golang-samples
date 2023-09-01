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

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	cloudevent "github.com/cloudevents/sdk-go/v2"
	"github.com/googleapis/google-cloudevents-go/cloud/auditdata"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestIamAuditlogs(t *testing.T) {
	auditlog := &auditdata.AuditLog{
		ServiceName: "iam.googleapis.com",
		MethodName:  "google.iam.admin.v1.CreateServiceAccountKey",
		AuthenticationInfo: &auditdata.AuthenticationInfo{
			PrincipalEmail: "user@example.com",
		},
		Request: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"name": structpb.NewStringValue("projects/-/serviceAccounts/service-account@my-project.iam.gserviceaccount.com")},
		},
	}
	logentry := &auditdata.LogEntryData{
		LogName:      "",
		ProtoPayload: auditlog,
	}

	event := cloudevent.NewEvent("1.0")
	event.SetID("1")
	event.SetSource("iam.googleapis.com")
	event.SetSubject("subject goes here")
	event.SetType("test")
	event.SetData("application/json", logentry)

	req, err := cloudevent.NewHTTPRequestFromEvent(context.Background(), "http://example.com", event)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	HandleCloudEvent(rr, req)

	want := "New Service Account Key created for projects/-/serviceAccounts/service-account@my-project.iam.gserviceaccount.com by user@example.com"
	if !strings.Contains(rr.Body.String(), want) {
		t.Errorf("want body to contain %s, got %s", want, rr.Body)
	}

}
