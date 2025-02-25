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

package sendgrid

// [START gae_mail_sendgrid]

import (
	"net/http"

	"google.golang.org/appengine/v2/urlfetch"
	"gopkg.in/sendgrid/sendgrid-go.v2"
)

func handler(w http.ResponseWriter, r *http.Request) {
	sg := sendgrid.NewSendGridClient("sendgrid_user", "sendgrid_key")
	ctx := r.Context()

	// Set http.Client to use the App Engine urlfetch client
	sg.Client = urlfetch.Client(ctx)

	message := sendgrid.NewMail()
	message.AddTo("example@email.com")
	message.SetSubject("Email From SendGrid")
	message.SetHTML("Through AppEngine")
	message.SetFrom("sendgrid@appengine.com")
	sg.Send(message)
}

// [END gae_mail_sendgrid]
