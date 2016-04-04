// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package counter

// [START push_queues_and_backends]
import (
	"net/http"
	"net/url"

	"google.golang.org/appengine"
	"google.golang.org/appengine/taskqueue"
)

func pushHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	key := r.FormValue("key")

	// Create a task pointed at a backend.
	t := taskqueue.NewPOSTTask("/path/to/my/worker/", url.Values{
		"key": {key},
	})
	host, err := appengine.ModuleHostname(ctx, "backend1", "", "")
	if err != nil {
		// Handle err
	}
	t.Header = http.Header{
		"Host": {host},
	}

	// Add the task to the default queue.
	if _, err := taskqueue.Add(ctx, t, ""); err != nil {
		// Handle err
	}
}

// [END push_queues_and_backends]
