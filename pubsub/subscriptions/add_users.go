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

package subscriptions

// [START pubsub_set_subscription_policy]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/pubsub/v2"
)

func addUsersToSubscription(w io.Writer, projectID, subID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	subName := fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subID)
	req := &iampb.GetIamPolicyRequest{
		Resource: subName,
	}
	policy, err := client.SubscriptionAdminClient.GetIamPolicy(ctx, req)
	if err != nil {
		return fmt.Errorf("error calling GetIamPolicy: %w", err)
	}
	b := &iampb.Binding{
		Role: "roles/editor",
		// Other valid prefixes are "serviceAccount:", "user:"
		// See the documentation for more values.
		Members: []string{"group:cloud-logs@google.com"},
	}
	policy.Bindings = append(policy.Bindings, b)

	setRequest := &iampb.SetIamPolicyRequest{
		Resource: subName,
		Policy:   policy,
	}
	_, err = client.SubscriptionAdminClient.SetIamPolicy(ctx, setRequest)
	if err != nil {
		return fmt.Errorf("error calling SetIamPolicy: %w", err)
	}
	fmt.Fprintln(w, "Added roles to subscription.")
	return nil
}

// [END pubsub_set_subscription_policy]
