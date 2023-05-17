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

// [START pubsublite_update_topic]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/pubsublite"
)

func updateTopic(w io.Writer, projectID, region, location, topicID, reservation string) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// NOTE: location can be either a region ("us-central1") or a zone ("us-central1-a")
	// For a list of valid locations, see https://cloud.google.com/pubsub/lite/docs/locations.
	// location := "us-central1"
	// topicID := "my-topic"
	// reservation := "projects/my-project-id/reservations/my-reservation"
	ctx := context.Background()
	client, err := pubsublite.NewAdminClient(ctx, region)
	if err != nil {
		return fmt.Errorf("pubsublite.NewAdminClient: %w", err)
	}
	defer client.Close()

	topicPath := fmt.Sprintf("projects/%s/locations/%s/topics/%s", projectID, location, topicID)
	// For ranges of fields in TopicConfigToUpdate, see https://pkg.go.dev/cloud.google.com/go/pubsublite/#TopicConfigToUpdate
	config := pubsublite.TopicConfigToUpdate{
		Name:                       topicPath,
		PartitionCount:             3, // Partition count cannot decrease.
		PublishCapacityMiBPerSec:   8,
		SubscribeCapacityMiBPerSec: 16,
		PerPartitionBytes:          60 * 1024 * 1024 * 1024,
		RetentionDuration:          24 * time.Hour,
		ThroughputReservation:      reservation,
	}
	updatedCfg, err := client.UpdateTopic(ctx, config)
	if err != nil {
		return fmt.Errorf("client.UpdateTopic got err: %w", err)
	}
	fmt.Fprintf(w, "Updated topic: %v\n", updatedCfg)
	return nil
}

// [END pubsublite_update_topic]
