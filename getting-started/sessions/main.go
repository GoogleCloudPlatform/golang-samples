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

// Command sessions starts an HTTP server that uses session state.
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	firestoregorilla "github.com/GoogleCloudPlatform/firestore-gorilla-go"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// app stores a sessions.Store. Create a new app with newApp.
type app struct {
	store sessions.Store
}

// colors are the random background colors that will be assigned to sessions.
var colors = []string{"red", "blue", "green", "yellow", "pink"}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal("GOOGLE_CLOUD_PROJECT must be set")
	}

	a, err := newApp(projectID)
	if err != nil {
		log.Fatalf("newApp: %v", err)
	}

	http.HandleFunc("/", a.index)

	log.Printf("Listening on :%v", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// newApp creates a new app.
func newApp(projectID string) (*app, error) {
	// For this sample, hashKey and blockKey are created using
	// GenerateRandomKey(), which does not automatically persist they keys.
	// New app instances will generate new keys, which will not be able to
	// decode cookies issued by other instances.
	//
	// Set these to random values and treat them as secrets (not hard-coded in
	// your source code).
	hashKey := securecookie.GenerateRandomKey(32)
	if hashKey == nil {
		log.Fatal("Failed to generate hashKey")
	}
	blockKey := securecookie.GenerateRandomKey(32)
	if blockKey == nil {
		log.Fatal("Failed to generate blockKey")
	}
	codecs := securecookie.CodecsFromPairs(hashKey, blockKey)

	ctx := context.Background()
	store, err := firestoregorilla.New(ctx, projectID, codecs...)
	if err != nil {
		log.Fatalf("firestoregorilla.New: %v", err)
	}

	return &app{store: store}, nil
}

// index uses sessions to assign users a random color and keep track of views.
func (a *app) index(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		return
	}

	session, err := a.store.Get(r, "my-sessions-name")
	if err != nil {
		log.Printf("store.Get: %v", err)
		http.Error(w, "Unable to get session", http.StatusInternalServerError)
		return
	}

	if session.IsNew {
		session.Values["views"] = 0
		session.Values["color"] = colors[rand.Intn(len(colors))]
	}
	session.Values["views"] = session.Values["views"].(int) + 1
	if err := session.Save(r, w); err != nil {
		log.Printf("Save: %v", err)
		// Don't return early so the user still gets a response.
	}

	fmt.Fprintf(w, "<body bgcolor=%v>Views %v</body>", session.Values["color"], session.Values["views"])
}
