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

// [START pubsub_old_version_create_subscription_with_smt]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

func createSubscriptionWithSMT(w io.Writer, projectID, subID string, topic *pubsub.Topic) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	// topic of type https://godoc.org/cloud.google.com/go/pubsub#Topic
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	code := `function redactSSN(message, metadata) {
			const data = JSON.parse(message.data);
			delete data['ssn'];
			message.data = JSON.stringify(data);
			return message;"
		}`
	transform := pubsub.MessageTransform{
		Transform: &pubsub.JavaScriptUDF{
			FunctionName: "redactSSN",
			Code:         code,
		},
	}
	cfg := pubsub.SubscriptionConfig{
		Topic:             topic,
		MessageTransforms: []pubsub.MessageTransform{transform},
	}
	sub, err := client.CreateSubscription(ctx, subID, cfg)
	if err != nil {
		return fmt.Errorf("CreateSubscription: %w", err)
	}
	fmt.Fprintf(w, "Created subscription with message transform: %v\n", sub)
	return nil
}

// [END pubsub_old_version_create_subscription_with_smt]
