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

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	functions.CloudEvent("HelloAuditLog", helloAuditLog)
}

// AuditLogEntry represents a LogEntry as described at
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
type AuditLogEntry struct {
	ProtoPayload *AuditLogProtoPayload `json:"protoPayload"`
}

// AuditLogProtoPayload represents AuditLog within the LogEntry.protoPayload
// See https://cloud.google.com/logging/docs/reference/audit/auditlog/rest/Shared.Types/AuditLog
type AuditLogProtoPayload struct {
	MethodName         string                 `json:"methodName"`
	ResourceName       string                 `json:"resourceName"`
	AuthenticationInfo map[string]interface{} `json:"authenticationInfo"`
}

// helloAuditLog receives a CloudEvent containing an AuditLogEntry, and logs a few fields.
func helloAuditLog(ctx context.Context, e event.Event) error {
	// Print out details from the CloudEvent itself
	// See https://github.com/cloudevents/spec/blob/v1.0.1/spec.md#subject
	// for details on the Subject field
	log.Printf("Event Type: %s", e.Type())
	log.Printf("Subject: %s", e.Subject())

	// Decode the Cloud Audit Logging message embedded in the CloudEvent
	logentry := &AuditLogEntry{}
	if err := e.DataAs(logentry); err != nil {
		ferr := fmt.Errorf("event.DataAs: %w", err)
		log.Print(ferr)
		return ferr
	}
	// Print out some of the information contained in the Cloud Audit Logging event
	// See https://cloud.google.com/logging/docs/audit#audit_log_entry_structure
	// for a full description of available fields.
	log.Printf("API Method: %s", logentry.ProtoPayload.MethodName)
	log.Printf("Resource Name: %s", logentry.ProtoPayload.ResourceName)
	if v, ok := logentry.ProtoPayload.AuthenticationInfo["principalEmail"]; ok {
		log.Printf("Principal: %s", v)
	}
	return nil
}

// [END functions_log_cloudevent]
