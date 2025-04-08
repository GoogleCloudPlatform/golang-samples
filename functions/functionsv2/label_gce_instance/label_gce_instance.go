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

// [START functions_label_gce_instance]

// Package helloworld provides a set of Cloud Functions samples.
package helloworld

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/protobuf/proto"
)

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

var client *compute.InstancesClient

func init() {
	// Create an Instances Client
	var err error
	client, err = compute.NewInstancesRESTClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create instances client: %s", err)
	}

	functions.CloudEvent("label-gce-instance", labelGceInstance)
}

// Cloud Function that receives GCE instance creation Audit Logs, and adds a
// `creator` label to the instance.
func labelGceInstance(ctx context.Context, ev event.Event) error {
	// Extract parameters from the Cloud Event and Cloud Audit Log data
	logentry := &AuditLogEntry{}
	if err := ev.DataAs(logentry); err != nil {
		err = fmt.Errorf("event.DataAs() : %w", err)
		log.Printf("Error parsing proto payload: %s", err)
		return err
	}
	payload := logentry.ProtoPayload
	creator, ok := payload.AuthenticationInfo["principalEmail"]
	if !ok {
		err := fmt.Errorf("principalEmail not found in cloud event payload: %v", payload)
		log.Printf("creator email not found: %s", err)
		return err
	}

	// Get relevant VM instance details from the event's `subject` property
	// Subject format:
	// compute.googleapis.com/projects/<PROJECT>/zones/<ZONE>/instances/<INSTANCE>
	paths := strings.Split(ev.Subject(), "/")
	if len(paths) < 6 {
		return fmt.Errorf("invalid event subject: %s", ev.Subject())
	}
	project := paths[2]
	zone := paths[4]
	instance := paths[6]

	// Sanitize the `creator` label value to match GCE label requirements
	// See https://cloud.google.com/compute/docs/labeling-resources#requirements
	labelSanitizer := regexp.MustCompile("[^a-z0-9_-]+")
	creatorstring := labelSanitizer.ReplaceAllString(strings.ToLower(creator.(string)), "_")

	// Get the newly-created VM instance's label fingerprint
	// This is a requirement of the Compute Engine API and avoids duplicate labels
	inst, err := client.Get(ctx, &computepb.GetInstanceRequest{
		Project:  project,
		Zone:     zone,
		Instance: instance,
	})
	if err != nil {
		err = fmt.Errorf("could not retrieve GCE instance: %s", err)
		log.Print(err)
		return err
	}
	if v, ok := inst.Labels["creator"]; ok {
		// Instance already has a creator label.
		log.Printf("instance %s already labeled with creator: %s", instance, v)
		return nil
	}

	// Add the creator label to the instance
	op, err := client.SetLabels(ctx, &computepb.SetLabelsInstanceRequest{
		Project:  project,
		Zone:     zone,
		Instance: instance,
		InstancesSetLabelsRequestResource: &computepb.InstancesSetLabelsRequest{
			LabelFingerprint: proto.String(inst.GetLabelFingerprint()),
			Labels: map[string]string{
				"creator": creatorstring,
			},
		},
	})
	if err != nil {
		log.Fatalf("Could not label GCE instance: %s", err)
	}
	log.Printf("Creator label added to %s in operation %v", instance, op)
	return nil
}

// [END functions_label_gce_instance]
