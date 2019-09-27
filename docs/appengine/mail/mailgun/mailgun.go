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

package mailgun

import (
	"fmt"
	"net/http"

	mailgun "github.com/mailgun/mailgun-go/v3"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func SendSimpleMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	httpc := urlfetch.Client(ctx)

	mg := mailgun.NewMailgun(
		"YOUR_DOMAIN_NAME", // Domain name
		"YOUR_API_KEY",     // API Key
	)
	mg.SetClient(httpc)

	msg, id, err := mg.Send(appengine.NewContext(r), mg.NewMessage(
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

	msg, id, err := mg.Send(appengine.NewContext(r), message)
	if err != nil {
		msg := fmt.Sprintf("Could not send message: %v, ID %v, %+v", err, id, msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Message sent!"))
}
