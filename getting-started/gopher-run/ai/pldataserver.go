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

//This package starts a Gopher Run player data server
package pldata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const projectID := "maralder-start"
//os.Getenv("GOOGLE_CLOUD_PROJECT")

type playdata struct {
	Act string  `json:"act"`
	H   float64 `json:h`
	X0  float64 `json:x0`
	Y0  float64 `json:y0`
	X1  float64 `json:x1`
	Y1  float64 `json:y1`
	X2  float64 `json:x2`
	Y2  float64 `json:y2`
	X3  float64 `json:x3`
	Y3  float64 `json:y3`
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Printf("firestore.NewClient: %v", err)
	}
	defer client.Close()
	//Read
	var d playdata
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&d)
	if err != nil {
		fmt.Fprint(w, "Error decoding JSON\n")
	}
	fmt.Fprint(w, "Act: "+d.Act)
	_, _, err = client.Collection("patterns").Add(ctx, map[string]interface{}{
		"act": d.Act,
		"h":   d.H,
		"x0":  d.X0,
		"y0":  d.Y0,
		"x1":  d.X1,
		"y1":  d.Y1,
		"x2":  d.X2,
		"y2":  d.Y2,
		"x3":  d.X3,
		"y3":  d.Y3,
	})
	if err != nil {
		log.Printf("Error setting data, %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		handlePost(w, r)
	} else {
		log.Printf("Unexpected request method: %v", r.Method)
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

func main() {
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
