// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sendgrid

// [START sample]

import (
	"net/http"

	"gopkg.in/sendgrid/sendgrid-go.v2"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func handler(w http.ResponseWriter, r *http.Request) {
	sg := sendgrid.NewSendGridClient("sendgrid_user", "sendgrid_key")
	ctx := appengine.NewContext(r)

	// Set http.Client to use the App Engine urlfetch client
	sg.Client = urlfetch.Client(ctx)

	message := sendgrid.NewMail()
	message.AddTo("example@email.com")
	message.SetSubject("Email From SendGrid")
	message.SetHTML("Through AppEngine")
	message.SetFrom("sendgrid@appengine.com")
	sg.Send(message)
}

// [END sample]
