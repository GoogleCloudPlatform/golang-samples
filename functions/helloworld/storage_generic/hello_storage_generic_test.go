// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
