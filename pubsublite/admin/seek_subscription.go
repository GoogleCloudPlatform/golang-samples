// Copyright 2021 Google LLC
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

package admin

// [START pubsublite_seek_subscription]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsublite"
)

// seekSubscription initiates a seek operation for a subscription.
func seekSubscription(w io.Writer, projectID, region, zone, subID string, seekTarget pubsublite.SeekTarget, waitForOperation bool) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// zone := "us-central1-a"
	// subID := "my-subscription"
	// seekTarget := pubsublite.Beginning
	// waitForOperation := false

	// Possible values for seekTarget:
	// - pubsublite.Beginning: replays from the beginning of all retained
	//   messages.
	// - pubsublite.End: skips past all current published messages.
	// - pubsublite.PublishTime(<time>): delivers messages with publish time
	//   greater than or equal to the specified timestamp.
	// - pubsublite.EventTime(<time>): seeks to the first message with event
	//   time greater than or equal to the specified timestamp.

	// Waiting for the seek operation to complete is optional. It indicates when
	// subscribers for all partitions are receiving messages from the seek
	// target. If subscribers are offline, the operation will complete once they
	// are online.

	ctx := context.Background()
	client, err := pubsublite.NewAdminClient(ctx, region)
	if err != nil {
		return fmt.Errorf("pubsublite.NewAdminClient: %w", err)
	}
	defer client.Close()

	// Initiate an out-of-band seek for a subscription to the specified target.
	// If an operation is returned, the seek has been successfully registered
	// and will eventually propagate to subscribers.
	subPath := fmt.Sprintf("projects/%s/locations/%s/subscriptions/%s", projectID, zone, subID)
	seekOp, err := client.SeekSubscription(ctx, subPath, seekTarget)
	if err != nil {
		return fmt.Errorf("client.SeekSubscription got err: %w", err)
	}
	fmt.Fprintf(w, "Seek operation initiated: %s\n", seekOp.Name())

	if waitForOperation {
		_, err = seekOp.Wait(ctx)
		if err != nil {
			return fmt.Errorf("seekOp.Wait got err: %w", err)
		}
		metadata, err := seekOp.Metadata()
		if err != nil {
			return fmt.Errorf("seekOp.Metadata got err: %w", err)
		}
		fmt.Fprintf(w, "Seek operation completed with metadata: %v\n", metadata)
	}
	return nil
}

// [END pubsublite_seek_subscription]
