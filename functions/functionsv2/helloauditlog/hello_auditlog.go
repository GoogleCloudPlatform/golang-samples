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
	auditevents "github.com/googleapis/google-cloudevents-go/cloud/audit/v1"
)

// HelloAuditLog receives a Cloud Audit Log event, and logs a few fields.
func HelloAuditLog(ctx context.Context, e event.Event) error {
	log.Printf("Event Type: %s", e.Type())
	log.Printf("Subject: %s", e.Subject())

	logentry := &auditevents.LogEntryData{}

	if err := e.DataAs(logentry); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}
	log.Printf("Resource Name: %s", *logentry.ProtoPayload.ResourceName)
	if logentry.ProtoPayload.RequestMetadata.CallerIP != nil {
		log.Printf("Caller IP: %s", *logentry.ProtoPayload.RequestMetadata.CallerIP)
	}
	if logentry.ProtoPayload.RequestMetadata.CallerSuppliedUserAgent != nil {
		log.Printf("User Agent: %s", *logentry.ProtoPayload.RequestMetadata.CallerSuppliedUserAgent)
	}
	return nil
}

// [END functions_log_cloudevent]
