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

// [START pubsub_test_topic_permissions]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/pubsub/v2"
)

func testPermissions(w io.Writer, projectID, topicID string) ([]string, error) {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %w", err)
	}

	req := &iampb.TestIamPermissionsRequest{
		Resource: fmt.Sprintf("projects/%s/topics/%s", projectID, topicID),
		Permissions: []string{
			"pubsub.topics.publish",
			"pubsub.topics.update",
		},
	}
	resp, err := client.TopicAdminClient.TestIamPermissions(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error calling TestIamPermissions: %w", err)
	}
	for _, perm := range resp.Permissions {
		fmt.Fprintf(w, "Allowed: %v\n", perm)
	}
	return resp.Permissions, nil
}

// [END pubsub_test_topic_permissions]
