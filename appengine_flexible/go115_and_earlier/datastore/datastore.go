// Copyright 2023 Google LLC
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

// [START gae_flex_datastore_app]

// Sample datastore demonstrates use of the cloud.google.com/go/datastore package from App Engine flexible.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"google.golang.org/appengine"
)

var datastoreClient *datastore.Client

func main() {
	ctx := context.Background()

	// Set this in app.yaml when running in production.
	projectID := os.Getenv("GCLOUD_DATASET_ID")

	var err error
	datastoreClient, err = datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ctx := context.Background()

	// Get a list of the most recent visits.
	visits, err := queryVisits(ctx, 10)
	if err != nil {
		msg := fmt.Sprintf("Could not get recent visits: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Record this visit.
	if err := recordVisit(ctx, time.Now(), r.RemoteAddr); err != nil {
		msg := fmt.Sprintf("Could not save visit: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Previous visits:")
	for _, v := range visits {
		fmt.Fprintf(w, "[%s] %s\n", v.Timestamp, v.UserIP)
	}
	fmt.Fprintln(w, "\nSuccessfully stored an entry of the current request.")
}

type visit struct {
	Timestamp time.Time
	UserIP    string
}

func recordVisit(ctx context.Context, now time.Time, userIP string) error {
	v := &visit{
		Timestamp: now,
		UserIP:    userIP,
	}

	k := datastore.IncompleteKey("Visit", nil)

	_, err := datastoreClient.Put(ctx, k, v)
	return err
}

func queryVisits(ctx context.Context, limit int64) ([]*visit, error) {
	// Print out previous visits.
	q := datastore.NewQuery("Visit").
		Order("-Timestamp").
		Limit(10)

	visits := make([]*visit, 0)
	_, err := datastoreClient.GetAll(ctx, q, &visits)
	return visits, err
}

// [END gae_flex_datastore_app]
