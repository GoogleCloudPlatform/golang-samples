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

// Sample mailgun is a demonstration on sending an e-mail from App Engine flexible environment.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	mailgun "github.com/mailgun/mailgun-go/v3"
	"google.golang.org/appengine"
)

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
		mustGetenv("MAILGUN_API_KEY"))
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}

// [START gae_flex_mailgun_simple_message]
func sendSimpleMessageHandler(w http.ResponseWriter, r *http.Request) {
	msg, id, err := mailgunClient.Send(r.Context(), mailgunClient.NewMessage(
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

// [END gae_flex_mailgun_simple_message]

// [START gae_flex_mailgun_complex_message]
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

	msg, id, err := mailgunClient.Send(r.Context(), message)
	if err != nil {
		msg := fmt.Sprintf("Could not send message: %v, ID %v, %+v", err, id, msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Message sent!"))
}

// [END gae_flex_mailgun_complex_message]
