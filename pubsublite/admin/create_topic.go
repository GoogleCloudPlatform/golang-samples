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

// [START pubsublite_create_topic]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsublite"
)

func createTopic(w io.Writer, projectID, region, zone, topicID, reservation string, regional bool) error {
	// projectID := "my-project-id"
	// region := "us-central1" // see https://cloud.google.com/pubsub/lite/docs/locations
	// zone := "us-central1-a"
	// topicID := "my-topic"
	// reservation := "projects/my-project-id/reservations/my-reservation"
	// regional := "true"
	ctx := context.Background()
	client, err := pubsublite.NewAdminClient(ctx, region)
	if err != nil {
		return fmt.Errorf("pubsublite.NewAdminClient: %v", err)
	}
	defer client.Close()

	const gib = 1 << 30

	var topicPath string
	if regional {
		topicPath = fmt.Sprintf("projects/%s/locations/%s/topics/%s", projectID, region, topicID)
	} else {
		topicPath = fmt.Sprintf("projects/%s/locations/%s/topics/%s", projectID, zone, topicID)
	}
	// For ranges of fields in TopicConfig, see https://pkg.go.dev/cloud.google.com/go/pubsublite/#TopicConfig
	topic, err := client.CreateTopic(ctx, pubsublite.TopicConfig{
		Name:                       topicPath,
		PartitionCount:             2, // Must be >= 1 and cannot decrease after creation.
		PublishCapacityMiBPerSec:   4,
		SubscribeCapacityMiBPerSec: 8,
		PerPartitionBytes:          30 * gib,
		RetentionDuration:          pubsublite.InfiniteRetention,
		ThroughputReservation:      reservation,
	})
	if err != nil {
		return fmt.Errorf("client.CreateTopic got err: %v", err)
	}
	if regional {
		fmt.Fprintf(w, "Created regional topic: %#v\n", topic)
	} else {
		fmt.Fprintf(w, "Created zonal topic: %#v\n", topic)
	}
	return nil
}

// [END pubsublite_create_topic]
