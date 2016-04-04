// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

// [START using_namespaces_with_the_Task_Queue]
import (
	"io"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/taskqueue"
)

type Counter struct {
	Count int64
}

func incrementCounter(ctx context.Context, name string) error {
	key := datastore.NewKey(ctx, "Counter", name, 0, nil)
	return datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		var ctr Counter
		err := datastore.Get(ctx, key, &ctr)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		ctr.Count++
		_, err = datastore.Put(ctx, key, &ctr)
		return err
	}, nil)
}

// taskQueueHandler serves /_ah/counter.
func taskQueueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	ctx := appengine.NewContext(r)
	err := incrementCounter(ctx, r.FormValue("counter_name"))
	if err != nil {
		// ... handle err
	}
}

func someRequest(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// Perform asynchronous requests to update counter.
	// (missing error handling here.)
	t := taskqueue.NewPOSTTask("/_ah/counter", map[string][]string{
		"counter_name": {"someRequest"},
	})

	taskqueue.Add(ctx, t, "")

	// temporarily use a new namespace
	{
		ctx, err := appengine.Namespace(ctx, "-global-")
		if err != nil {
			// ... handle err
		}
		taskqueue.Add(ctx, t, "")
	}

	io.WriteString(w, "Counters will be updated.\n")
}

func init() {
	http.HandleFunc("/_ah/counter", taskQueueHandler)
	http.HandleFunc("/some_request", someRequest)
}

// [END using_namespaces_with_the_Task_Queue]
