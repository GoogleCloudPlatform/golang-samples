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

// [START pubsub_create_push_subscription]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsub"
)

func createWithEndpoint(w io.Writer, projectID, subName string, topic *pubsub.Topic, endpoint string) error {
	// projectID := "my-project-id"
	// subName := projectID + "-example-sub"
	// topic of type https://godoc.org/cloud.google.com/go/pubsub#Topic
	// endpoint := "https://my-test-project.appspot.com/push"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	sub, err := client.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 10 * time.Second,
		PushConfig:  pubsub.PushConfig{Endpoint: endpoint},
	})
	if err != nil {
		return fmt.Errorf("CreateSubscription: %v", err)
	}
	fmt.Fprintf(w, "Created subscription: %v\n", sub)
	return nil
}

// [END pubsub_create_push_subscription]
