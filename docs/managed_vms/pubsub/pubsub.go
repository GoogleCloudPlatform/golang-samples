// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"

	"google.golang.org/appengine"
	"google.golang.org/cloud"
	"google.golang.org/cloud/pubsub"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

var (
	pubsubCtx context.Context

	// Pub/Sub topic to publish messages to. Populated by the PUBSUB_TOPIC
	// environment variable.
	topic string

	// Messages received by this instance.
	messagesMu sync.Mutex
	messages   []string
)

const maxMessages = 10

func main() {
	ctx := context.Background()

	// Set this in app.yaml when running in production.
	projectID := os.Getenv("GCLOUD_PROJECT")
	topic = os.Getenv("PUBSUB_TOPIC")

	var err error
	pubsubCtx, err = pubsubContext(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", listHandler)
	http.HandleFunc("/pubsub/publish", publishHandler)
	http.HandleFunc("/pubsub/push", pushHandler)

	appengine.Main()
}

func pubsubContext(ctx context.Context, projectID string) (context.Context, error) {
	httpClient, err := google.DefaultClient(ctx, pubsub.ScopePubSub)
	if err != nil {
		return nil, err
	}
	return cloud.WithContext(ctx, projectID, httpClient), nil
}

type pushRequest struct {
	Message struct {
		Attributes map[string]string
		Data       []byte
		ID         string `json:"message_id"`
	}
	Subscription string
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	msg := &pushRequest{}
	if err := json.NewDecoder(r.Body).Decode(msg); err != nil {
		http.Error(w, fmt.Sprintf("Could not decode body: %v", err), http.StatusBadRequest)
		return
	}

	messagesMu.Lock()
	defer messagesMu.Unlock()
	// Limit to ten.
	messages = append(messages, string(msg.Message.Data))
	if len(messages) > maxMessages {
		messages = messages[len(messages)-maxMessages:]
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	messagesMu.Lock()
	defer messagesMu.Unlock()

	if err := tmpl.Execute(w, messages); err != nil {
		log.Printf("Could not execute template: %v", err)
	}
}

func publishHandler(w http.ResponseWriter, r *http.Request) {
	msg := &pubsub.Message{
		Data: []byte(r.FormValue("payload")),
	}

	if _, err := pubsub.Publish(pubsubCtx, topic, msg); err != nil {
		http.Error(w, fmt.Sprintf("Could not publish message: %v", err), 500)
		return
	}

	fmt.Fprint(w, "Message published.")
}

var tmpl = template.Must(template.New("").Parse(`<!DOCTYPE html>
<html>
  <head>
    <title>Pub/Sub</title>
  </head>
  <body>
    <div>
      <p>Last ten messages received by this instance:</p>
      <ul>
      {{ range . }}
          <li>{{ . }}</li>
      {{ end }}
      </ul>
    </div>
    <!-- [START form] -->
    <form method="post" action="/pubsub/publish">
      <textarea name="payload" placeholder="Enter message here"></textarea>
      <input type="submit">
    </form>
    <!-- [END form] -->
    <p>Note: if the application is running across multiple instances, each
      instance will have its own list of messages.</p>
  </body>
</html>`))
