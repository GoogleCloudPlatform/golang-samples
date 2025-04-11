// Copyright 2020 Google LLC
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

package samples

// [START datastore_admin_client_create]
import (
	"context"
	"fmt"
	"io"

	admin "cloud.google.com/go/datastore/admin/apiv1"
)

// clientCreate creates a new Datastore admin client.
func clientCreate(w io.Writer) (*admin.DatastoreAdminClient, error) {
	ctx := context.Background()
	client, err := admin.NewDatastoreAdminClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("admin.NewDatastoreAdminClient: %w", err)
	}
	// Close client when done using it.
	// defer client.Close()
	fmt.Fprintf(w, "Admin client created\n")
	return client, nil
}

// [END datastore_admin_client_create]
