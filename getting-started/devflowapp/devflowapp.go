// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build go1.8

// Example app with code portable betwee different execution environments.
// See README.md for details.

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/golang-samples/getting-started/devflowapp/services"
)

// Checks messages sent to a user. The identity of the user is found in the
// HTTP request parametes.
func handleCheckMessages(w http.ResponseWriter, r *http.Request) {
	uvalues, ok := r.URL.Query()["user"]
	if !ok {
		fmt.Fprintf(w, "<p>Please include a value for 'user'</p>")
	} else {
		messageService := services.GetMessageService()
		user := uvalues[0]
		messages, err := services.CheckMessages(messageService, user)
		if err != nil {
			fmt.Fprintf(w, "<p>%v</p>", err)
			return
		}
		fmt.Fprintf(w, "<p>You have %d message(s)</p>", len(messages))
	}
}

// Handle a HTTP request for the root URL
func handleDefault(w http.ResponseWriter, r *http.Request) {
	pathSend := "/send?user=Friend1&friend=Friend2&text=We+miss+you!"
	pathCheck := "/messages?user=Friend2"
	fmt.Fprintf(w, "<p>Please try one of these two options:</p>"+
		"<ol>"+
		"<li><a href=\"%s\">Send a message</a></li>"+
		"<li><a href=\"%s\">Get messages</a></li>"+
		"</ol>", pathSend, pathCheck)
}

// Handle a HTTP request to send a message to a user. The identify of the
// sender and recipient and retrieved from HTTP request parametes.
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
	message := services.Message{
		User:   user,
		Friend: friend,
		Text:   text,
		Id:     -1}
	err := services.SendUserMessage(messageService, message)
	if err != nil {
		fmt.Fprintf(w, "<p>There was an error sending the message</p>\n")
		return
	}

	fmt.Fprintf(w, "<p>Message sent</p>\n")
}

// Entry point to the application
func main() {

	http.HandleFunc("/", handleDefault)
	http.HandleFunc("/messages", handleCheckMessages)
	http.HandleFunc("/send", handleSend)
	http.HandleFunc("/_ah/health", healthCheckHandler)
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Health check for the load balancer
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}
