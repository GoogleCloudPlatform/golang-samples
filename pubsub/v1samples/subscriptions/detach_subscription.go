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

package subscriptions

// [START pubsub_old_version_detach_subscription]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

func detachSubscription(w io.Writer, projectID, subName string) error {
	// projectID is the project which contains the topic you manage.
	// This might differ from the project which contains the subscription
	// you wish to detach, which can exist in any GCP project.
	// projectID := "my-project-id"
	// subName := "projects/some-project/subscriptions/my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	// Call DetachSubscription, which detaches a subscription from
	// a topic. This can only be done if you have the
	// `pubsub.topics.detachSubscription` role on the topic.
	_, err = client.DetachSubscription(ctx, subName)
	if err != nil {
		return fmt.Errorf("detach subscription failed: %w", err)
	}

	fmt.Fprintf(w, "Detached subscription %s", subName)
	return nil
}

// [END pubsub_old_version_detach_subscription]
