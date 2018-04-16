// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample mailjet is a demonstration on sending an e-mail from App Engine standard environment.
package mailjet

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

// [START import]
import "github.com/mailjet/mailjet-apiv3-go"

// [END import]

func init() {
	// Check env variables are set.
	mustGetenv("MJ_APIKEY_PUBLIC")
	mustGetenv("MJ_APIKEY_PRIVATE")

	http.HandleFunc("/send", sendEmail)
}

var fromEmail = mustGetenv("MJ_FROM_EMAIL")

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}

func sendEmail(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	mailjetClient := mailjet.NewMailjetClient(
		mustGetenv("MJ_APIKEY_PUBLIC"),
		mustGetenv("MJ_APIKEY_PRIVATE"),
	)

	mailjetClient.SetClient(urlfetch.Client(ctx))

	to := r.FormValue("to")
	if to == "" {
		http.Error(w, "Missing 'to' parameter.", http.StatusBadRequest)
		return
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: fromEmail,
				Name:  "Mailjet Pilot",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: to,
					Name:  "passenger 1",
				},
			},
			Subject:  "Your email flight plan!",
			TextPart: "Dear passenger, welcome to Mailjet! May the delivery force be with you!",
			HTMLPart: "<h3>Dear passenger, welcome to Mailjet!</h3><br />May the delivery force be with you!",
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	resp, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		msg := fmt.Sprintf("Could not send mail: %v", err)
		http.Error(w, msg, 500)
		return
	}

	fmt.Fprintf(w, "%d email(s) sent!", len(resp.ResultsV31))
}
