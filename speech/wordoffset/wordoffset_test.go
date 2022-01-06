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

package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	speech "cloud.google.com/go/speech/apiv1"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestSyncLocal(t *testing.T) {
	testutil.SystemTest(t)

	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := syncWords(client, &out, "../testdata/quit.raw"); err != nil {
		t.Fatal(err)
	}
	if got, want := out.String(), `Word: "quit" (startTime=`; !strings.Contains(got, want) {
		t.Errorf("got %q, want to contain %q", got, want)
	}
}

func TestAsyncGCS(t *testing.T) {
	testutil.SystemTest(t)

	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	if err := asyncWords(client, &out, "gs://python-docs-samples-tests/speech/audio.raw"); err != nil {
		t.Fatal(err)
	}
	if got, want := out.String(), `Word: "Brooklyn" (startTime=`; !strings.Contains(got, want) {
		t.Errorf("got %q, want to contain %q", got, want)
	}
	if got, want := out.String(), `Word: "Bridge" (startTime=`; !strings.Contains(got, want) {
		t.Errorf("got %q, want to contain %q", got, want)
	}
}
