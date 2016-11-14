// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample mailgun is a demonstration on sending an e-mail from App Engine flexible environment.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"google.golang.org/appengine"
)

// [START import]
import "github.com/mailgun/mailgun-go"

// [END import]

func main() {
	http.HandleFunc("/send_simple", sendSimpleMessageHandler)
	http.HandleFunc("/send_complex", sendComplexMessageHandler)

	appengine.Main()
}

var (
	mailgunClient mailgun.Mailgun
	mailgunDomain string
)

func init() {
	mailgunDomain = mustGetenv("MAILGUN_DOMAIN_NAME")
	mailgunClient = mailgun.NewMailgun(
		mailgunDomain,
		mustGetenv("MAILGUN_API_KEY"),
		mustGetenv("MAILGUN_PUBLIC_KEY"))
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}

func sendSimpleMessageHandler(w http.ResponseWriter, r *http.Request) {
	msg, id, err := mailgunClient.Send(mailgunClient.NewMessage(
		/* From */ fmt.Sprintf("Excited User <mailgun@%s>", mailgunDomain),
		/* Subject */ "Hello",
		/* Body */ "Testing some Mailgun awesomness!",
		/* To */ "bar@example.com", "YOU@"+mailgunDomain,
	))
	if err != nil {
		msg := fmt.Sprintf("Could not send message: %v, ID %v, %+v", err, id, msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Message sent!"))
}

func sendComplexMessageHandler(w http.ResponseWriter, r *http.Request) {
	message := mailgunClient.NewMessage(
		/* From */ fmt.Sprintf("Excited User <mailgun@%s>", mailgunDomain),
		/* Subject */ "Hello",
		/* Body */ "Testing some Mailgun awesomness!",
		/* To */ "foo@example.com",
	)
	message.AddCC("baz@example.com")
	message.AddBCC("bar@example.com")
	message.SetHtml("<html>HTML version of the body</html>")
	message.AddReaderAttachment("files/test.txt",
		ioutil.NopCloser(strings.NewReader("foo")))

	msg, id, err := mailgunClient.Send(message)
	if err != nil {
		msg := fmt.Sprintf("Could not send message: %v, ID %v, %+v", err, id, msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Message sent!"))
}
