// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Simple CLI to run the createTask function via the Cloud Tasks API.
// Used for one-off testing and development.

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 4 {
		fmt.Println("Usage: Must include 3 arguments for projectID, locationID, and queueID")
		os.Exit(1)
	}
	projectID := os.Args[1]
	locationID := os.Args[2]
	queueID := os.Args[3]

	message := ""
	if len(os.Args) > 4 {
		message = os.Args[4]
	}

	task, err := createTask(projectID, locationID, queueID, message)
	if err != nil {
		log.Fatalf("createTask: %v", err)
	}

	fmt.Printf("Create Task: %s\n", task.GetName())
}
