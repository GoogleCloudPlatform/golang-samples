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

// [START pubsub_old_version_publish_avro_records]
import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/linkedin/goavro/v2"
)

func publishAvroRecords(w io.Writer, projectID, topicID, avscFile string) error {
	// projectID := "my-project-id"
	// topicID := "my-topic"
	// avscFile = "path/to/an/avro/schema/file(.avsc)/formatted/in/json"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}

	avroSource, err := os.ReadFile(avscFile)
	if err != nil {
		return fmt.Errorf("os.ReadFile err: %w", err)
	}
	codec, err := goavro.NewCodec(string(avroSource))
	if err != nil {
		return fmt.Errorf("goavro.NewCodec err: %w", err)
	}
	record := map[string]interface{}{"name": "Alaska", "post_abbr": "AK"}

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
		msg, err = codec.BinaryFromNative(nil, record)
		if err != nil {
			return fmt.Errorf("codec.BinaryFromNative err: %w", err)
		}
	case pubsub.EncodingJSON:
		msg, err = codec.TextualFromNative(nil, record)
		if err != nil {
			return fmt.Errorf("codec.TextualFromNative err: %w", err)
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
	fmt.Fprintf(w, "Published avro record: %s\n", string(msg))
	return nil
}

// [END pubsub_old_version_publish_avro_records]
