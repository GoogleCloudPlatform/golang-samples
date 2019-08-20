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

// [START pubsub_get_subscription_policy]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
)

func policy(w io.Writer, projectID, subName string) (*iam.Policy, error) {
	// projectID := "my-project-id"
	// subName := projectID + "-example-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	policy, err := client.Subscription(subName).IAM().Policy(ctx)
	if err != nil {
		return nil, fmt.Errorf("Subscription: %v", err)
	}
	for _, role := range policy.Roles() {
		fmt.Fprintf(w, "%q: %q\n", role, policy.Members(role))
	}
	return policy, nil
}

// [END pubsub_get_subscription_policy]
