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

package publish

// [START pubsublite_publish_custom_attributes]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsublite/pscompat"
)

func publishWithCustomAttributes(w io.Writer, projectID, zone, topicID string) error {
	// projectID := "my-project-id"
	// zone := "us-central1-a"
	// topicID := "my-topic"
	ctx := context.Background()
	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projectID, zone, topicID)

	// Create the publisher client.
	publisher, err := pscompat.NewPublisherClient(ctx, topicPath)
	if err != nil {
		return fmt.Errorf("pscompat.NewPublisherClient error: %w", err)
	}

	// Ensure the publisher will be shut down.
	defer publisher.Stop()

	// Publish a message with custom attributes.
	result := publisher.Publish(ctx, &pubsub.Message{
		Data: []byte("message-with-custom-attributes"),
		Attributes: map[string]string{
			"year":   "2020",
			"author": "unknown",
		},
	})

	// Get blocks until the result is ready.
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("publish error: %w", err)
	}

	fmt.Fprintf(w, "Published a message with custom attributes: %v\n", id)
	return publisher.Error()
}

// [END pubsublite_publish_custom_attributes]
