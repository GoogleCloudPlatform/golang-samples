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
)

func TestHelloRemoteConfig(t *testing.T) {
	tests := []struct {
		updateType string
		origin     string
		version    string
	}{
		{
			updateType: "my type",
			origin:     "my origin",
			version:    "1.0.0",
		},
		{
			updateType: "my type",
			origin:     "my origin",
			version:    "2.0.0",
		},
	}

	for _, test := range tests {
		r, w, _ := os.Pipe()
		log.SetOutput(w)
		originalFlags := log.Flags()
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

		e := RemoteConfigEvent{
			UpdateType:    test.updateType,
			UpdateOrigin:  test.origin,
			VersionNumber: test.version,
		}
		HelloRemoteConfig(context.Background(), e)

		w.Close()
		log.SetOutput(os.Stderr)
		log.SetFlags(originalFlags)

		out, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}

		got := string(out)
		if !strings.Contains(got, test.updateType) {
			t.Errorf("HelloRemoteConfig(%v) got %q, want to contain UpdateType %q", e, got, test.updateType)
		}
		if !strings.Contains(got, test.origin) {
			t.Errorf("HelloRemoteConfig(%v) got %q, want to contain origin %q", e, got, test.origin)
		}
		if !strings.Contains(got, test.version) {
			t.Errorf("HelloRemoteConfig(%v) got %q, want to contain version %q", e, got, test.version)
		}
	}
}
