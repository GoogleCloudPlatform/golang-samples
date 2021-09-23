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

// [START pubsub_create_topic_with_schema]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

func createWithSchema(w io.Writer, projectID, topicID, schemaID string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	tc := &pubsub.TopicConfig{
		SchemaSettings: &pubsub.SchemaSettings{
			Schema:   fmt.Sprintf("projects/%s/schemas/%s", projectID, schemaID),
			Encoding: pubsub.EncodingJSON,
		},
	}
	t, err := client.CreateTopicWithConfig(ctx, topicID, tc)
	if err != nil {
		return fmt.Errorf("client.CreateTopicWithConfig: %v", err)
	}
	fmt.Fprintf(w, "Topic created with schema: %v\n", t)
	return nil
}

// [END pubsub_create_topic_with_schema]
