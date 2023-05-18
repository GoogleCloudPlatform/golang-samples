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

// [START getting_started_background_app_main]

// Command index is an HTTP app that displays all previous translations
// (stored in Firestore) and has a form to request new translations. On form
// submission, the request is sent to Pub/Sub to be processed in the background.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/getting-started/background"
)

// topicName is the Pub/Sub topic to publish requests to. The Cloud Function to
// process translation requests should be subscribed to this topic.
const topicName = "translate"

// An app holds the clients and parsed templates that can be reused between
// requests.
type app struct {
	pubsubClient    *pubsub.Client
	pubsubTopic     *pubsub.Topic
	firestoreClient *firestore.Client
	tmpl            *template.Template
}

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatalf("GOOGLE_CLOUD_PROJECT must be set")
	}

	a, err := newApp(projectID, "index")
	if err != nil {
		log.Fatalf("newApp: %v", err)
	}

	http.HandleFunc("/", a.index)
	http.HandleFunc("/request-translation", a.requestTranslation)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on localhost:%v", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// newApp creates a new app.
func newApp(projectID, templateDir string) (*app, error) {
	ctx := context.Background()

	pubsubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %w", err)
	}

	pubsubTopic := pubsubClient.Topic(topicName)

	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("firestore.NewClient: %w", err)
	}

	// Template referenced relative to the module/app root.
	tmpl, err := template.ParseFiles(filepath.Join(templateDir, "index.html"))
	if err != nil {
		return nil, fmt.Errorf("template.New: %w", err)
	}

	return &app{
		pubsubClient: pubsubClient,
		pubsubTopic:  pubsubTopic,

		firestoreClient: firestoreClient,
		tmpl:            tmpl,
	}, nil
}

// [END getting_started_background_app_main]

// [START getting_started_background_app_list]

// index lists the current translations.
func (a *app) index(w http.ResponseWriter, r *http.Request) {
	docs, err := a.firestoreClient.Collection("translations").Documents(r.Context()).GetAll()
	if err != nil {
		log.Printf("GetAll: %v", err)
		http.Error(w, fmt.Sprintf("Error getting translations: %v", err), http.StatusInternalServerError)
		return
	}

	var translations []background.Translation
	for _, d := range docs {
		t := background.Translation{}
		if err := d.DataTo(&t); err != nil {
			log.Printf("DataTo: %v", err)
			http.Error(w, "Error reading translations", http.StatusInternalServerError)
			return
		}
		translations = append(translations, t)
	}

	if err := a.tmpl.Execute(w, translations); err != nil {
		log.Printf("tmpl.Execute: %v", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

// [END getting_started_background_app_list]

// [START getting_started_background_app_request]

// requestTranslation parses the request, validates it, and sends it to Pub/Sub.
func (a *app) requestTranslation(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("ParseForm: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	v := r.PostFormValue("v")
	if v == "" {
		log.Printf("Empty value")
		http.Error(w, "Empty value", http.StatusBadRequest)
		return
	}
	acceptableLanguages := map[string]bool{
		"de": true,
		"en": true,
		"es": true,
		"fr": true,
		"ja": true,
		"sw": true,
	}
	lang := r.PostFormValue("lang")
	if !acceptableLanguages[lang] {
		log.Printf("Unsupported language: %v", lang)
		http.Error(w, fmt.Sprintf("Unsupported language: %v", lang), http.StatusBadRequest)
		return
	}

	log.Printf("Translation requested: %q -> %s", v, lang)

	t := background.Translation{
		Original: v,
		Language: lang,
	}
	msg, err := json.Marshal(t)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
		http.Error(w, "Error requesting translation", http.StatusInternalServerError)
		return
	}

	res := a.pubsubTopic.Publish(r.Context(), &pubsub.Message{Data: msg})
	if _, err := res.Get(r.Context()); err != nil {
		log.Printf("Publish.Get: %v", err)
		http.Error(w, "Error requesting translation", http.StatusInternalServerError)
		return
	}
}

// [END getting_started_background_app_request]
