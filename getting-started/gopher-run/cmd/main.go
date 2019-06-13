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

//The cmd command starts a Gopher Run game server
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func main() {
	//Open http port
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "maralder-start")
	if err != nil {
		log.Fatalf("Error opening client, %v", err)
	}
	defer client.Close()
	//Get snapshots
	iter := client.Collection("teams").Snapshots(ctx)
	defer iter.Stop()
	for {
		docsnap, err := iter.Next()
		if err != nil {
			log.Fatalf("Error checking snapshots, %v", err)
		}
		fmt.Println(docsnap.Changes)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	projectID := "maralder-start"
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()
	//Read
	var d struct {
		Type  string `json:"type"`
		Name  string `json:"name"`
		Score int    `json:"score"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Error decoding JSON\n")
		printData(ctx, client, w)
		return
	}
	fmt.Fprintf(w, "Hello %s\n", html.EscapeString(d.Name))
	if d.Type == "add" {
		dsnap, err := client.Collection("teams").Doc(d.Name).Get(ctx)
		if err != nil {
			return
		}
		var s struct {
			Team  string
			Score int
		}
		dsnap.DataTo(&s)
		_, err = client.Collection("teams").Doc(d.Name).Set(ctx, map[string]interface{}{
			"name":  d.Name,
			"score": s.Score + d.Score,
		})
		if err != nil {
			log.Fatalf("Error setting data, %v", err)
		}
	}
	if d.Type == "set" {
		_, err = client.Collection("teams").Doc(d.Name).Set(ctx, map[string]interface{}{
			"name":  d.Name,
			"score": d.Score,
		})
		if err != nil {
			log.Fatalf("Error setting data, %v", err)
		}
	}
	if d.Type == "delete" {
		_, err := client.Collection("teams").Doc(d.Name).Delete(ctx)
		if err != nil {
			log.Fatalf("Error removing data, %v", err)
		}
	}
	printData(ctx, client, w)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	//Open database
	projectID := "maralder-start"
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Error creating client, %v", err)
	}
	defer client.Close()
	printData(ctx, client, w)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		handlePost(w, r)
	}
	if r.Method == "GET" {
		handleGet(w, r)
	}
}

func printData(ctx context.Context, client *firestore.Client, w io.Writer) {
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
