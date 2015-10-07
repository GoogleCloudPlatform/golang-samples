// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// This App Engine application uses its default service account to list all
// the BigQuery projects accessible via the BigQuery REST API.
package sample

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	// all requests are handled by handler.
	http.HandleFunc("/", handle)
}

func handle(w http.ResponseWriter, r *http.Request) {
	// create a new App Engine context from the request.
	c := appengine.NewContext(r)

	// obtain the list of project names.
	names, err := projects(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// print it to the output.
	fmt.Fprintln(w, strings.Join(names, "\n"))
}

// projects returns a list with the names of all the Big Query projects visible
// with the given context.
func projects(c context.Context) ([]string, error) {
	// create a new authenticated HTTP client over urlfetch.
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(c, bigquery.BigqueryScope),
			Base:   &urlfetch.Transport{Context: c},
		},
	}

	// create the BigQuery service.
	bq, err := bigquery.New(client)
	if err != nil {
		return nil, fmt.Errorf("create service: %v", err)
	}

	// list the projects.
	list, err := bq.Projects.List().Do()
	if err != nil {
		return nil, fmt.Errorf("list projects: %v", err)
	}

	// prepare the list of names.
	var names []string
	for _, p := range list.Projects {
		names = append(names, p.FriendlyName)
	}
	return names, nil
}
