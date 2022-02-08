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

func updateTopic(w io.Writer, projectID, region, zone, topicID, reservation string, regional bool) error {
	// projectID := "my-project-id"
	// region := "us-central1"
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

	var topicPath string
	if regional {
		topicPath = fmt.Sprintf("projects/%s/locations/%s/topics/%s", projectID, region, topicID)
	} else {
		topicPath = fmt.Sprintf("projects/%s/locations/%s/topics/%s", projectID, zone, topicID)
	}
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
		return fmt.Errorf("client.UpdateTopic got err: %v", err)
	}
	if regional {
		fmt.Fprintf(w, "Updated regional topic: %v\n", updatedCfg)
	} else {
		fmt.Fprintf(w, "Updated zonal topic: %v\n", updatedCfg)
	}
	return nil
}

// [END pubsublite_update_topic]
