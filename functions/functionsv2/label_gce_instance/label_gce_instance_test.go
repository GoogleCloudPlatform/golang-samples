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
	"encoding/json"
	"testing"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/googleapis/gax-go/v2"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

// FakeInstancesClient is a GceInstancesClient that does not call GCE, and returns stub values.
type FakeInstancesClient struct {
	instance *computepb.Instance
	labels   map[string]string
}

func (c FakeInstancesClient) Get(_ context.Context, _ *computepb.GetInstanceRequest, _ ...gax.CallOption) (*computepb.Instance, error) {
	return c.instance, nil
}

func (c FakeInstancesClient) SetLabels(_ context.Context, r *computepb.SetLabelsInstanceRequest, _ ...gax.CallOption) (*compute.Operation, error) {
	lset := r.GetInstancesSetLabelsRequestResource().GetLabels()
	for k, v := range lset {
		c.labels[k] = v
	}
	return new(compute.Operation), nil
}

func TestLabelGceInstance(t *testing.T) {
	tests := []struct {
		name       string
		payload    *AuditLogProtoPayload
		instance   *computepb.Instance
		newCreator string
	}{
		{
			name: "no label",
			payload: &AuditLogProtoPayload{
				MethodName:   "",
				ResourceName: "",
				AuthenticationInfo: map[string]interface{}{
					"principalEmail": "user@example.com",
				},
			},
			instance: &computepb.Instance{
				Hostname:         new(string),
				Id:               new(uint64),
				Kind:             new(string),
				LabelFingerprint: new(string),
				Labels:           map[string]string{},
				Name:             new(string),
				SelfLink:         new(string),
				Zone:             new(string),
			},
			newCreator: "user_example_com",
		},
		{
			name: "existing-creator",
			payload: &AuditLogProtoPayload{
				MethodName:   "",
				ResourceName: "",
				AuthenticationInfo: map[string]interface{}{
					"principalEmail": "user@example.com",
				},
			},
			instance: &computepb.Instance{
				Hostname:         new(string),
				Id:               new(uint64),
				Kind:             new(string),
				LabelFingerprint: new(string),
				Labels: map[string]string{
					"creator": "existing_creator",
				},
				Name:     new(string),
				SelfLink: new(string),
				Zone:     new(string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fake := &FakeInstancesClient{
				instance: tt.instance,
				labels:   make(map[string]string),
			}
			client = fake
			e := event.New()
			e.SetSubject("compute.googleapis.com/projects/PROJECT/zones/ZONE/instances/INSTANCE")
			e.SetType("google.cloud.audit.log.v1.written")
			auditlog := &AuditLogEntry{
				ProtoPayload: tt.payload,
			}
			eventdata, err := json.Marshal(auditlog)
			if err != nil {
				t.Fatalf("failed to marshal json for test %s: %s", tt.name, err)
			}
			e.SetDataContentType("application/json")
			e.SetData("application/json", eventdata)

			if err := LabelGceInstance(context.Background(), e); err != nil {
				t.Fatalf("LabelGceInstance(%s): unexpected error %s", tt.name, err)
			}
			// check that we updated creator label if expected.
			if tt.newCreator != "" {
				newvalue, ok := fake.labels["creator"]
				if !ok || newvalue != tt.newCreator {
					t.Fatalf("LabelGceInstance(%s): incorrect creator label: got %s, want %s", tt.name, newvalue, tt.newCreator)
				}
			}

		})
	}
}
