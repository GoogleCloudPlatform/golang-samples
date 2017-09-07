// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/datastore"
)

var (
	datastoreClient *datastore.Client

	payloadKey = datastore.NameKey("CloudTaskPayload", "Singleton", nil)
)

type Payload struct {
	Body string
}

func main() {
	ctx := context.Background()
	var err error

	projectID := os.Getenv("GCLOUD_PROJECT")
	if projectID == "" {
		projectID, err = metadata.ProjectID()
		if err != nil {
			log.Fatalf("Could not get project ID: %v", err)
		}
	}

	datastoreClient, err = datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Could not initialize datastore client: %v", err)
	}

	http.HandleFunc("/", rootHandle)
	http.HandleFunc("/payload", payloadHandle)
	http.HandleFunc("/_ah/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) })
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var rootTmpl = template.Must(template.New("").Parse(`<!doctype html>
{{if .}}Last message: {{.Body}}{{else}}No messages yet.{{end}}`))

func rootHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var last Payload
	err := datastoreClient.Get(ctx, payloadKey, &last)
	if err == datastore.ErrNoSuchEntity {
		rootTmpl.Execute(w, nil)
		return
	}
	if err != nil {
		handleError(w, r, err)
		return
	}
	rootTmpl.Execute(w, last)
}

func payloadHandle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, r, err)
		return
	}

	payload := &Payload{
		Body: string(body),
	}
	_, err = datastoreClient.Put(ctx, payloadKey, payload)
	if err != nil {
		handleError(w, r, err)
		return
	}

	fmt.Fprint(w, "OK")
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "text/plain")
	http.Error(w, fmt.Sprintf("Unexpected error: %v", err), 500)

	log.Printf("Error (%s): %v", r.URL.Path, err)
	return
}
