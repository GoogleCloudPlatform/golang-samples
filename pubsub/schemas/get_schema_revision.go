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

// [START pubsub_get_schema_revision]
import (
	"context"
	"fmt"
	"io"

	pubsub "cloud.google.com/go/pubsub/v2/apiv1"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
)

func getSchemaRevision(w io.Writer, projectID, schemaID string) error {
	// projectID := "my-project-id"
	// schemaID := my-schema@c7cfa2a8 // with revision
	ctx := context.Background()
	client, err := pubsub.NewSchemaClient(ctx)
	if err != nil {
		return fmt.Errorf("pubsub.NewSchemaClient: %w", err)
	}
	defer client.Close()

	req := &pubsubpb.GetSchemaRequest{
		Name: fmt.Sprintf("projects/%s/schemas/%s", projectID, schemaID),
		View: pubsubpb.SchemaView_FULL,
	}
	s, err := client.GetSchema(ctx, req)
	if err != nil {
		return fmt.Errorf("client.GetSchema revision: %w", err)
	}
	fmt.Fprintf(w, "Got schema revision: %#v\n", s)
	return nil
}

// [END pubsub_get_schema_revision]
