// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START datastore_quickstart]
// Sample datastore-quickstart fetches an entity from Google Cloud Datastore.
package main

import (
	"fmt"
	"log"

	// Imports the Google Cloud Datastore client package.
	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
)

type Task struct {
	Description string
}

func main() {
	ctx := context.Background()

	// Setx your Google Cloud Platform project ID.
	projectID := "YOUR_PROJECT_ID"

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the kind for the new entity.
	kind := "Task"
	// Sets the name/ID for the new entity.
	name := "sampletask1"
	// Creates a Key instance.
	taskKey := datastore.NameKey(kind, name, nil)

	// Creates a Task instance.
	task := Task{
		Description: "Buy milk",
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, taskKey, &task); err != nil {
		log.Fatalf("Failed to save task: %v", err)
	}

	fmt.Printf("Saved %v: %v\n", taskKey, task.Description)
}

// [END datastore_quickstart]
