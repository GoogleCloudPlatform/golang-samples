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

// [START pubsub_old_version_create_topic_with_schema]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

func createTopicWithSchema(w io.Writer, projectID, topicID, schemaID string, encoding pubsub.SchemaEncoding) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	// schemaID := "my-schema-id"
	// encoding := pubsub.EncodingJSON
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}

	tc := &pubsub.TopicConfig{
		SchemaSettings: &pubsub.SchemaSettings{
			Schema:   fmt.Sprintf("projects/%s/schemas/%s", projectID, schemaID),
			Encoding: encoding,
		},
	}
	t, err := client.CreateTopicWithConfig(ctx, topicID, tc)
	if err != nil {
		return fmt.Errorf("CreateTopicWithConfig: %w", err)
	}
	fmt.Fprintf(w, "Topic with schema created: %#v\n", t)
	return nil
}

// [END pubsub_old_version_create_topic_with_schema]
