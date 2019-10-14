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

// [START pubsub_test_subscription_permissions]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

func testPermissions(w io.Writer, projectID, subID string) ([]string, error) {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	sub := client.Subscription(subID)
	perms, err := sub.IAM().TestPermissions(ctx, []string{
		"pubsub.subscriptions.consume",
		"pubsub.subscriptions.update",
	})
	if err != nil {
		return nil, fmt.Errorf("TestPermissions: %v", err)
	}
	for _, perm := range perms {
		fmt.Fprintf(w, "Allowed: %v\n", perm)
	}
	// [END pubsub_test_subscription_permissions]
	return perms, nil
}
