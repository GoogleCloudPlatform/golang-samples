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

// [START pubsub_list_schemas]
import (
	"context"
	"fmt"
	"io"

	pubsub "cloud.google.com/go/pubsub/v2/apiv1"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/api/iterator"
)

func listSchemas(w io.Writer, projectID string) ([]*pubsubpb.Schema, error) {
	// projectID := "my-project-id"
	ctx := context.Background()
	client, err := pubsub.NewSchemaClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewSchemaClient: %w", err)
	}
	defer client.Close()

	var schemas []*pubsubpb.Schema

	req := &pubsubpb.ListSchemasRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		View:   pubsubpb.SchemaView_FULL,
	}
	schemaIter := client.ListSchemas(ctx, req)
	for {
		sc, err := schemaIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("schemaIter.Next: %w", err)
		}
		fmt.Fprintf(w, "Got schema: %#v\n", sc)
		schemas = append(schemas, sc)
	}

	fmt.Fprintf(w, "Got %d schemas", len(schemas))
	return schemas, nil
}

// [END pubsub_list_schemas]
