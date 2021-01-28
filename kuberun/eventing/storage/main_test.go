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

package main

import (
	"io/ioutil"
	"os"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func TestReceive(t *testing.T) {
	tests := []struct {
		subject string
		want    string
	}{
		{subject: "objects/go-test.txt", want: "Detected change in GCS bucket: objects/go-test.txt\n"},
	}
	for _, test := range tests {
		old := os.Stdout // keep backup of the real stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		defer func() {
			os.Stdout = old
		}()

		event := cloudevents.NewEvent()
		event.SetID("test-id")
		event.SetSpecVersion("1.0")
		event.SetSource("https://localhost")
		event.SetType("example.type")
		event.SetSubject(test.subject)
		event.SetData(cloudevents.ApplicationJSON, map[string]string{"hello": "world"})
		Receive(event)

		w.Close()

		out, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}
		if got := string(out); got != test.want {
			t.Errorf("Receive(%q): got %q, want %q", test.subject, got, test.want)
		}
	}
}
