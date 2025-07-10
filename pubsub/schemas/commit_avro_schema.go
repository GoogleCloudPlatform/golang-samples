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

// [START pubsub_commit_avro_schema]
import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/pubsub"
)

// commitAvroSchema commits a new Avro schema revision to an existing schema.
func commitAvroSchema(w io.Writer, projectID, schemaID, avscFile string) error {
	// projectID := "my-project-id"
	// schemaID := "my-schema-id"
	// avscFile = "path/to/an/avro/schema/file(.avsc)/formatted/in/json"
	ctx := context.Background()
	client, err := pubsub.NewSchemaClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewSchemaClient: %w", err)
	}
	defer client.Close()

	// Read an Avro schema file formatted in JSON as a byte slice.
	avscSource, err := os.ReadFile(avscFile)
	if err != nil {
		return fmt.Errorf("error reading from file: %s", avscFile)
	}

	config := pubsub.SchemaConfig{
		Name:       fmt.Sprintf("projects/%s/schemas/%s", projectID, schemaID),
		Type:       pubsub.SchemaAvro,
		Definition: string(avscSource),
	}
	s, err := client.CommitSchema(ctx, schemaID, config)
	if err != nil {
		return fmt.Errorf("CommitSchema: %w", err)
	}
	fmt.Fprintf(w, "Committed a schema using an Avro schema: %#v\n", s)
	return nil
}

// [END pubsub_commit_avro_schema]
