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
import (
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// [END import]

func main() {
	http.HandleFunc("/sendmail", sendMailHandler)

	appengine.Main()
}

var sendgridKey string

func init() {
	sendgridKey = os.Getenv("SENDGRID_API_KEY")
	if sendgridKey == "" {
		log.Fatal("SENDGRID_API_KEY environment variable not set.")
	}
}

func sendMailHandler(w http.ResponseWriter, r *http.Request) {
	to := &mail.Email{Address: "example@email.com"}
	from := &mail.Email{Address: "sendgrid@appengine.com"}
	subject := "Email from SendGrid"
	content := mail.NewContent("text/plain", "Through App Engine")

	body := mail.NewV3MailInit(from, subject, to, content)

	req := sendgrid.GetRequest(sendgridKey, "/v3/mail/send", "")
	req.Method = "POST"
	req.Body = mail.GetRequestBody(body)
	if _, err := sendgrid.API(req); err != nil {
		http.Error(w, fmt.Sprintf("could not send mail: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "email sent successfully.")
}
