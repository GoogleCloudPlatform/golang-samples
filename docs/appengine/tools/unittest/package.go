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
