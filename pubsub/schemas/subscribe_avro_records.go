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

// [START pubsub_subscribe_avro_records]
import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/linkedin/goavro/v2"
)

func subscribeWithAvroSchema(w io.Writer, projectID, subID, avscFile string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	// avscFile = "path/to/an/avro/schema/file(.avsc)/formatted/in/json"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}

	avroSchema, err := os.ReadFile(avscFile)
	if err != nil {
		return fmt.Errorf("os.ReadFile err: %w", err)
	}
	codec, err := goavro.NewCodec(string(avroSchema))
	if err != nil {
		return fmt.Errorf("goavro.NewCodec err: %w", err)
	}

	sub := client.Subscription(subID)
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var mu sync.Mutex
	sub.Receive(ctx2, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		encoding := msg.Attributes["googclient_schemaencoding"]

		var state map[string]interface{}
		if encoding == "BINARY" {
			data, _, err := codec.NativeFromBinary(msg.Data)
			if err != nil {
				fmt.Fprintf(w, "codec.NativeFromBinary err: %v\n", err)
				msg.Nack()
				return
			}
			fmt.Fprintf(w, "Received a binary-encoded message:\n%#v\n", data)
			state = data.(map[string]interface{})
		} else if encoding == "JSON" {
			data, _, err := codec.NativeFromTextual(msg.Data)
			if err != nil {
				fmt.Fprintf(w, "codec.NativeFromTextual err: %v\n", err)
				msg.Nack()
				return
			}
			fmt.Fprintf(w, "Received a JSON-encoded message:\n%#v\n", data)
			state = data.(map[string]interface{})
		} else {
			fmt.Fprintf(w, "Unknown message type(%s), nacking\n", encoding)
			msg.Nack()
			return
		}
		fmt.Fprintf(w, "%s is abbreviated as %s\n", state["name"], state["post_abbr"])
		msg.Ack()
	})
	return nil
}

// [END pubsub_subscribe_avro_records]
