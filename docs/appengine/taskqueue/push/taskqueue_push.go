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

// [START intro]

package counter

import (
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
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
{{range .}}
<p>{{.Name}}: {{.Count}}</p>
{{end}}
<p>Start a new counter:</p>
<form action="/" method="POST">
<input type="text" name="name">
<input type="submit" value="Add">
</form>
`

// [END intro]
