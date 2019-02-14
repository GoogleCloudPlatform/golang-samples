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

// Sample sendgrid is a demonstration on sending an e-mail from App Engine flexible environment.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/appengine"
)

// [START gae_flex_sendgrid_import]
import "gopkg.in/sendgrid/sendgrid-go.v2"

// [END gae_flex_sendgrid_import]

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

// [START gae_flex_sendgrid]
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

// [END gae_flex_sendgrid]
