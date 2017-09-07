// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	cloudtasks "google.golang.org/api/cloudtasks/v2beta2"
)

var (
	tasksService *cloudtasks.Service
)

func main() {
	ctx := context.Background()

	queueName := os.Args[1] // "projects/$PROJECT_ID/locations/$LOCATION_ID/queues/$QUEUE_ID"
	payload := os.Args[2]

	httpClient, err := google.DefaultClient(ctx, cloudtasks.CloudPlatformScope)
	if err != nil {
		log.Fatalf("Could not get HTTP client with app default credentials: %v", err)
	}
	tasksService, err = cloudtasks.New(httpClient)
	if err != nil {
		log.Fatalf("Could not initialize cloudtasks service: %v", err)
	}

	taskReq := &cloudtasks.CreateTaskRequest{
		Task: &cloudtasks.Task{
			AppEngineHttpRequest: &cloudtasks.AppEngineHttpRequest{
				RelativeUrl: "/payload",
				HttpMethod:  "POST",
				Payload:     base64.StdEncoding.EncodeToString([]byte(payload)),
			},
		},
	}

	_, err = tasksService.Projects.Locations.Queues.Tasks.Create(queueName, taskReq).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Couldn't add task: %v", err)
	}

	fmt.Println("Successfully added task.")
}
