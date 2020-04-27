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

package topics

// [START pubsub_publisher_batch_settings]
import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
)

func publishWithSettings(w io.Writer, projectID, topicID string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	var results []*pubsub.PublishResult
	var resultErrors []error
	t := client.Topic(topicID)
	t.PublishSettings.ByteThreshold = 5000
	t.PublishSettings.CountThreshold = 10
	t.PublishSettings.DelayThreshold = 100 * time.Millisecond

	for i := 0; i < 10; i++ {
		result := t.Publish(ctx, &pubsub.Message{
			Data: []byte("Message " + strconv.Itoa(i)),
		})
		results = append(results, result)
	}
	// The Get method blocks until a server-generated ID or
	// an error is returned for the published message.
	for i, res := range results {
		id, err := res.Get(ctx)
		if err != nil {
			resultErrors = append(resultErrors, err)
			fmt.Fprintf(w, "Failed to publish: %v", err)
			continue
		}
		fmt.Fprintf(w, "Published message %d; msg ID: %v\n", i, id)
	}
	if len(resultErrors) != 0 {
		return fmt.Errorf("Get: %v", resultErrors[len(resultErrors)-1])
	}
	fmt.Fprintf(w, "Published messages with batch settings.")
	return nil
}

// [END pubsub_publisher_batch_settings]
