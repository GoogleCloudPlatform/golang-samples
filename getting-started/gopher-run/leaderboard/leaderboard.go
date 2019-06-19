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

// Package leaderboard starts a Gopher Run leaderboard server.
package leaderboard

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type scoredata struct {
	Name     string  `json:"name"`
	Team     string  `json:"team"`
	Coins    int     `json:"coins"`
	Distance float32 `json:"distance"`
}

// Handler chooses a handler function based on the request method.
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		handlePost(w, r)
	} else if r.Method == "GET" {
		handleGet(w, r)
	}
}

// handlePost adds a new score to the database.
func handlePost(w http.ResponseWriter, r *http.Request) {
	projectID := "maralder-start"
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()
	//Read
	var d scoredata
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&d)
	if err != nil {
		fmt.Fprint(w, "Error decoding JSON\n")
	}
	fmt.Fprint(w, "Act: "+d.Name)
	_, _, err = client.Collection("leaderboard").Add(ctx, map[string]interface{}{
		"name":     d.Name,
		"team":     d.Team,
		"coins":    d.Coins,
		"distance": d.Distance,
	})
	if err != nil {
		log.Fatalf("Error setting data, %v", err)
	}
}

// handleGet retrieves the top scores from the database.
func handleGet(w http.ResponseWriter, r *http.Request) {
	projectID := "maralder-start"
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Error creating client, %v", err)
	}
	defer client.Close()
	iter := client.Collection("teams").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed iteration %v", err)
		}
		fmt.Fprint(w, doc.Data())
	}
}
