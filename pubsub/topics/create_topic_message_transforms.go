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

// [START pubsub_create_topic_with_smt]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)

// createTopicWithSMT creates a topic with a single message transform function applied.
func createTopicWithSMT(w io.Writer, projectID, topicID string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	code := `function redactSSN(message, metadata) {
			const data = JSON.parse(message.data);
			delete data['ssn'];
			message.data = JSON.stringify(data);
			return message;
		}`
	transform := &pubsubpb.MessageTransform{
		Transform: &pubsubpb.MessageTransform_JavascriptUdf{
			JavascriptUdf: &pubsubpb.JavaScriptUDF{
				FunctionName: "redactSSN",
				Code:         code,
			},
		},
	}

	topic := &pubsubpb.Topic{
		Name:              fmt.Sprintf("projects/%s/topics/%s", projectID, topicID),
		MessageTransforms: []*pubsubpb.MessageTransform{transform},
	}

	topic, err = client.TopicAdminClient.CreateTopic(ctx, topic)
	if err != nil {
		return fmt.Errorf("CreateTopic: %w", err)
	}

	fmt.Fprintf(w, "Created topic with message transform: %v\n", topic)
	return nil
}

// [END pubsub_create_topic_with_smt]
