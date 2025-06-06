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

package sample

// [START gae_mail_intro]
import (
	"bytes"
	"fmt"
	"net/http"

	"google.golang.org/appengine/v2/log"
	"google.golang.org/appengine/v2/mail"
)

const confirmMessage = `
Thank you for creating an account!
Please confirm your email address by clicking on the link below:

%s
`

func confirm(_ http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	addr := r.FormValue("email")
	url := createConfirmationURL(r)
	msg := &mail.Message{
		Sender:  "Example.com Support <support@example.com>",
		To:      []string{addr},
		Subject: "Confirm your registration",
		Body:    fmt.Sprintf(confirmMessage, url),
	}
	if err := mail.Send(ctx, msg); err != nil {
		log.Errorf(ctx, "Couldn't send email: %v", err)
	}
}

// [END gae_mail_intro]

func createConfirmationURL(_ *http.Request) string {
	return ""
}

// [START gae_mail_init]
func init() {
	http.HandleFunc("/_ah/mail/", incomingMail)
}

// [END gae_mail_init]

// [START gae_mail_incoming_mail]
func incomingMail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()
	var b bytes.Buffer
	if _, err := b.ReadFrom(r.Body); err != nil {
		log.Errorf(ctx, "Error reading body: %v", err)
		return
	}
	log.Infof(ctx, "Received mail: %v", b)
}

// [END gae_mail_incoming_mail]
