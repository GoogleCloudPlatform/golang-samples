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

// Package e2e contains end-to-end tests for Go programs running on Google Cloud Platform.
// See README.md for details on running the tests.
package e2e

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/aeintegrate"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestHelloWorld(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	helloworld := &aeintegrate.App{
		Name:      "hw",
		Dir:       tc.Path("appengine", "go11x", "helloworld"),
		ProjectID: tc.ProjectID,
	}
	defer helloworld.Cleanup()

	bodyShouldContain(t, helloworld, "/", "Hello, World!")
}

func TestGo11xStatic(t *testing.T) {
	tc := testutil.EndToEndTest(t)

	helloworld := &aeintegrate.App{
		Name:      "static",
		Dir:       tc.Path("appengine", "go11x", "static"),
		ProjectID: tc.ProjectID,
	}
	defer helloworld.Cleanup()

	bodyShouldContain(t, helloworld, "/", "The Gopher Network")
}

func bodyShouldContain(t *testing.T, p *aeintegrate.App, path, shouldContain string) {
	if p.Deployed() {
		t.Fatalf("[%s] expected non-deployed app", p.Name)
	}

	if err := p.Deploy(); err != nil {
		t.Fatalf("could not deploy %s: %v", p.Name, err)
	}

	url, _ := p.URL("")
	log.Printf("(%s) Deployed to %s", p.Name, url)

	testutil.Retry(t, 20, 10*time.Second, func(r *testutil.R) {
		resp, err := p.Get(path)
		if err != nil {
			r.Errorf("Get: %v", err)
			return
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			r.Errorf("could not read body: %v", err)
			return
		}
		if got := string(b); !strings.Contains(got, shouldContain) {
			r.Errorf("got\n----\n%s\n----Want to contain:\n----%s\n", got, shouldContain)
		}
	})
}
