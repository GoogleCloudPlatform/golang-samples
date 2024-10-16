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

package main

// [START datastore_update_entity]
import (
	"context"
	"log"

	"cloud.google.com/go/datastore"
)

// MarkDone marks the task done with the given ID.
func MarkDone(projectID string, databaseID string, taskID int64) error {
	ctx := context.Background()
	client, err := datastore.NewClientWithDatabase(ctx, projectID, databaseID)
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}
	defer client.Close()
	// Create a key using the given integer ID.
	key := datastore.IDKey("Task", taskID, nil)

	// In a transaction load each task, set done to true and store.
	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var task Task
		if err := tx.Get(key, &task); err != nil {
			return err
		}
		task.Done = true
		_, err := tx.Put(key, &task)
		return err
	})
	return err
}

// [END datastore_update_entity]
