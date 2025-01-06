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

// [START pubsub_old_version_get_topic_policy]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
)

func policy(w io.Writer, projectID, topicID string) (*iam.Policy, error) {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	policy, err := client.Topic(topicID).IAM().Policy(ctx)
	if err != nil {
		return nil, fmt.Errorf("Policy: %w", err)
	}
	for _, role := range policy.Roles() {
		fmt.Fprint(w, policy.Members(role))
	}
	return policy, nil
}

// [END pubsub_old_version_get_topic_policy]
