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

// [START pubsublite_quickstart_subscriber]

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsublite/pscompat"
)

func main() {
	// NOTE: Set these flags for an existing Pub/Sub Lite subscription containing
	// published messages when running this sample.
	projectID := flag.String("project_id", "", "Cloud Project ID")
	zone := flag.String("zone", "", "Cloud Zone where the topic resides, e.g. us-central1-a")
	subscriptionID := flag.String("subscription_id", "", "Existing Pub/Sub Lite subscription")
	timeout := flag.Duration("timeout", 90*time.Second, "The duration to receive messages")
	flag.Parse()

	ctx := context.Background()
	subscriptionPath := fmt.Sprintf("projects/%s/locations/%s/subscriptions/%s", *projectID, *zone, *subscriptionID)

	// Configure flow control settings. These settings apply per partition.
	// The message stream is paused based on the maximum size or number of
	// messages that the subscriber has already received, whichever condition is
	// met first.
	settings := pscompat.ReceiveSettings{
		// 10 MiB. Must be greater than the allowed size of the largest message
		// (1 MiB).
		MaxOutstandingBytes: 10 * 1024 * 1024,
		// 1,000 outstanding messages. Must be > 0.
		MaxOutstandingMessages: 1000,
	}

	// Create the subscriber client.
	subscriber, err := pscompat.NewSubscriberClientWithSettings(ctx, subscriptionPath, settings)
	if err != nil {
		log.Fatalf("pscompat.NewSubscriberClientWithSettings error: %v", err)
	}

	// Listen for messages until the timeout expires.
	log.Printf("Listening to messages on %s for %v...\n", subscriptionPath, *timeout)
	cctx, cancel := context.WithTimeout(ctx, *timeout)
	defer cancel()
	var receiveCount int32

	// Receive blocks until the context is cancelled or an error occurs.
	if err := subscriber.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		// NOTE: May be called concurrently; synchronize access to shared memory.
		atomic.AddInt32(&receiveCount, 1)

		// Metadata decoded from the message ID contains the partition and offset.
		metadata, err := pscompat.ParseMessageMetadata(msg.ID)
		if err != nil {
			log.Fatalf("Failed to parse %q: %v", msg.ID, err)
		}

		fmt.Printf("Received (partition=%d, offset=%d): %s\n", metadata.Partition, metadata.Offset, string(msg.Data))
		msg.Ack()
	}); err != nil {
		log.Fatalf("SubscriberClient.Receive error: %v", err)
	}

	fmt.Printf("Received %d messages\n", receiveCount)
}

// [END pubsublite_quickstart_subscriber]
