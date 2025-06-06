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

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
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
	req := &pubsubpb.GetTopicRequest{
		Topic: fmt.Sprintf("projects/%s/topics/%s", projectID, topicID),
	}
	t, err := client.TopicAdminClient.GetTopic(ctx, req)
	if err != nil {
		return fmt.Errorf("got err in GetTopic: %w", err)
	}
	encoding := t.SchemaSettings.Encoding

	var msg []byte
	switch encoding {
	case pubsubpb.Encoding_BINARY:
		msg, err = proto.Marshal(state)
		if err != nil {
			return fmt.Errorf("proto.Marshal err: %w", err)
		}
	case pubsubpb.Encoding_JSON:
		msg, err = protojson.Marshal(state)
		if err != nil {
			return fmt.Errorf("protojson.Marshal err: %w", err)
		}
	default:
		return fmt.Errorf("invalid encoding: %v", encoding)
	}

	publisher := client.Publisher(t.GetName())
	result := publisher.Publish(ctx, &pubsub.Message{
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
