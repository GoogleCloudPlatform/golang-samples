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
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
)

// app stores a sessions.Store. Create a new app with newApp.
type app struct {
	client *firestore.Client
	tmpl   *template.Template
}

type session struct {
	Greetings string `json:"greeting"`
	Views     int    `json:"views"`
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

	tmpl, err := template.New("Index").Parse(`<body>{{.views}} {{if eq .views 1}}view{{else}}views{{end}} for "{{.greeting}}"</body>`)
	if err != nil {
		return nil, fmt.Errorf("template.New: %w", err)
	}

	return &app{
		client: client,
		tmpl:   tmpl,
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

	ctx := context.Background()

	var id string
	var session session
	var cookie *http.Cookie
	var isNewSession bool

	// collectionName and cookieName are a non-empty identifiers for this app's sessions. Set them to
	// something descriptive for your app.
	//
	// collectionName is used as the Firestore
	// collection name that stores the sessions.
	//
	// cookieName is used as the Key
	// name that contains the session's id value.

	collectionName := "hello-views"
	cookieName := "session_id"

	// If err is different to nil, it means the cookie has not been set, so it will be created.
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			// isNewSession flag is set to true
			isNewSession = true
		} else {
			log.Printf("Error getting cookie: %v", err)
			http.Error(w, "Error accessing session", http.StatusInternalServerError)
			return
		}
	}

	// If isNewSession flag is true, the session will be created
	if isNewSession {
		// Get unique id for new document
		id = a.client.Collection(collectionName).NewDoc().ID

		session.Greetings = greetings[rand.Intn(len(greetings))]
		session.Views = 1

		// New document is created
		_, err := a.client.Collection(collectionName).Doc(id).Set(ctx, session)
		if err != nil {
			log.Printf("client.Collection.Doc.Set error: %v", err)
			// Don't return early so the user still gets a response.
		}

		// Cookie is set
		cookie = &http.Cookie{
			Name:  cookieName,
			Value: id,
		}
		http.SetCookie(w, cookie)
	} else {
		// The session exists

		// Get session
		doc, err := a.client.Collection(collectionName).Doc(cookie.Value).Get(ctx)
		if err != nil {
			log.Printf("client.Collection.Doc.Get error: %v", err)
			// Don't return early so the user still gets a response.
		}

		// Unmarshal documents's content to local type
		err = doc.DataTo(&session)
		if err != nil {
			log.Printf("doc.DataTo error: %v", err)
			// Don't return early so the user still gets a response.
		}

		// Add 1 to current views value
		session.Views++

		// Update document
		_, err = a.client.Collection(collectionName).Doc(cookie.Value).Set(ctx, session)
		if err != nil {
			log.Printf("client.Collection.Doc.Set error: %v", err)
			// Don't return early so the user still gets a response.
		}
	}

	templateData := map[string]interface{}{
		"views":    session.Views,
		"greeting": session.Greetings,
	}

	if err := a.tmpl.Execute(w, templateData); err != nil {
		log.Printf("Execute: %v", err)
	}
}

// [END getting_started_sessions_handler]
