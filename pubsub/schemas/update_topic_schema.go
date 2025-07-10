// Copyright 2023 Google LLC
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

package schema

// [START pubsub_update_topic_schema]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

func updateTopicSchema(w io.Writer, projectID, topicID, firstRevisionID, lastRevisionID string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	// firstRevisionID := "my-revision-id"
	// lastRevisionID := "my-revision-id"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	t := client.Topic(topicID)

	// This updates the first / last revision ID for the topic's schema.
	// To clear the schema entirely, use a zero valued (empty) SchemaSettings.
	tc := pubsub.TopicConfigToUpdate{
		SchemaSettings: &pubsub.SchemaSettings{
			FirstRevisionID: firstRevisionID,
			LastRevisionID:  lastRevisionID,
		},
	}

	gotTopicCfg, err := t.Update(ctx, tc)
	if err != nil {
		fmt.Fprintf(w, "topic.Update err: %v\n", gotTopicCfg)
		return err
	}
	fmt.Fprintf(w, "Updated topic with schema: %#v\n", gotTopicCfg)
	return nil
}

// [END pubsub_update_topic_schema]
