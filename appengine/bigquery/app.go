// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// This App Engine application uses its default service account to list all
// the BigQuery datasets accessible via the BigQuery REST API.
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
	http.HandleFunc("/", handle)
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// create a new App Engine context from the request.
	ctx := appengine.NewContext(r)

	// obtain the list of dataset names.
	names, err := datasets(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text")

	if len(names) == 0 {
		fmt.Fprintf(w, "no datasets visible")
	} else {
		fmt.Fprintf(w, "datasets:\n\t"+strings.Join(names, "\n\t"))
	}
}

// datasets returns a list with the ids of all the Big Query datasets visible
// with the given context.
func datasets(ctx context.Context) ([]string, error) {
	// create a new authenticated HTTP client over urlfetch.
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(ctx, bigquery.BigqueryScope),
			Base:   &urlfetch.Transport{Context: ctx},
		},
	}

	// create the BigQuery service.
	bq, err := bigquery.New(client)
	if err != nil {
		return nil, fmt.Errorf("could not create service: %v", err)
	}

	// obtain the current application id, the BigQuery id is the same.
	appID := appengine.AppID(ctx)

	// prepare the list of ids.
	var ids []string
	datasets, err := bq.Datasets.List(appID).Do()
	if err != nil {
		return nil, fmt.Errorf("could not list datasets for %q: %v", appID, err)
	}
	for _, d := range datasets.Datasets {
		ids = append(ids, d.Id)
	}
	return ids, nil
}
