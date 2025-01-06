// Copyright 2025 Google LLC
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

package topics

// [START pubsub_old_version_publisher_with_compression]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

func publishWithCompression(w io.Writer, projectID, topicID string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	// msg := "Hello World"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %w", err)
	}
	defer client.Close()

	t := client.Topic(topicID)
	// Enable compression and configure the compression threshold to 10 bytes (default to 240 B).
	// Publish requests of sizes > 10 B (excluding the request headers) will get compressed.
	t.PublishSettings.EnableCompression = true
	t.PublishSettings.CompressionBytesThreshold = 10
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte("This is a test message"),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("pubsub: result.Get: %w", err)
	}
	fmt.Fprintf(w, "Published a message; msg ID: %v\n", id)
	return nil
}

// [END pubsub_old_version_publisher_with_compression]
