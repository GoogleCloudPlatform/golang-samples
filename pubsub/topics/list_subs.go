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

// [START pubsub_list_topic_subscriptions]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/api/iterator"
)

func listSubscriptions(w io.Writer, projectID, topicID string) error {
	// projectID := "my-project-id"
	// topicName := "projects/sample-248520/topics/ocr-go-test-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	req := &pubsubpb.ListTopicSubscriptionsRequest{
		Topic: fmt.Sprintf("projects/%s/topics/%s", projectID, topicID),
	}
	it := client.TopicAdminClient.ListTopicSubscriptions(ctx, req)
	for {
		sub, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("error listing topic subscriptions: %w", err)
		}
		fmt.Fprintf(w, "got subscription: %s\n", sub)
	}
	return nil
}

// [END pubsub_list_topic_subscriptions]
