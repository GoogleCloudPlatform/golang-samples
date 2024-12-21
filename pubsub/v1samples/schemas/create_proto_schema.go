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

// [START pubsub_old_version_create_proto_schema]
import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/pubsub"
)

// createProtoSchema creates a schema resource from a schema proto file.
func createProtoSchema(w io.Writer, projectID, schemaID, protoFile string) error {
	// projectID := "my-project-id"
	// schemaID := "my-schema"
	// protoFile = "path/to/a/proto/schema/file(.proto)/formatted/in/protocol/buffers"
	ctx := context.Background()
	client, err := pubsub.NewSchemaClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewSchemaClient: %w", err)
	}
	defer client.Close()

	protoSource, err := os.ReadFile(protoFile)
	if err != nil {
		return fmt.Errorf("error reading from file: %s", protoFile)
	}

	config := pubsub.SchemaConfig{
		Type:       pubsub.SchemaProtocolBuffer,
		Definition: string(protoSource),
	}
	s, err := client.CreateSchema(ctx, schemaID, config)
	if err != nil {
		return fmt.Errorf("CreateSchema: %w", err)
	}
	fmt.Fprintf(w, "Schema created: %#v\n", s)
	return nil
}

// [END pubsub_old_version_create_proto_schema]
