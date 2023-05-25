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

// [START datastore_add_entity]
import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/datastore"
)

// Task is the model used to store tasks in the datastore.
type Task struct {
	Desc    string    `datastore:"description"`
	Created time.Time `datastore:"created"`
	Done    bool      `datastore:"done"`
	id      int64     // The integer ID used in the datastore.
}

// AddTask adds a task with the given description to the datastore,
// returning the key of the newly created entity.
func AddTask(projectID string, desc string) (*datastore.Key, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}
	defer client.Close()
	task := &Task{
		Desc:    desc,
		Created: time.Now(),
	}
	key := datastore.IncompleteKey("Task", nil)
	return client.Put(ctx, key, task)

}
