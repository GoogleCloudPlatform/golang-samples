// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START example_handler]
package demo

import (
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/xmpp"
)

func init() {
	xmpp.Handle(handleChat)
}

func handleChat(ctx context.Context, m *xmpp.Message) {
	if strings.HasPrefix(m.Body, "hello") {
		reply := &xmpp.Message{
			To:   []string{m.Sender},
			Body: "hey there!",
		}
		err := reply.Send(ctx)
		if err != nil {
			log.Errorf(ctx, "Sending reply: %v", err)
		}
	}
}

// [END example_handler]
