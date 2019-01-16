// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START functions_pubsub_unit_test]

package helloworld

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"
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
		HelloPubSub(context.Background(), m)

		w.Close()
		log.SetOutput(os.Stderr)
		log.SetFlags(originalFlags)

		out, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}
		if got := string(out); got != test.want {
			t.Errorf("HelloPubSub(%q) = %q, want %q", test.data, got, test.want)
		}
	}
}

// [END functions_pubsub_unit_test]
