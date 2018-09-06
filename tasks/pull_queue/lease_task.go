// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package snippets is a collection of sample code snippets.
package snippets

// [START cloud_tasks_lease_and_acknowledge_task]

import (
	"context"
	"fmt"
	"io"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta2"
	duration "github.com/golang/protobuf/ptypes/duration"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta2"
)

// runTaskPull leases the next task from the specified pull queue.
func taskLease(w io.Writer, projectID, locationID, queueID string) (*taskspb.Task, error) {
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
	fmt.Fprintln(w, "Leased task:", leasedTask.GetName())

	return leasedTask, nil
}

// taskAck acknowledges the provided Task for use in conjunction with taskLease().
func taskAck(w io.Writer, task *taskspb.Task) error {
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
	fmt.Fprintln(w, "Acknowledged task:", task.GetName())

	return nil
}

// [END cloud_tasks_lease_and_acknowledge_task]
