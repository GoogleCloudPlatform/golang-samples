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
	"fmt"
	"log"

	"github.com/cloudevents/sdk-go/v2/event"
)

// AuditLogEntry represents a LogEntry as described at
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
type AuditLogEntry struct {
	ProtoPayload *AuditLogProtoPayload `json:"protoPayload"`
}

// AuditLogProtoPayload represents AuditLog within the LogEntry.protoPayload
// See https://cloud.google.com/logging/docs/reference/audit/auditlog/rest/Shared.Types/AuditLog
type AuditLogProtoPayload struct {
	MethodName      string                 `json:"methodName"`
	ResourceName    string                 `json:"resourceName"`
	Request         map[string]interface{} `json:"request"`
	RequestMetadata map[string]interface{} `json:"requestMetadata"`
}

// HelloAuditLog receives a Cloud Audit Log event, and logs a few fields.
func HelloAuditLog(ctx context.Context, e event.Event) error {
	log.Printf("Event Type: %s", e.Type())
	log.Printf("Subject: %s", e.Subject())

	logentry := &AuditLogEntry{}

	if err := e.DataAs(logentry); err != nil {
		ferr := fmt.Errorf("event.DataAs: %v", err)
		log.Print(ferr)
		return ferr
	}
	log.Printf("Method Name: %s", logentry.ProtoPayload.MethodName)
	log.Printf("Resource Name: %s", logentry.ProtoPayload.ResourceName)
	if v, ok := logentry.ProtoPayload.Request["@type"]; ok {
		log.Printf("Request Type: %s", v)
	}
	if v, ok := logentry.ProtoPayload.RequestMetadata["callerIp"]; ok {
		log.Printf("Caller IP: %s", v)
	}
	if v, ok := logentry.ProtoPayload.RequestMetadata["callerSuppliedUserAgent"]; ok {
		log.Printf("User Agent: %s", v)
	}
	return nil
}

// [END functions_log_cloudevent]
