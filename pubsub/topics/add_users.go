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

// [START pubsub_set_topic_policy]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/pubsub/v2"
)

func addUsers(w io.Writer, projectID, topicID string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	topicName := fmt.Sprintf("projects/%s/topics/%s", projectID, topicID)
	req := &iampb.GetIamPolicyRequest{
		Resource: topicName,
	}
	policy, err := client.TopicAdminClient.GetIamPolicy(ctx, req)
	if err != nil {
		return fmt.Errorf("error calling GetIamPolicy: %w", err)
	}
	b1 := &iampb.Binding{
		Role:    "roles/viewer",
		Members: []string{"allUsers"},
	}
	b2 := &iampb.Binding{
		Role: "roles/editor",
		// Other valid prefixes are "serviceAccount:", "user:"
		// See the documentation for more values.
		Members: []string{"group:cloud-logs@google.com"},
	}
	policy.Bindings = append(policy.Bindings, b1, b2)

	setRequest := &iampb.SetIamPolicyRequest{
		Resource: topicName,
		Policy:   policy,
	}
	_, err = client.TopicAdminClient.SetIamPolicy(ctx, setRequest)
	if err != nil {
		return fmt.Errorf("error calling SetIamPolicy: %w", err)
	}
	return nil
}

// [END pubsub_set_topic_policy]
