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

// This App Engine application uses its default service account to list all
// the BigQuery datasets accessible via the BigQuery REST API.
package sample

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
)

func init() {
	http.HandleFunc("/", handle)
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Create a new App Engine context from the request.
	ctx := appengine.NewContext(r)

	// Get the list of dataset names.
	names, err := datasets(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	if len(names) == 0 {
		fmt.Fprintf(w, "No datasets visible")
	} else {
		fmt.Fprintf(w, "Datasets:\n\t"+strings.Join(names, "\n\t"))
	}
}

// datasets returns a list with the IDs of all the Big Query datasets visible
// with the given context.
func datasets(ctx context.Context) ([]string, error) {
	// Get the current application ID, which is the same as the project ID.
	projectID := appengine.AppID(ctx)

	// Create the BigQuery service.
	bq, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("could not create service: %v", err)
	}

	// Return a list of IDs.
	var ids []string
	it := bq.Datasets(ctx)
	for {
		ds, err := it.Next()
		if err == iterator.Done {
			return ids, nil
		} else if err != nil {
			return nil, err
		}
		ids = append(ids, ds.DatasetID)
	}
}
