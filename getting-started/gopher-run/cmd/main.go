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

// The cmd command starts a Gopher Run leaderboard server
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/getting-started/gopher-run/leaderboard"
)

type app struct {
	projectID string
	bucket    *storage.BucketHandle
	fsClient  *firestore.Client
}

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal("No GOOGLE_CLOUD_PROJECT variable")
	}
	a, err := newApp(projectID)
	if err != nil {
		log.Fatalf("newApp: %v", err)
	}
	http.HandleFunc("/leaderboard/post", a.addScore)
	http.HandleFunc("/leaderboard/get", a.topScores)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Starting server: localhost:%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func (a *app) addScore(w http.ResponseWriter, r *http.Request) {
	var d leaderboard.ScoreData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&d); err != nil {
		log.Printf("decoder.Decode: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	r.Body.Close()
	fmt.Fprint(w, "Act: "+d.Name)
	if err := leaderboard.AddScore(r.Context(), a.fsClient, d); err != nil {
		log.Printf("leaderboar.AddScore: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}

// topScores retrieves top 10 scores from the database, return as \n-separated jsons.
func (a *app) topScores(w http.ResponseWriter, r *http.Request) {
	scores, err := leaderboard.TopScores(r.Context(), a.fsClient)
	if err != nil {
		log.Printf("leaderboard.TopScores: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	for _, obj := range scores {
		j, err := json.Marshal(obj)
		if err != nil {
			log.Printf("json.Marshal: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%v\n", string(j))
	}
}

func newApp(projectID string) (*app, error) {
	ctx := context.Background()
	fsClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %v", err)
	}
	csClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	bName := os.Getenv("GOPHER_RUN_BUCKET")
	if bName == "" {
		return nil, fmt.Errorf("env variable GOPHER_RUN_BUCKET must be set")
	}
	bucket := csClient.Bucket(bName)
	return &app{projectID: projectID, fsClient: fsClient, bucket: bucket}, nil
}
