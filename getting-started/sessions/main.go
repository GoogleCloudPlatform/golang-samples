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

// [START getting_started_sessions_setup]

// Command sessions starts an HTTP server that uses session state.
package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firestoregorilla "github.com/GoogleCloudPlatform/firestore-gorilla-sessions"
	"github.com/gorilla/sessions"
)

// app stores a sessions.Store. Create a new app with newApp.
type app struct {
	store sessions.Store
	tmpl  *template.Template
}

// greetings are the random greetings that will be assigned to sessions.
var greetings = []string{
	"Hello World",
	"Hallo Welt",
	"Ciao Mondo",
	"Salut le Monde",
	"Hola Mundo",
}

// [END getting_started_sessions_setup]

// [START getting_started_sessions_main]

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

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// newApp creates a new app.
func newApp(projectID string) (*app, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	store, err := firestoregorilla.New(ctx, client)
	if err != nil {
		log.Fatalf("firestoregorilla.New: %v", err)
	}

	tmpl, err := template.New("Index").Parse(`<body>{{.views}} {{if eq .views 1.0}}view{{else}}views{{end}} for "{{.greeting}}"</body>`)
	if err != nil {
		return nil, fmt.Errorf("template.New: %w", err)
	}

	return &app{
		store: store,
		tmpl:  tmpl,
	}, nil
}

// [END getting_started_sessions_main]

// [START getting_started_sessions_handler]

// index uses sessions to assign users a random greeting and keep track of
// views.
func (a *app) index(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		return
	}

	// name is a non-empty identifier for this app's sessions. Set it to
	// something descriptive for your app. It is used as the Firestore
	// collection name that stores the sessions.
	name := "hello-views"
	session, err := a.store.Get(r, name)
	if err != nil {
		// Could not get the session. Log an error and continue, saving a new
		// session.
		log.Printf("store.Get: %v", err)
	}

	if session.IsNew {
		// firestoregorilla uses JSON, which unmarshals numbers as float64s.
		session.Values["views"] = float64(0)
		session.Values["greeting"] = greetings[rand.Intn(len(greetings))]
	}
	session.Values["views"] = session.Values["views"].(float64) + 1
	if err := session.Save(r, w); err != nil {
		log.Printf("Save: %v", err)
		// Don't return early so the user still gets a response.
	}

	if err := a.tmpl.Execute(w, session.Values); err != nil {
		log.Printf("Execute: %v", err)
	}
}

// [END getting_started_sessions_handler]
