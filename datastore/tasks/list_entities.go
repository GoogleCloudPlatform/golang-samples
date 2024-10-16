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

// [START datastore_retrieve_entities]
import (
	"context"
	"log"

	"cloud.google.com/go/datastore"
)

// ListTasks returns all the tasks in ascending order of creation time.
func ListTasks(projectID string, databaseID string) ([]*Task, error) {
	ctx := context.Background()
	//client, err := datastore.NewClient(ctx, projectID)
	client, err := datastore.NewClientWithDatabase(ctx, projectID, databaseID)
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}

	var tasks []*Task
	// Create a query to fetch all Task entities, ordered by "created".
	query := datastore.NewQuery("Task").Order("created")
	keys, err := client.GetAll(ctx, query, &tasks)
	if err != nil {
		return nil, err
	}

	// Set the id field on each Task from the corresponding key.
	for i, key := range keys {
		tasks[i].id = key.ID
	}

	return tasks, nil
}

// [END datastore_retrieve_entities]
