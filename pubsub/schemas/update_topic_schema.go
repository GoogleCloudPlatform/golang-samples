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

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func updateTopicSchema(w io.Writer, projectID, topicID, firstRevisionID, lastRevisionID string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic" // an existing topic that has schema settings attached to it.
	// firstRevisionID := "my-revision-id"
	// lastRevisionID := "my-revision-id"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}

	// This updates the first / last revision ID for the topic's schema.
	// To clear the schema entirely, use a zero valued (empty) SchemaSettings
	// with the same field mask.
	req := &pubsubpb.UpdateTopicRequest{
		Topic: &pubsubpb.Topic{
			Name: fmt.Sprintf("projects/%s/topics/%s", projectID, topicID),
			SchemaSettings: &pubsubpb.SchemaSettings{
				FirstRevisionId: firstRevisionID,
				LastRevisionId:  lastRevisionID,
			},
		},
		// Construct a field mask to indicate which field to update in the topic.
		// Fields are specified relative to the topic
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"schema_settings.first_revision_id", "schema_settings.last_revision_id"},
		},
	}
	gotTopicCfg, err := client.TopicAdminClient.UpdateTopic(ctx, req)
	if err != nil {
		fmt.Fprintf(w, "topic.Update err: %v\n", gotTopicCfg)
		return err
	}
	fmt.Fprintf(w, "Updated topic with schema: %#v\n", gotTopicCfg)
	return nil
}

// [END pubsub_update_topic_schema]
