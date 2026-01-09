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

package schema

// [START pubsub_subscribe_proto_messages]
import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"cloud.google.com/go/pubsub/v2"
	statepb "github.com/GoogleCloudPlatform/golang-samples/internal/pubsub/schemas"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func subscribeWithProtoSchema(w io.Writer, projectID, subID string) error {
	// projectID := "my-project-id"
	// subID := "my-sub"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}

	// Create an instance of the message to be decoded (a single U.S. state).
	state := &statepb.State{}

	sub := client.Subscriber(subID)
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var mu sync.Mutex
	sub.Receive(ctx2, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		encoding := msg.Attributes["googclient_schemaencoding"]

		if encoding == "BINARY" {
			if err := proto.Unmarshal(msg.Data, state); err != nil {
				fmt.Fprintf(w, "proto.Unmarshal err: %v\n", err)
				msg.Nack()
				return
			}
			fmt.Printf("Received a binary-encoded message:\n%#v\n", state)
		} else if encoding == "JSON" {
			if err := protojson.Unmarshal(msg.Data, state); err != nil {
				fmt.Fprintf(w, "proto.Unmarshal err: %v\n", err)
				msg.Nack()
				return
			}
			fmt.Fprintf(w, "Received a JSON-encoded message:\n%#v\n", state)
		} else {
			fmt.Fprintf(w, "Unknown message type(%s), nacking\n", encoding)
			msg.Nack()
			return
		}
		fmt.Fprintf(w, "%s is abbreviated as %s\n", state.Name, state.PostAbbr)
		msg.Ack()
	})
	return nil
}

// [END pubsub_subscribe_proto_messages]
