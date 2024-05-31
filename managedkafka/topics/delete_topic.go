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

// [START managedkafka_delete_topic]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/managedkafka/apiv1/managedkafkapb"
	"google.golang.org/api/option"

	managedkafka "cloud.google.com/go/managedkafka/apiv1"
)

func deleteTopic(w io.Writer, projectID, region, clusterID, topicID string, opts ...option.ClientOption) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// clusterID := "my-cluster"
	// topicID := "my-topic"
	ctx := context.Background()
	client, err := managedkafka.NewClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("managedkafka.NewClient got err: %w", err)
	}
	defer client.Close()

	clusterPath := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectID, region, clusterID)
	topicPath := fmt.Sprintf("%s/topics/%s", clusterPath, topicID)
	req := &managedkafkapb.DeleteTopicRequest{
		Name: topicPath,
	}
	if err := client.DeleteTopic(ctx, req); err != nil {
		return fmt.Errorf("client.DeleteTopic got err: %w", err)
	}
	fmt.Fprint(w, "Deleted topic\n")
	return nil
}

// [END managedkafka_delete_topic]
