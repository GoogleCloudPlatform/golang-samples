// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START sample_code]
// This sample gets the app displays 5 log Records at a time, including all
// AppLogs, with a Next link to let the user page through the results using the
// Record's Offset property.
package app

import (
	"encoding/base64"
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	http.HandleFunc("/", handler)
}

const recordsPerPage = 5

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// Set up a data structure to pass to the HTML template.
	var data struct {
		Records []*log.Record
		Offset  string // base-64 encoded string
	}

	// Set up a log.Query.
	query := &log.Query{AppLogs: true}

	// Get the incoming offset param from the Next link to advance through
	// the logs. (The first time the page is loaded there won't be any offset.)
	if offset := r.FormValue("offset"); offset != "" {
		query.Offset, _ = base64.URLEncoding.DecodeString(offset)
	}

	// Run the query, obtaining a Result iterator.
	res := query.Run(ctx)

	// Iterate through the results populating the data struct.
	for i := 0; i < recordsPerPage; i++ {
		rec, err := res.Next()
		if err == log.Done {
			break
		}
		if err != nil {
			log.Errorf(ctx, "Reading log records: %v", err)
			break
		}

		data.Records = append(data.Records, rec)
		if i == recordsPerPage-1 {
			data.Offset = base64.URLEncoding.EncodeToString(rec.Offset)
		}
	}

	// Render the template to the HTTP response.
	if err := tmpl.Execute(w, data); err != nil {
		log.Errorf(ctx, "Rendering template: %v", err)
	}
}

var tmpl = template.Must(template.New("").Parse(`
	{{range .Records}}
		<h2>Request Log</h2>
		<p>{{.EndTime}}: {{.IP}} {{.Method}} {{.Resource}}</p>
		{{with .AppLogs}}
			<h3>App Logs:</h3>
			<ul>
			{{range .}}
				<li>{{.Time}}: {{.Message}}</li>
			<{{end}}
			</ul>
		{{end}}
	{{end}}
	{{with .Offset}}
		<a href="?offset={{.}}">Next</a>
	{{end}}
`))

// [END sample_code]
