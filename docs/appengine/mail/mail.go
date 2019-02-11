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

// [START intro_1]
import (
	"bytes"
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/mail"
)

func confirm(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
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

const confirmMessage = `
Thank you for creating an account!
Please confirm your email address by clicking on the link below:

%s
`

// [END intro_1]

func createConfirmationURL(r *http.Request) string {
	return ""
}

// [START intro_3]
func init() {
	http.HandleFunc("/_ah/mail/", incomingMail)
}

// [END intro_3]

// [START intro_4]
func incomingMail(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	defer r.Body.Close()
	var b bytes.Buffer
	if _, err := b.ReadFrom(r.Body); err != nil {
		log.Errorf(ctx, "Error reading body: %v", err)
		return
	}
	log.Infof(ctx, "Received mail: %v", b)
}

// [END intro_4]
