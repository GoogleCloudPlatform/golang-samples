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

// Sample twilio demonstrates sending and receiving SMS, receiving calls via Twilio from App Engine flexible environment.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/appengine"
)

// [START gae_flex_twilio_import]
import (
	"bitbucket.org/ckvist/twilio/twiml"
	"bitbucket.org/ckvist/twilio/twirest"
)

// [END gae_flex_twilio_import]

func main() {
	http.HandleFunc("/call/receive", receiveCallHandler)
	http.HandleFunc("/sms/send", sendSMSHandler)
	http.HandleFunc("/sms/receive", receiveSMSHandler)

	appengine.Main()
}

var (
	twilioClient = twirest.NewClient(
		mustGetenv("TWILIO_ACCOUNT_SID"),
		mustGetenv("TWILIO_AUTH_TOKEN"))
	twilioNumber = mustGetenv("TWILIO_NUMBER")
)

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
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
