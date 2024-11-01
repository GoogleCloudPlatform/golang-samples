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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/getting-started/gopher-run/generator"
	"github.com/GoogleCloudPlatform/golang-samples/getting-started/gopher-run/leaderboard"
	"golang.org/x/oauth2/google"
)

type playData struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

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
	http.HandleFunc("/predict", a.predictionRequest)
	http.HandleFunc("/pldata", a.addPlayData)
	http.HandleFunc("/bggenerator", a.sendGeneratedBackground)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static/gorun/"))))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Starting server on port %v\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func (a *app) predictionRequest(w http.ResponseWriter, r *http.Request) {
	client, err := google.DefaultClient(r.Context(), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		log.Printf("DefaultClient: %v", err)
	}
	resp, err := client.Post("https://ml.googleapis.com/v1/projects/"+a.projectID+"/models/playerdata_linear_classification:predict", "application/json", r.Body)
	if err != nil {
		log.Printf("client.Post: %v", err)
		http.Error(w, "Prediction server error", http.StatusInternalServerError)
		return
	}
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

// Concatenate data from the top runs and start a Cloud ML training job on it
func (a *app) submitTrainingJob(ctx context.Context) error {
	bkt := a.bucket
	topPlayers, err := leaderboard.TopScores(ctx, a.fsClient)
	if err != nil {
		log.Printf("leaderboard.TopScores: %v", err)
		return err
	}
	// Add dummy data at the beginning for consistency in which actions are assigned numbers 0,1,2,3 for prediction output
	data := []byte(`idle,0,0,0,0,0,0,0,0,0,0,0
roll,0,0,0,0,0,0,0,0,0,0,0
jump,0,0,0,0,0,0,0,0,0,0,0
unroll,0,0,0,0,0,0,0,0,0,0,0`)
	for _, player := range topPlayers {
		pld, err := bkt.Object("pldata/" + player.Name + "_pldata.csv").NewReader(ctx)
		if err != nil {
			log.Printf("NewReader: %v", err)
			continue
		}
		defer pld.Close()
		old, err := bkt.Object("pldata.csv").NewReader(ctx)
		if err != nil {
			log.Printf("NewReader: %v", err)
			continue
		}
		defer old.Close()
		b, err := ioutil.ReadAll(pld)
		if err != nil {
			log.Printf("ioutil.ReadAll: %v", err)
			continue
		}
		data = append(data, []byte("\n")...))
		data = append(data, b...)
	}
	newObj := bkt.Object("pldata.csv").NewWriter(ctx)
	defer newObj.Close()
	newObj.Write(data)
	return nil
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
	top, err := leaderboard.AddScore(r.Context(), a.fsClient, d)
	if err != nil {
		log.Printf("leaderboar.AddScore: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
	a.submitTrainingJob(r.Context())
	fmt.Fprint(w, top)
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

func (a *app) addPlayData(w http.ResponseWriter, r *http.Request) {
	var d playData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&d); err != nil {
		fmt.Fprintf(w, "decoder.Decode: %v", err)
	}
	r.Body.Close()
	new := a.bucket.Object("pldata/" + d.Name + "_pldata.csv").NewWriter(r.Context())
	defer new.Close()
	new.Write([]byte(d.Value))
	fmt.Fprint(w, "Recieved data\n")
}

// sendGeneratedBackground returns cloud/hill placements.
func (a *app) sendGeneratedBackground(w http.ResponseWriter, r *http.Request) {
	var d generator.RequestData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&d); err != nil {
		fmt.Fprintf(w, "decoder.Decode: %v\n", err)
		log.Printf("decoder.Decode: %v\n", err)
		return
	}
	r.Body.Close()
	objs := generator.GenerateBackground(d.Xmin, d.Xmax, d.Speed)
	s := ""
	for _, obj := range objs {
		s += obj.String() + "\n"
	}
	fmt.Fprint(w, s)
}

func newApp(projectID string) (*app, error) {
	ctx := context.Background()
	fsClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %w", err)
	}
	csClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	bName := os.Getenv("GOPHER_RUN_BUCKET")
	if bName == "" {
		return nil, fmt.Errorf("env variable GOPHER_RUN_BUCKET must be set")
	}
	bucket := csClient.Bucket(bName)
	return &app{projectID: projectID, fsClient: fsClient, bucket: bucket}, nil
}
