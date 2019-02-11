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

package helloworld

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/functions/metadata"
)

func TestHelloGCSInfo(t *testing.T) {
	r, w, _ := os.Pipe()
	log.SetOutput(w)
	originalFlags := log.Flags()
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	name := "hello_gcs.txt"
	e := GCSEvent{
		Name: name,
	}
	meta := &metadata.Metadata{
		EventID: "event ID",
	}
	ctx := metadata.NewContext(context.Background(), meta)

	HelloGCSInfo(ctx, e)

	w.Close()
	log.SetOutput(os.Stderr)
	log.SetFlags(originalFlags)

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}

	got := string(out)
	wants := []string{
		"File: " + name,
		"Event ID: " + meta.EventID,
	}
	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Errorf("HelloGCSInfo(%v) = %q, want to contain %q", e, got, want)
		}
	}
}
