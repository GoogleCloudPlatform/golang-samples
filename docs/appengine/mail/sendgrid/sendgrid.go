// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sendgrid

// [START sample]

import (
	"fmt"
	"net/http"

	"github.com/sendgrid/rest"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

// TODO: put your sendgrid key here.
var sendgridKey string

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	hc := urlfetch.Client(ctx)

	to := &mail.Email{Address: "example@email.com"}
	from := &mail.Email{Address: "sendgrid@appengine.com"}
	subject := "Email from SendGrid"
	content := mail.NewContent("text/plain", "Through App Engine")

	body := mail.NewV3MailInit(from, subject, to, content)

	req := sendgrid.GetRequest(sendgridKey, "/v3/mail/send", "")
	req.Method = "POST"
	req.Body = mail.GetRequestBody(body)

	hreq := rest.BuildRequestObject(req)
	resp, err := hc.Do(hreq)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not send mail: %v", err), http.StatusInternalServerError)
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		http.Error(w, fmt.Sprintf("could not send mail: status code %v", err), resp.StatusCode)
	}

	fmt.Fprintf(w, "email sent successfully.")
}

// [END sample]
