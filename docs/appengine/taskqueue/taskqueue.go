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

// [START gae_tasks_within_transactions]
import (
	"context"
	"net/http"
	"net/url"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/taskqueue"
)

func f(ctx context.Context) {
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		t := taskqueue.NewPOSTTask("/worker", url.Values{
			// ...
		})
		// Use the transaction's context when invoking taskqueue.Add.
		_, err := taskqueue.Add(ctx, t, "")
		if err != nil {
			// Handle error
		}
		// ...
		return nil
	}, nil)
	if err != nil {
		// Handle error
	}
	// ...
}

// [END gae_tasks_within_transactions]

func example() {
	var ctx context.Context

	// [START gae_purging_tasks]
	// Purge entire queue...
	err := taskqueue.Purge(ctx, "queue1")
	// [END gae_purging_tasks]

	// [START gae_deleting_tasks]
	// Delete an individual task...
	t := &taskqueue.Task{Name: "foo"}
	err = taskqueue.Delete(ctx, t, "queue1")
	// [END gae_deleting_tasks]
	_ = err

	// [START gae_taskqueue_host]
	h := http.Header{}
	h.Add("Host", "versionHostname")
	task := taskqueue.Task{
		Header: h,
	}
	// [END gae_taskqueue_host]
	_ = task
}
