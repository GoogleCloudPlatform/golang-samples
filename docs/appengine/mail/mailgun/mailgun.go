// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package mailgun

import (
	"fmt"
	"net/http"

	"github.com/mailgun/mailgun-go"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func SendSimpleMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	httpc := urlfetch.Client(ctx)

	mg := mailgun.NewMailgun(
		"YOUR_DOMAIN_NAME", // Domain name
		"YOUR_API_KEY",     // API Key
		"YOUR_PUBLIC_KEY",  // Public Key
	)
	mg.SetClient(httpc)

	msg, id, err := mg.Send(mg.NewMessage(
		/* From */ "Excited User <mailgun@YOUR_DOMAIN_NAME>",
		/* Subject */ "Hello",
		/* Body */ "Testing some Mailgun awesomness!",
		/* To */ "bar@example.com", "YOU@YOUR_DOMAIN_NAME",
	))
	if err != nil {
		msg := fmt.Sprintf("Could not send message: %v, ID %v, %+v", err, id, msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Message sent!"))
}

func SendComplexMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	httpc := urlfetch.Client(ctx)

	mg := mailgun.NewMailgun(
		"YOUR_DOMAIN_NAME", // Domain name
		"YOUR_API_KEY",     // API Key
		"YOUR_PUBLIC_KEY",  // Public Key
	)
	mg.SetClient(httpc)

	message := mg.NewMessage(
		/* From */ "Excited User <mailgun@YOUR_DOMAIN_NAME>",
		/* Subject */ "Hello",
		/* Body */ "Testing some Mailgun awesomness!",
		/* To */ "foo@example.com",
	)
	message.AddCC("baz@example.com")
	message.AddBCC("bar@example.com")
	message.SetHtml("<html>HTML version of the body</html>")
	message.AddAttachment("files/test.jpg")
	message.AddAttachment("files/test.txt")

	msg, id, err := mg.Send(message)
	if err != nil {
		msg := fmt.Sprintf("Could not send message: %v, ID %v, %+v", err, id, msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Message sent!"))
}
