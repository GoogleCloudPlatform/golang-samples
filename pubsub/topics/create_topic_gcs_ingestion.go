// Copyright 2024 Google LLC
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

// [START pubsub_create_topic_with_cloud_storage_ingestion]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func createTopicWithCloudStorageIngestion(w io.Writer, projectID, topicID, bucket, matchGlob, minimumObjectCreateTime, delimiter string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	// bucket := "my-bucket"
	// matchGlob := "**.txt"
	// minimumObjectCreateTime := "2006-01-02T15:04:05Z"
	// delimiter := ","

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	minCreateTime, err := time.Parse(time.RFC3339, minimumObjectCreateTime)
	if err != nil {
		return err
	}

	topicpb := &pubsubpb.Topic{
		Name: fmt.Sprintf("projects/%s/topics/%s", projectID, topicID),
		IngestionDataSourceSettings: &pubsubpb.IngestionDataSourceSettings{
			Source: &pubsubpb.IngestionDataSourceSettings_CloudStorage_{
				CloudStorage: &pubsubpb.IngestionDataSourceSettings_CloudStorage{
					Bucket: bucket,
					// Alternatively, can be Avro or PubSubAvro formats. See
					InputFormat: &pubsubpb.IngestionDataSourceSettings_CloudStorage_TextFormat_{
						TextFormat: &pubsubpb.IngestionDataSourceSettings_CloudStorage_TextFormat{
							Delimiter: &delimiter,
						},
					},
					MatchGlob:               matchGlob,
					MinimumObjectCreateTime: timestamppb.New(minCreateTime),
				},
			},
		},
	}
	t, err := client.TopicAdminClient.CreateTopic(ctx, topicpb)
	if err != nil {
		return fmt.Errorf("CreateTopic: %w", err)
	}
	fmt.Fprintf(w, "Cloud storage topic created: %v\n", t)
	return nil
}

// [END pubsub_create_topic_with_cloud_storage_ingestion]
