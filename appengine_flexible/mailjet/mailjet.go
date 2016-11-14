// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample mailjet is a demonstration on sending an e-mail from App Engine flexible environment.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/appengine"
)

// [START import]
import "github.com/mailjet/mailjet-apiv3-go"

// [END import]

func main() {
	http.HandleFunc("/send", sendEmail)

	appengine.Main()
}

var (
	mailjetClient = mailjet.NewMailjetClient(
		mustGetenv("MJ_APIKEY_PUBLIC"),
		mustGetenv("MJ_APIKEY_PRIVATE"),
	)

	fromEmail = mustGetenv("MJ_FROM_EMAIL")
)

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}

func sendEmail(w http.ResponseWriter, r *http.Request) {
	to := r.FormValue("to")
	if to == "" {
		http.Error(w, "Missing 'to' parameter.", http.StatusBadRequest)
		return
	}

	m := &mailjet.InfoSendMail{
		FromEmail: fromEmail,
		FromName:  "App Engine Mailjet sample",
		Subject:   "Your email flight plan!",
		TextPart:  "Dear passenger, welcome to Mailjet! May the delivery force be with you!",
		HTMLPart:  "<h3>Dear passenger, welcome to Mailjet!</h3><br/>May the delivery force be with you!",
		Recipients: []mailjet.Recipient{
			{
				Email: to,
			},
		},
	}

	resp, err := mailjetClient.SendMail(m)
	if err != nil {
		msg := fmt.Sprintf("Could not send mail: %v", err)
		http.Error(w, msg, 500)
		return
	}

	fmt.Fprintf(w, "%d email(s) sent!", len(resp.Sent))
}
