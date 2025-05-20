// Copyright 2025 Google LLC
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

package main

// [START getting_started_sessions_setup]
import (
	"context"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
)

// app stores a sessions.Store. Create a new app with newApp.
type app struct {
	tmpl         *template.Template
	collectionID string
	projectID    string
}

// session stores the client's session information.
// This type is also used for executing the template.
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

	// collectionID is a non-empty identifier for this app, it is used as the Firestore
	// collection name that stores the sessions.
	//
	// Set it to something more descriptive for your app.
	collectionID := "hello-views"

	a, err := newApp(projectID, collectionID)
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
func newApp(projectID, collectionID string) (app, error) {
	tmpl, err := template.New("Index").Parse(`<body>{{.Views}} {{if eq .Views 1}}view{{else}}views{{end}} for "{{.Greetings}}"</body>`)
	if err != nil {
		log.Fatalf("template.New: %v", err)
	}

	return app{
		tmpl:         tmpl,
		collectionID: collectionID,
		projectID:    projectID,
	}, nil
}

// [END getting_started_sessions_main]

// [START getting_started_sessions_handler]

// index uses sessions to assign users a random greeting and keep track of
// views.
func (a *app) index(w http.ResponseWriter, r *http.Request) {
	// Allows requests only for the root path ("/") to prevent duplicate calls.
	if r.RequestURI != "/" {
		return
	}

	var session session
	var doc *firestore.DocumentRef

	isNewSession := false

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, a.projectID)
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	// cookieName is a non-empty identifier for this app, it is used as the key name
	// that contains the session's id value.
	//
	// Set it to something more descriptive for your app.
	cookieName := "session_id"

	// If err is different to nil, it means the cookie has not been set, so it will be created.
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		// isNewSession flag is set to true
		isNewSession = true
	}

	// If isNewSession flag is true, the session will be created
	if isNewSession {
		// Get unique id for new document
		doc = client.Collection(a.collectionID).NewDoc()

		session.Greetings = greetings[rand.Intn(len(greetings))]
		session.Views = 1

		// Cookie is set
		cookie = &http.Cookie{
			Name:  cookieName,
			Value: doc.ID,
		}
		http.SetCookie(w, cookie)
	} else {
		// The session exists

		// Retrieve document from collection by ID
		docSnapshot, err := client.Collection(a.collectionID).Doc(cookie.Value).Get(ctx)
		if err != nil {
			log.Printf("doc.Get error: %v", err)
			http.Error(w, "Error getting session", http.StatusInternalServerError)
			return
		}

		// Unmarshal documents's content to local type
		err = docSnapshot.DataTo(&session)
		if err != nil {
			log.Printf("doc.DataTo error: %v", err)
			http.Error(w, "Error parsing session", http.StatusInternalServerError)
			return
		}

		doc = docSnapshot.Ref

		// Add 1 to current views value
		session.Views++
	}

	// The document is created/updated
	_, err = doc.Set(ctx, session)
	if err != nil {
		log.Printf("doc.Set error: %v", err)
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	if err := a.tmpl.Execute(w, session); err != nil {
		log.Printf("Execute: %v", err)
	}
}

// [END getting_started_sessions_handler]
