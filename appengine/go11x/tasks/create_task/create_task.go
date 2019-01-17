// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START cloud_tasks_appengine_create_task]

// Command create_task constructs and adds a task to an App Engine Queue.
package main

import (
	"context"
	"fmt"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta3"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta3"
)

// createTask creates a new task in your App Engine queue.
func createTask(projectID, locationID, queueID, message string) (*taskspb.Task, error) {
	// Create a new Cloud Tasks client instance.
	// See https://godoc.org/cloud.google.com/go/cloudtasks/apiv2beta3
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %v", err)
	}

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)

	// Build the Task payload.
	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2beta3#CreateTaskRequest
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2beta3#AppEngineHttpRequest
			PayloadType: &taskspb.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &taskspb.AppEngineHttpRequest{
					HttpMethod:  taskspb.HttpMethod_POST,
					RelativeUri: "/task_handler",
				},
			},
		},
	}

	// Add a payload message if one is present.
	if message != "" {
		req.Task.GetAppEngineHttpRequest().Body = []byte(message)
	}

	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %v", err)
	}

	return createdTask, nil
}

// [END cloud_tasks_appengine_create_task]
