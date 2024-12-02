// Copyright 2021 Google LLC
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

// [START functions_cloudevent_pubsub_unit_test]

package helloworld

import (
	"context"
	"io"
	"log"
	"os"
	"testing"

	"github.com/cloudevents/sdk-go/v2/event"
)

func TestHelloPubSub(t *testing.T) {
	tests := []struct {
		data string
		want string
	}{
		{want: "Hello, World!\n"},
		{data: "Go", want: "Hello, Go!\n"},
	}
	for _, test := range tests {
		r, w, _ := os.Pipe()
		log.SetOutput(w)
		originalFlags := log.Flags()
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

		m := PubSubMessage{
			Data: []byte(test.data),
		}
		msg := MessagePublishedData{
			Message: m,
		}
		e := event.New()
		e.SetDataContentType("application/json")
		e.SetData(e.DataContentType(), msg)

		helloPubSub(context.Background(), e)

		w.Close()
		log.SetOutput(os.Stderr)
		log.SetFlags(originalFlags)

		out, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}
		if got := string(out); got != test.want {
			t.Errorf("HelloPubSub(%q) = %q, want %q", test.data, got, test.want)
		}
	}
}

// [END functions_cloudevent_pubsub_unit_test]
