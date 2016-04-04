// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START example_sending]
package demo

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/xmpp"
)

func init() {
	http.HandleFunc("/send", sendHandler)
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	m := &xmpp.Message{
		To:   []string{"example@gmail.com"},
		Body: "Someone has sent you a gift: http://example.com/gifts/",
	}
	err := m.Send(ctx)
	if err != nil {
		// Send an email message instead...
	}
}

// [END example_sending]
