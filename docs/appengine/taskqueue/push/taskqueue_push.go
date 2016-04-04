// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START intro]
package counter

import (
	"html/template"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/worker", worker)
}

type Counter struct {
	Name  string
	Count int
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if name := r.FormValue("name"); name != "" {
		t := taskqueue.NewPOSTTask("/worker", map[string][]string{"name": {name}})
		if _, err := taskqueue.Add(ctx, t, ""); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	q := datastore.NewQuery("Counter")
	var counters []Counter
	if _, err := q.GetAll(ctx, &counters); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := handlerTemplate.Execute(w, counters); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// OK
}

func worker(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	name := r.FormValue("name")
	key := datastore.NewKey(ctx, "Counter", name, 0, nil)
	var counter Counter
	if err := datastore.Get(ctx, key, &counter); err == datastore.ErrNoSuchEntity {
		counter.Name = name
	} else if err != nil {
		log.Errorf(ctx, "%v", err)
		return
	}
	counter.Count++
	if _, err := datastore.Put(ctx, key, &counter); err != nil {
		log.Errorf(ctx, "%v", err)
	}
}

var handlerTemplate = template.Must(template.New("handler").Parse(handlerHTML))

const handlerHTML = `
{{repeat .}}
<p>{{.Name}}: {{.Count}}</p>
{{end}}
<p>Start a new counter:</p>
<form action="/" method="POST">
<input type="text" name="name">
<input type="submit" value="Add">
</form>
`

// [END intro]

func example() {
	var ctx context.Context
	var t *taskqueue.Task
	_ = t

	// [START deferred_tasks]
	var expensiveFunc = delay.Func("some-arbitrary-key", func(ctx context.Context, a string, b int) {
		// do something expensive!
	})

	// Somewhere else
	expensiveFunc.Call(ctx, "Hello, world!", 42)
	// [END deferred_tasks]

	// [START URL_endpoints]
	t = &taskqueue.Task{Path: "/path/to/my/worker"}
	t = &taskqueue.Task{Path: "/path?a=b&c=d", Method: "GET"}
	// [END URL_endpoints]
}
