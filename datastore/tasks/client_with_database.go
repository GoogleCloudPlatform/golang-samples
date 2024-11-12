// Copyright 2024 Google LLC
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

package main

// [START datastore_create_client_with_db]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/datastore"
)

// createClientWithDatabase creates a new client that references a
// custom database, e.g. not 'default'.
func createClientWithDatabase(w io.Writer, project, databaseID string) error {

	ctx := context.Background()
	client, err := datastore.NewClientWithDatabase(ctx, project, databaseID)
	if err != nil {
		return err
	}
	defer client.Close()

	// Sets the kind for the new entity.
	kind := "Task"
	// Sets the name/ID for the new entity.
	name := "sampletask1"
	// Creates a Key instance.
	taskKey := datastore.NameKey(kind, name, nil)

	// Creates a Task instance.
	task := struct {
		Description string
	}{
		Description: "Buy milk",
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, taskKey, &task); err != nil {
		return err
	}

	fmt.Fprintf(w, "Saved %v: %v\n", taskKey, task.Description)

	return nil
}

// [END datastore_create_client_with_db]
