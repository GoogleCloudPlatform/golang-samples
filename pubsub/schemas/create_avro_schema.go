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

// [START pubsub_create_avro_schema]
import (
	"context"
	"fmt"
	"io"
	"os"

	pubsub "cloud.google.com/go/pubsub/v2/apiv1"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)

// createAvroSchema creates a schema resource from a JSON-formatted Avro schema file.
func createAvroSchema(w io.Writer, projectID, schemaID, avscFile string) error {
	// projectID := "my-project-id"
	// schemaID := "my-schema"
	// avscFile = "path/to/an/avro/schema/file(.avsc)/formatted/in/json"
	ctx := context.Background()
	client, err := pubsub.NewSchemaClient(ctx)
	if err != nil {
		return fmt.Errorf("pubsub.NewSchemaClient: %w", err)
	}
	defer client.Close()

	avscSource, err := os.ReadFile(avscFile)
	if err != nil {
		return fmt.Errorf("error reading from file: %s", avscFile)
	}

	req := &pubsubpb.CreateSchemaRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		Schema: &pubsubpb.Schema{
			Type:       pubsubpb.Schema_AVRO,
			Definition: string(avscSource),
		},
		SchemaId: schemaID,
	}
	s, err := client.CreateSchema(ctx, req)
	if err != nil {
		return fmt.Errorf("error calling CreateSchema: %w", err)
	}
	fmt.Fprintf(w, "Schema created: %#v\n", s)
	return nil
}

// [END pubsub_create_avro_schema]
