// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START package_example_1]
package newsletter

import (
	"reflect"
	"testing"

	"google.golang.org/appengine/mail"
)

func TestComposeNewsletter(t *testing.T) {
	want := &mail.Message{
		Sender:  "newsletter@appspot.com",
		To:      []string{"User <user@example.com>"},
		Subject: "Weekly App Engine Update",
		Body:    "Don't forget to test your code!",
	}
	if msg := composeNewsletter(); !reflect.DeepEqual(msg, want) {
		t.Errorf("composeMessage() = %+v, want %+v", msg, want)
	}
}

// [END package_example_1]

func composeNewsletter() *mail.Message {
	return &mail.Message{
		Sender:  "newsletter@appspot.com",
		To:      []string{"User <user@example.com>"},
		Subject: "Weekly App Engine Update",
		Body:    "Don't forget to test your code!",
	}
}
