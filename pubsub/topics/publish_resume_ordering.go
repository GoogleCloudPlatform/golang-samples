// Copyright 2019 Google LLC
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

// [START pubsub_resume_publish_with_ordering_keys]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub/v2"
	"google.golang.org/api/option"
)

func resumePublishWithOrderingKey(w io.Writer, projectID, topicID string) {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()

	// Pub/Sub's ordered delivery guarantee only applies when publishes for an ordering key are in the same region
	// For list of locational endpoints for Pub/Sub, see https://cloud.google.com/pubsub/docs/reference/service_apis_overview#list_of_locational_endpoints
	client, err := pubsub.NewClient(ctx, projectID,
		option.WithEndpoint("us-east1-pubsub.googleapis.com:443"))
	if err != nil {
		fmt.Fprintf(w, "pubsub.NewClient: %v", err)
		return
	}
	defer client.Close()

	// client.Publisher can be passed a topic ID (e.g. "my-topic") or
	// a fully qualified name (e.g. "projects/my-project/topics/my-topic").
	// If a topic ID is provided, the project ID from the client is used.
	// Reuse this publisher for all publish calls to send messages in batches.
	publisher := client.Publisher(topicID)
	publisher.EnableMessageOrdering = true
	key := "some-ordering-key"

	result := publisher.Publish(ctx, &pubsub.Message{
		Data:        []byte("some-message"),
		OrderingKey: key,
	})
	_, err = result.Get(ctx)
	if err != nil {
		// Error handling code can be added here.
		fmt.Printf("Failed to publish: %s\n", err)

		// Resume publish on an ordering key that has had unrecoverable errors.
		// After such an error publishes with this ordering key will fail
		// until this method is called.
		publisher.ResumePublish(key)
	}

	fmt.Fprint(w, "Published a message with ordering key successfully\n")
}

// [END pubsub_resume_publish_with_ordering_keys]
