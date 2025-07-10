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

// [START pubsub_publish_proto_messages]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
	statepb "github.com/GoogleCloudPlatform/golang-samples/internal/pubsub/schemas"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func publishProtoMessages(w io.Writer, projectID, topicID string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}

	state := &statepb.State{
		Name:     "Alaska",
		PostAbbr: "AK",
	}

	// Get the topic encoding type.
	t := client.Topic(topicID)
	cfg, err := t.Config(ctx)
	if err != nil {
		return fmt.Errorf("topic.Config err: %w", err)
	}
	encoding := cfg.SchemaSettings.Encoding

	var msg []byte
	switch encoding {
	case pubsub.EncodingBinary:
		msg, err = proto.Marshal(state)
		if err != nil {
			return fmt.Errorf("proto.Marshal err: %w", err)
		}
	case pubsub.EncodingJSON:
		msg, err = protojson.Marshal(state)
		if err != nil {
			return fmt.Errorf("protojson.Marshal err: %w", err)
		}
	default:
		return fmt.Errorf("invalid encoding: %v", encoding)
	}

	result := t.Publish(ctx, &pubsub.Message{
		Data: msg,
	})
	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("result.Get: %w", err)
	}
	fmt.Fprintf(w, "Published proto message with %#v encoding: %s\n", encoding, string(msg))
	return nil
}

// [END pubsub_publish_proto_messages]
