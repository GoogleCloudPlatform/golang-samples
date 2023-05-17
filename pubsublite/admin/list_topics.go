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

// [START pubsublite_list_topics]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsublite"
	"google.golang.org/api/iterator"
)

func listTopics(w io.Writer, projectID, region, zone string) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// zone := "us-central1-a"
	ctx := context.Background()
	client, err := pubsublite.NewAdminClient(ctx, region)
	if err != nil {
		return fmt.Errorf("pubsublite.NewAdminClient: %w", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, zone)
	topicIter := client.Topics(ctx, parent)
	for {
		topic, err := topicIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("topicIter.Next got err: %w", err)
		}
		fmt.Fprintf(w, "Got topic: %#v\n", topic)
	}
	return nil
}

// [END pubsublite_list_topics]
