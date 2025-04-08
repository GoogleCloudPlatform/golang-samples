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
