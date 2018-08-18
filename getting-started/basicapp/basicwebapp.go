// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

// Example app with code portable betwee different execution environments

import (
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/profiler"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"

	"github.com/GoogleCloudPlatform/golang-samples/getting-started/basicapp/services"
)

func handleCheckMessages(w http.ResponseWriter, r *http.Request) {
	uvalues, ok := r.URL.Query()["user"]
	if !ok {
		fmt.Fprintf(w, "<p>Please include a value for 'user'</p>")
	} else {
		messageService := services.GetMessageService()
		user := uvalues[0]
		ctx := r.Context()
		_, span := trace.StartSpan(ctx, "CheckMessages")
		defer span.End()
		messages, err := services.CheckMessages(messageService, user)
		if err != nil {
			fmt.Fprintf(w, "<p>%v</p>", err)
			return
		}
		fmt.Fprintf(w, "<p>You have %d message(s)</p>", len(messages))
	}
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
	pathSend := "/send?user=Friend1&friend=Friend2&text=We+miss+you!"
	pathCheck := "/messages?user=Friend2"
	fmt.Fprintf(w, "<p>Please try one of these two options:</p>"+
		"<ol>"+
		"<li><a href=\"%s\">Send a message</a></li>"+
		"<li><a href=\"%s\">Get messages</a></li>"+
		"</ol>", pathSend, pathCheck)
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	messageService := services.GetMessageService()
	uvalues, ok := r.URL.Query()["user"]
	user := "Friend1"
	if ok {
		user = uvalues[0]
	}
	fvalues, ok := r.URL.Query()["friend"]
	friend := "Friend2"
	if ok {
		friend = fvalues[0]
	}
	tvalues, ok := r.URL.Query()["text"]
	text := "We miss you!"
	if ok {
		text = tvalues[0]
	}
	message := services.Message{user, friend, text, -1}
	ctx := r.Context()
	_, span := trace.StartSpan(ctx, "SendUserMessage")
	defer span.End()
	err := services.SendUserMessage(messageService, message)
	if err != nil {
		fmt.Fprintf(w, "<p>There was an error sending the message</p>\n")
		return
	}

	fmt.Fprintf(w, "<p>Message sent</p>\n")
}

func main() {
	// Profiler initialization
	if err := profiler.Start(profiler.Config{
		Service:        "basicapp",
		ServiceVersion: "0.001",
	}); err != nil {
		log.Printf("main: running without the profiler enabled\n")
	}
	ocSetup()

	http.HandleFunc("/", handleDefault)
	http.HandleFunc("/messages", handleCheckMessages)
	http.HandleFunc("/send", handleSend)
	http.HandleFunc("/_ah/health", healthCheckHandler)
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", &ochttp.Handler{}))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

// Initialization for OpenCensus for tracing
func ocSetup() {

	log.Printf("Configure OpenCensus")

	// Create and register exporters for monitoring metrics and trace
	se, err := stackdriver.NewExporter(stackdriver.Options{})
	if err != nil {
		log.Printf("Error creating Stackdriver exporter: %v", err)
		return
	}
	trace.RegisterExporter(se)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	log.Printf("ocSetup: exporter registerd")
}
