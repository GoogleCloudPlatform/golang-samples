// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/taskqueue"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// [START adding_tasks_to_a_pull_queue]
	t := &taskqueue.Task{
		Payload: []byte("hello world"),
		Method:  "PULL",
	}
	_, err := taskqueue.Add(ctx, t, "pull-queue")
	// [END adding_tasks_to_a_pull_queue]
	_ = err

	// [START leasing_tasks_1]
	tasks, err := taskqueue.Lease(ctx, 100, "pull-queue", 3600)
	// [END leasing_tasks_1]

	// [START leasing_tasks_2]
	_, err = taskqueue.Add(ctx, &taskqueue.Task{
		Payload: []byte("parse"), Method: "PULL", Tag: "parse",
	}, "pull-queue")
	_, err = taskqueue.Add(ctx, &taskqueue.Task{
		Payload: []byte("render"), Method: "PULL", Tag: "render",
	}, "pull-queue")

	// leases render tasks, but not parse
	tasks, err = taskqueue.LeaseByTag(ctx, 100, "pull-queue", 3600, "render")

	// Leases up to 100 tasks that have same tag.
	// Tag is that of "oldest" task by ETA.
	tasks, err = taskqueue.LeaseByTag(ctx, 100, "pull-queue", 3600, "")
	// [END leasing_tasks_2]

	// [START deleting_tasks_1]
	tasks, err = taskqueue.Lease(ctx, 100, "pull-queue", 3600)
	// Perform some work with the tasks here

	taskqueue.DeleteMulti(ctx, tasks, "pull-queue")
	// [END deleting_tasks_1]

}
