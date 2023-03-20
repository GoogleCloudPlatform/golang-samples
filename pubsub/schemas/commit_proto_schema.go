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

// [START pubsub_commit_proto_schema]
import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"cloud.google.com/go/pubsub"
)

// commitProtoSchema commits a new proto schema revision to an existing schema.
func commitProtoSchema(w io.Writer, projectID, schemaID, protoFile string) error {
	// projectID := "my-project-id"
	// schemaID := "my-schema"
	// protoFile = "path/to/a/proto/schema/file(.proto)/formatted/in/protocol/buffers"
	ctx := context.Background()
	client, err := pubsub.NewSchemaClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewSchemaClient: %v", err)
	}
	defer client.Close()

	// Read a proto file as a byte slice.
	protoSource, err := ioutil.ReadFile(protoFile)
	if err != nil {
		return fmt.Errorf("error reading from file: %s", protoFile)
	}

	config := pubsub.SchemaConfig{
		Name:       fmt.Sprintf("projects/%s/schemas/%s", projectID, schemaID),
		Type:       pubsub.SchemaProtocolBuffer,
		Definition: string(protoSource),
	}
	s, err := client.CommitSchema(ctx, schemaID, config)
	if err != nil {
		return fmt.Errorf("CommitSchema: %v", err)
	}
	fmt.Fprintf(w, "Committed a schema using a protobuf schema: %#v\n", s)
	return nil
}

// [END pubsub_commit_proto_schema]
