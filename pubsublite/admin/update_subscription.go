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

// [START pubsublite_update_subscription]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsublite"
)

func updateSubscription(w io.Writer, projectID, region, location, subID string) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// NOTE: location can be either a region ("us-central1") or a zone ("us-central1-a")
	// For a list of valid locations, see https://cloud.google.com/pubsub/lite/docs/locations.
	// location := "us-central1"
	// subID := "my-subscription"
	ctx := context.Background()
	client, err := pubsublite.NewAdminClient(ctx, region)
	if err != nil {
		return fmt.Errorf("pubsublite.NewAdminClient: %w", err)
	}
	defer client.Close()

	subPath := fmt.Sprintf("projects/%s/locations/%s/subscriptions/%s", projectID, location, subID)
	config := pubsublite.SubscriptionConfigToUpdate{
		Name:                subPath,
		DeliveryRequirement: pubsublite.DeliverAfterStored,
	}
	updatedCfg, err := client.UpdateSubscription(ctx, config)
	if err != nil {
		return fmt.Errorf("client.UpdateSubscription got err: %w", err)
	}
	fmt.Fprintf(w, "Updated subscription: %#v\n", updatedCfg)
	return nil
}

// [END pubsublite_update_subscription]
