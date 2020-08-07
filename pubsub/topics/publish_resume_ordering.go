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

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func resumePublishWithOrderingKey(w io.Writer, projectID, topicID string) {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()

	// Sending messages to the same region ensures they are received in order
	// even when multiple publishers are used.
	client, err := pubsub.NewClient(ctx, projectID,
		option.WithEndpoint("us-east1-pubsub.googleapis.com:443"))
	if err != nil {
		fmt.Fprintf(w, "pubsub.NewClient: %v", err)
		return
	}
	defer client.Close()

	t := client.Topic(topicID)
	t.EnableMessageOrdering = true
	key := "some-ordering-key"

	res := t.Publish(ctx, &pubsub.Message{
		Data:        []byte("some-message"),
		OrderingKey: key,
	})
	_, err = res.Get(ctx)
	if err != nil {
		// Error handling code can be added here.
		fmt.Printf("Failed to publish: %s\n", err)

		// Resume publish on an ordering key that has had unrecoverable errors.
		// Until this method is not called, subsequent publishes with this ordering
		// key will fail.
		t.ResumePublish(key)
	}

	fmt.Fprint(w, "Published a message with ordering key successfully\n")
}

// [END pubsub_resume_publish_with_ordering_keys]
