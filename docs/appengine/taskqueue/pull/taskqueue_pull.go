// Copyright 2019 Google LLC
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

package sample

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/taskqueue"
)

func addTaskHandler(_ http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// [START gae_adding_tasks_to_pull_queue]
	t := &taskqueue.Task{
		Payload: []byte("hello world"),
		Method:  "PULL",
	}
	_, err := taskqueue.Add(ctx, t, "pull-queue")
	// [END gae_adding_tasks_to_a_pull_queue]
	_ = err

	// [START gae_leasing_tasks_1]
	tasks, err := taskqueue.Lease(ctx, 100, "pull-queue", 3600)
	// [END gae_leasing_tasks_1]

	// [START gae_leasing_tasks_2]
	_, err = taskqueue.Add(ctx, &taskqueue.Task{
		Payload: []byte("parse"), Method: "PULL", Tag: "parse",
	}, "pull-queue")
	_, err = taskqueue.Add(ctx, &taskqueue.Task{
		Payload: []byte("render"), Method: "PULL", Tag: "render",
	}, "pull-queue")

	// leases render tasks, but not parse
	tasks, err = taskqueue.LeaseByTag(ctx, 100, "pull-queue", 3600, "render")

	// Leases up to 100 tasks that have same tag.
	tasks, err = taskqueue.LeaseByTag(ctx, 100, "pull-queue", 3600, "")
	// [END gae_leasing_tasks_2]

	// [START gae_deleting_tasks_1]
	tasks, err = taskqueue.Lease(ctx, 100, "pull-queue", 3600)
	// Perform some work with the tasks here

	taskqueue.DeleteMulti(ctx, tasks, "pull-queue")
	// [END gae_deleting_tasks_1]

}
