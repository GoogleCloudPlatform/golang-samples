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

// Package subscription is a tool to manage Google Cloud Pub/Sub subscriptions by using the Pub/Sub API.
// See more about Google Cloud Pub/Sub at https://cloud.google.com/pubsub/docs/overview.
package subscription

// [START pubsub_subscriber_error_listener]

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

func pullMsgsError(client *pubsub.Client, subName string) error {
	ctx := context.Background()
	// If the service returns a non-retryable error, Receive returns that error after
	// all of the outstanding calls to the handler have returned.
	err := client.Subscription(subName).Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Got message: %q\n", string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		return err
	}
	return nil
}

// [END pubsub_subscriber_error_listener]
