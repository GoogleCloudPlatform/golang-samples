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

// Package topics is a tool to manage Google Cloud Pub/Sub topics by using the Pub/Sub API.
// See more about Google Cloud Pub/Sub at https://cloud.google.com/pubsub/docs/overview.
package topics

// [START pubsub_delete_topic]

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

func delete(client *pubsub.Client, topic string) error {
	ctx := context.Background()
	t := client.Topic(topic)
	if err := t.Delete(ctx); err != nil {
		return fmt.Errorf("Delete: %v", err)
	}
	fmt.Printf("Deleted topic: %v\n", t)
	return nil
}

// [END pubsub_delete_topic]
