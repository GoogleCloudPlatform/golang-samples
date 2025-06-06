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

// [START pubsub_list_schema_revisions]
import (
	"context"
	"fmt"
	"io"

	pubsub "cloud.google.com/go/pubsub/v2/apiv1"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/api/iterator"
)

func listSchemaRevisions(w io.Writer, projectID, schemaID string) ([]*pubsubpb.Schema, error) {
	// projectID := "my-project-id"
	// schemaID := "my-schema-id"
	ctx := context.Background()
	client, err := pubsub.NewSchemaClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewSchemaClient: %w", err)
	}
	defer client.Close()

	var schemas []*pubsubpb.Schema

	req := &pubsubpb.ListSchemaRevisionsRequest{
		Name: fmt.Sprintf("projects/%s/schemas/%s", projectID, schemaID),
		View: pubsubpb.SchemaView_FULL,
	}
	schemaIter := client.ListSchemaRevisions(ctx, req)
	for {
		sc, err := schemaIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("schemaIter.Next: %w", err)
		}
		fmt.Fprintf(w, "Got schema revision: %#v\n", sc)
		schemas = append(schemas, sc)
	}

	fmt.Fprintf(w, "Got %d schema revisions", len(schemas))
	return schemas, nil
}

// [END pubsub_list_schema_revisions]
