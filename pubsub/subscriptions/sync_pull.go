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

// [START pubsub_subscriber_sync_pull]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
	pubsubV1 "cloud.google.com/go/pubsub/apiv1"
	pb "google.golang.org/genproto/googleapis/pubsub/v1"
)

func pullMsgsSync(w io.Writer, projectID, subName string, topic *pubsub.Topic, maxMessages int32) ([]string, error) {
	// projectID := "my-project-id"
	// subName := projectID + "-example-sub"
	// topic of type https://godoc.org/cloud.google.com/go/pubsub#Topic
	ctx := context.Background()

	// Publish 10 messages on the topic.
	var results []*pubsub.PublishResult
	for i := 0; i < 15; i++ {
		res := topic.Publish(ctx, &pubsub.Message{
			Data: []byte(fmt.Sprintf("hello world #%d", i)),
		})
		results = append(results, res)
	}

	// Check that all messages were published.
	for _, r := range results {
		_, err := r.Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("Get: %v", err)
		}
	}

	// Instantiate a client
	subClient, err := pubsubV1.NewSubscriberClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Client instantiation error: %v", err)
	}
	sub := fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subName)

	req := &pb.PullRequest{
		Subscription:      sub,
		ReturnImmediately: false,
		MaxMessages:       maxMessages,
	}

	resp, err := subClient.Pull(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Pull error: %v", err)
	}
	var msgs, ackIDs []string
	for _, msg := range resp.GetReceivedMessages() {
		ackIDs = append(ackIDs, msg.GetAckId())
		message := string(msg.GetMessage().Data)
		fmt.Printf("Got message %q\n", message)
		msgs = append(msgs, message)
	}

	subClient.Acknowledge(ctx, &pb.AcknowledgeRequest{
		Subscription: sub,
		AckIds:       ackIDs,
	})

	return msgs, nil
}

// [END pubsub_subscriber_sync_pull]
