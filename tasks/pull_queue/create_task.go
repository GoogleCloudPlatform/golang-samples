// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package snippets is a collection of sample code snippets.
package snippets

// [START cloud_tasks_create_task]

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta2"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta2"
)

// taskCreate creates a new Task on the specified pull queue.
func taskCreate(w io.Writer, projectID, locationID, queueID string) (*taskspb.Task, error) {
	// Create a new Cloud Tasks client instance.
	// See https://godoc.org/cloud.google.com/go/cloudtasks/apiv2beta2
	ctx := context.Background()
	c, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewCloudTasksClient: %v", err)
	}

	// Message to be sent as the task payload.
	message := base64.StdEncoding.EncodeToString([]byte("a message for the recipient"))

	// Construct the expected form of the Queue ID.
	queueName := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)

	// Cloud Tasks Go Client uses protobuf.
	// See https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2beta2#CreateTaskRequest
	req := &taskspb.CreateTaskRequest{
		Parent: queueName,
		Task: &taskspb.Task{
			PayloadType: &taskspb.Task_PullMessage{
				PullMessage: &taskspb.PullMessage{
					Payload: []byte(message),
				},
			},
		},
	}
	createdTask, err := c.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("CreateCloudTask: %v", err)
	}

	fmt.Fprintln(w, "Created task:", createdTask.GetName())

	return createdTask, nil
}

// [END cloud_tasks_create_task]
