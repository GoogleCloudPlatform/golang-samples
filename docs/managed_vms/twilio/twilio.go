// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/appengine"
)

// [START import]
import (
	"bitbucket.org/ckvist/twilio/twiml"
	"bitbucket.org/ckvist/twilio/twirest"
)

// [END import]

func main() {
	http.HandleFunc("/call/receive", receiveCallHandler)
	http.HandleFunc("/sms/send", sendSMSHandler)
	http.HandleFunc("/sms/receive", receiveSMSHandler)

	appengine.Main()
}

var (
	twilioClient *twirest.TwilioClient
	twilioNumber string
)

func init() {
	twilioSID := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if twilioSID == "" || twilioToken == "" {
		log.Fatal("TWILIO_ACCOUNT_SID and TWILIO_AUTH_TOKEN environment variables must be set.")
	}

	twilioNumber = os.Getenv("TWILIO_NUMBER")
	if twilioNumber == "" {
		log.Fatal("TWILIO_NUMBER environment variable must be set.")
	}
	twilioClient = twirest.NewClient(twilioSID, twilioToken)
}

func receiveCallHandler(w http.ResponseWriter, r *http.Request) {
	resp := twiml.NewResponse()
	resp.Action(twiml.Say{Text: "Hello from App Engine!"})
	resp.Send(w)
}

func sendSMSHandler(w http.ResponseWriter, r *http.Request) {
	to := r.FormValue("to")
	if to == "" {
		http.Error(w, "Missing 'to' parameter.", http.StatusBadRequest)
		return
	}

	msg := twirest.SendMessage{
		Text: "Hello from App Engine!",
		From: twilioNumber,
		To:   to,
	}

	resp, err := twilioClient.Request(msg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not send SMS: %v", err), 500)
		return
	}

	fmt.Fprintf(w, "SMS sent successfully. Response:\n%#v", resp.Message)
}

func receiveSMSHandler(w http.ResponseWriter, r *http.Request) {
	sender := r.FormValue("From")
	body := r.FormValue("Body")

	resp := twiml.NewResponse()
	resp.Action(twiml.Message{
		Body: fmt.Sprintf("Hello, %s, you said: %s", sender, body),
		From: twilioNumber,
		To:   sender,
	})
	resp.Send(w)
}
