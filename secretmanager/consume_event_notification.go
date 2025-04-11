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

package secretmanager

// [START secretmanager_consume_event_notification]
import (
	"context"
	"fmt"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Attributes PubSubAttributes `json:"attributes"`
	Data       []byte           `json:"data"`
}

// PubSubAttributes are attributes from the Pub/Sub event.
type PubSubAttributes struct {
	SecretId  string `json:"secretId"`
	EventType string `json:"eventType"`
}

// ConsumeEventNotification demonstrates how to consume and process the Pub/Sub
// notification from Secret Manager.
func ConsumeEventNotification(ctx context.Context, m PubSubMessage) (string, error) {
	// The printed value will be visible in Cloud Logging:
	//
	//     https://cloud.google.com/functions/docs/monitoring/logging
	//
	eventType := m.Attributes.EventType
	secretID := m.Attributes.SecretId
	data := m.Data

	return fmt.Sprintf("Received %s for %s. New metadata: %q.",
		eventType, secretID, data), nil
}

// [END secretmanager_consume_event_notification]
