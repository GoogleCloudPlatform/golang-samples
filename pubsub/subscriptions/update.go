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

// [START pubsub_update_push_configuration]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func updateEndpoint(w io.Writer, projectID, subscriptionName, endpoint string) error {
	// projectID := "my-project-id"
	// subscriptionName := "projects/my-project/subscriptions/my-sub"
	// endpoint := "https://my-test-project.appspot.com/push"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	req := &pubsubpb.UpdateSubscriptionRequest{
		Subscription: &pubsubpb.Subscription{
			Name: subscriptionName,
			PushConfig: &pubsubpb.PushConfig{
				PushEndpoint: endpoint,
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"push_config"},
		},
	}
	subConfig, err := client.SubscriptionAdminClient.UpdateSubscription(ctx, req)
	if err != nil {
		return fmt.Errorf("Update: %w", err)
	}
	fmt.Fprintf(w, "Updated subscription config: %v\n", subConfig)
	return nil
}

// [END pubsub_update_push_configuration]
