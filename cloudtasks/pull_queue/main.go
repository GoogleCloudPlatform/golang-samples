// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Program pull_queue provides a basic CLI for basic Cloud Tasks API interaction.
package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta2"
	duration "github.com/golang/protobuf/ptypes/duration"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta2"
)

// usage powers the help documentation when a help operation is used or the wrong number of arguments specified.
func usage() {
	fmt.Println("Usage of tasks_cli:")
	fmt.Println()
	fmt.Println("\t$> tasks create $PROJECT_ID $LOCATION_ID $QUEUE_ID")
	fmt.Println("\t$> tasks pull $PROJECT_ID $LOCATION_ID $QUEUE_ID")
	fmt.Println("\t$> tasks help")
	fmt.Println("\nFor more information, see https://cloud.google.com/cloud-tasks/docs")
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 || (len(args) < 3 && args[0] != "help") {
		usage()
		os.Exit(1)
	}

	switch args[0] {
	case "create":
		runTaskCreate(args[1:])
	case "pull":
		runTaskLeaseAndAck(args[1:])
	default:
		usage()
	}

	fmt.Println()
}

// runTaskCreate is invoked by the CLI for the "create" operation.
func runTaskCreate(args []string) {
	_, err := taskCreate(args[0], args[1], args[2])
	if err != nil {
		log.Fatalf("Error creating task: %s\n", err)
	}
}

// runTaskPull is invoked by the CLI for the "pull" operation.
func runTaskLeaseAndAck(args []string) {
	task, err := taskLease(args[0], args[1], args[2])
	if err != nil {
		log.Fatalf("Error leasing task: %s\n", err)
	}
	if task == nil {
		log.Println("No tasks available for lease")
	} else if taskAck(task) != nil {
		log.Fatalf("Error acknowledging task: %s\n", err)
	}
}

// [START cloud_tasks_create_task]

// taskCreate creates a new Task on the specified pull queue.
func taskCreate(projectID, locationID, queueID string) (*taskspb.Task, error) {
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

	fmt.Println("Created task:", createdTask.GetName())

	return createdTask, nil
}

// [END cloud_tasks_create_task]

// [START cloud_tasks_lease_and_acknowledge_task]

// runTaskPull leases the next task from the specified pull queue.
func taskLease(projectID, locationID, queueID string) (*taskspb.Task, error) {
	// Create a new Cloud Tasks client instance.
	// See https://godoc.org/cloud.google.com/go/cloudtasks/apiv2beta2
	ctx := context.Background()
	c, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewCloudTasksClient: %v", err)
	}

	// Construct the expected form of the Queue ID.
	queueName := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)

	// Cloud Tasks Go Client uses protobuf.
	// See https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2beta2#LeaseTasksRequest
	req := &taskspb.LeaseTasksRequest{
		Parent:        queueName,
		MaxTasks:      1,
		LeaseDuration: &duration.Duration{Seconds: 600},
		ResponseView:  taskspb.Task_FULL,
	}

	// See https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2beta2#LeaseTasksResponse
	resp, err := c.LeaseTasks(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("LeaseCloudTask: %v", err)
	}

	// If no tasks are available, nothing further to be done.
	if len(resp.Tasks) == 0 {
		return nil, nil
	}

	// Leasing tasks allows retrieval of one or more tasks. The Tasks property is always a slice
	// even if a single task is leased.
	leasedTask := resp.Tasks[0]

	// See the full code on Github for the implementation of toJsonString.
	fmt.Println("Leased task:", leasedTask.GetName())

	return leasedTask, nil
}

// taskAck acknowledges the provided Task for use in conjunction with taskLease().
func taskAck(task *taskspb.Task) error {
	// Create a new Cloud Tasks client instance.
	// See https://godoc.org/cloud.google.com/go/cloudtasks/apiv2beta2
	ctx := context.Background()
	c, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewCloudTasksClient: %v", err)
	}

	// See https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2beta2#AcknowledgeTaskRequest
	req := &taskspb.AcknowledgeTaskRequest{
		Name:         task.GetName(),
		ScheduleTime: task.GetScheduleTime(),
	}

	if err := c.AcknowledgeTask(ctx, req); err != nil {
		return fmt.Errorf("AcknowledgeCloudTask: %v", err)
	}
	fmt.Println("Acknowledged task:", task.GetName())

	return nil
}

// [END cloud_tasks_lease_and_acknowledge_task]
