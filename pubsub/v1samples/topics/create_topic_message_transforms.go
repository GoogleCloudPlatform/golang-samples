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

// [START pubsub_old_version_create_topic_with_smt]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

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
			return message;"
		}`
	transform := pubsub.MessageTransform{
		Transform: &pubsub.JavaScriptUDF{
			FunctionName: "redactSSN",
			Code:         code,
		},
	}
	cfg := &pubsub.TopicConfig{
		MessageTransforms: []pubsub.MessageTransform{transform},
	}
	t, err := client.CreateTopicWithConfig(ctx, topicID, cfg)
	if err != nil {
		return fmt.Errorf("CreateTopic: %w", err)
	}
	fmt.Fprintf(w, "Created topic with message transform: %v\n", t)
	return nil
}

// [END pubsub_old_version_create_topic_with_smt]
