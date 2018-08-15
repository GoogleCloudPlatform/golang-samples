// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample sendgrid is a demonstration on sending an e-mail from App Engine flexible environment.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/appengine"
)

// [START import]
import "gopkg.in/sendgrid/sendgrid-go.v2"

// [END import]

func main() {
	http.HandleFunc("/sendmail", sendMailHandler)

	appengine.Main()
}

var sendgridClient *sendgrid.SGClient

func init() {
	sendgridKey := os.Getenv("SENDGRID_API_KEY")
	if sendgridKey == "" {
		log.Fatal("SENDGRID_API_KEY environment variable not set.")
	}
	sendgridClient = sendgrid.NewSendGridClientWithApiKey(sendgridKey)
}

func sendMailHandler(w http.ResponseWriter, r *http.Request) {
	m := sendgrid.NewMail()
	m.AddTo("example@email.com")
	m.SetSubject("Email From SendGrid")
	m.SetHTML("Through AppEngine")
	m.SetFrom("sendgrid@appengine.com")

	if err := sendgridClient.Send(m); err != nil {
		http.Error(w, fmt.Sprintf("could not send mail: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "email sent successfully.")
}
