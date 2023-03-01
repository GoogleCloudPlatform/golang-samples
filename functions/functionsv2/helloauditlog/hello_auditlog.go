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
	"log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/google-cloudevents-go/cloud/auditdata"
	"google.golang.org/protobuf/encoding/protojson"
)

func init() {
	functions.CloudEvent("HelloAuditLog", helloAuditLog)
}

// helloAuditLog receives a CloudEvent containing an AuditLogEntry, and logs a few fields.
func helloAuditLog(ctx context.Context, e event.Event) error {
	// Print out details from the CloudEvent itself
	// See https://github.com/cloudevents/spec/blob/v1.0.1/spec.md#subject
	// for details on the Subject field
	log.Printf("Event Type: %s", e.Type())
	log.Printf("Subject: %s", e.Subject())

	// Decode the Cloud Audit Logging message embedded in the CloudEvent
	logentry := &auditdata.LogEntryData{}
	err := protojson.Unmarshal(e.Data(), logentry)
	if err != nil {
		log.Print(err)
		return err
	}
	// Print out some of the information contained in the Cloud Audit Logging event
	// See https://cloud.google.com/logging/docs/audit#audit_log_entry_structure
	// for a full description of available fields.
	log.Printf("API Method: %s", logentry.ProtoPayload.MethodName)
	log.Printf("Resource Name: %s", logentry.ProtoPayload.ResourceName)
	log.Printf("Principal: %s", logentry.ProtoPayload.AuthenticationInfo.PrincipalEmail)
	return nil
}

// [END functions_log_cloudevent]
