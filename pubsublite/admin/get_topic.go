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

// [START pubsublite_get_topic]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsublite"
)

func getTopic(w io.Writer, projectID, region, zone, topicID string, regional bool) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// zone := "us-central1-a"
	// topicID := "my-topic"
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
	topic, err := client.Topic(ctx, topicPath)
	if err != nil {
		return fmt.Errorf("client.Topic got err: %v", err)
	}
	fmt.Fprintf(w, "Got topic: %#v\n", *topic)
	return nil
}

// [END pubsublite_get_topic]
