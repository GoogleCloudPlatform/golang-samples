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

// [START pubsub_delete_schema]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsub"
)

func deleteSchema(w io.Writer, projectID, schemaID string) error {
	// projectID := "my-project-id"
	// schemaID := "my-schema"
	ctx := context.Background()
	client, err := pubsub.NewSchemaClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewSchemaClient: %w", err)
	}
	defer client.Close()

	if err := client.DeleteSchema(ctx, schemaID); err != nil {
		return fmt.Errorf("client.DeleteSchema: %w", err)
	}
	fmt.Fprintf(w, "Deleted schema: %s", schemaID)
	return nil
}

// [END pubsub_delete_schema]
