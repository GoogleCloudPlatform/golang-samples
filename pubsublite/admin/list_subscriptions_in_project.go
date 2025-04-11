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

// [START pubsublite_list_subscriptions_in_project]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsublite"
	"google.golang.org/api/iterator"
)

func listSubscriptionsInProject(w io.Writer, projectID, region, location string) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// NOTE: location can be either a region ("us-central1") or a zone ("us-central1-a")
	// For a list of valid locations, see https://cloud.google.com/pubsub/lite/docs/locations.
	// location := "us-central1"
	ctx := context.Background()
	client, err := pubsublite.NewAdminClient(ctx, region)
	if err != nil {
		return fmt.Errorf("pubsublite.NewAdminClient: %w", err)
	}
	defer client.Close()

	// To list topics in a region, set location to a cloud region instead.
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, location)
	subIter := client.Subscriptions(ctx, parent)
	for {
		sub, err := subIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("subIter.Next() got err: %w", err)
		}
		fmt.Fprintf(w, "Got subscription config: %v", sub)
	}
	return nil
}

// [END pubsublite_list_subscriptions_in_project]
