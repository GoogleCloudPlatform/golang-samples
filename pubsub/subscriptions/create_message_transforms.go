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

// [START pubsub_create_subscription_with_smt]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)

// createSubscriptionWithSMT creates a subscription with a single message transform function applied.
func createSubscriptionWithSMT(w io.Writer, projectID, topicID, subID string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	// subID := "my-sub"
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
			return message;
		}`

	transform := &pubsubpb.MessageTransform{
		Transform: &pubsubpb.MessageTransform_JavascriptUdf{
			JavascriptUdf: &pubsubpb.JavaScriptUDF{
				FunctionName: "redactSSN",
				Code:         code,
			},
		},
	}

	sub := &pubsubpb.Subscription{
		Name:              fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subID),
		Topic:             fmt.Sprintf("projects/%s/topics/%s", projectID, topicID),
		MessageTransforms: []*pubsubpb.MessageTransform{transform},
	}
	sub, err = client.SubscriptionAdminClient.CreateSubscription(ctx, sub)
	if err != nil {
		return fmt.Errorf("CreateSubscription: %w", err)
	}
	fmt.Fprintf(w, "Created subscription with message transform: %v\n", sub)
	return nil
}

// [END pubsub_create_subscription_with_smt]
