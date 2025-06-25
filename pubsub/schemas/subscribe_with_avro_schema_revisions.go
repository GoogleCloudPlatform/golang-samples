// Copyright 2023 Google LLC
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

// [START pubsub_subscribe_avro_records_with_revisions]
import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"cloud.google.com/go/pubsub/v2"
	schema "cloud.google.com/go/pubsub/v2/apiv1"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"github.com/linkedin/goavro/v2"
)

func subscribeWithAvroSchemaRevisions(w io.Writer, projectID, subID string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}

	schemaClient, err := schema.NewSchemaClient(ctx)
	if err != nil {
		return fmt.Errorf("pubsub.NewSchemaClient: %w", err)
	}

	// Create the cache for the codecs for different revision IDs.
	revisionCodecs := make(map[string]*goavro.Codec)

	sub := client.Subscriber(subID)
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var mu sync.Mutex
	sub.Receive(ctx2, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		name := msg.Attributes["googclient_schemaname"]
		revision := msg.Attributes["googclient_schemarevisionid"]

		codec, ok := revisionCodecs[revision]
		// If the codec doesn't exist in the map, this is the first time we
		// are seeing this revision. We need to fetch the schema and cache the
		// codec. It would be more typical to do this asynchronously, but is
		// shown here in a synchronous way to ease readability.
		if !ok {
			s := &pubsubpb.GetSchemaRequest{
				Name: fmt.Sprintf("%s@%s", name, revision),
				View: pubsubpb.SchemaView_FULL,
			}
			schema, err := schemaClient.GetSchema(ctx, s)
			if err != nil {
				fmt.Fprintf(w, "Nacking, cannot read message without schema: %v\n", err)
				msg.Nack()
				return
			}
			codec, err = goavro.NewCodec(schema.Definition)
			if err != nil {
				msg.Nack()
				fmt.Fprintf(w, "goavro.NewCodec err: %v\n", err)
			}
			revisionCodecs[revision] = codec
		}

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

// [END pubsub_subscribe_avro_records_with_revisions]
