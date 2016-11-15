// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package e2e contains end-to-end tests for Go programs running on Google Cloud Platform.
// See README.md for details on running the tests.
package e2e

import (
	"bytes"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/time/rate"

	"github.com/GoogleCloudPlatform/golang-samples/internal/aeintegrate"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// env:flex deployments are quite flaky when done in parallel.
// Offset each deployment by some amount of time.
var limit = rate.NewLimiter(rate.Every(15*time.Second), 1)

func TestHelloWorld(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	limit.Wait(context.Background())

	helloworld := &aeintegrate.App{
		Name:      "hw",
		Dir:       tc.Path("appengine_flexible", "helloworld"),
		ProjectID: tc.ProjectID,
	}
	defer helloworld.Cleanup()

	bodyShouldContain(t, helloworld, "/", "Hello world!")
}

func TestDatastore(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	limit.Wait(context.Background())

	datastore := &aeintegrate.App{
		Name:      "ds",
		Dir:       tc.Path("appengine_flexible", "datastore"),
		ProjectID: tc.ProjectID,
		Env: map[string]string{
			"GCLOUD_DATASET_ID": tc.ProjectID,
		},
	}
	defer datastore.Cleanup()

	bodyShouldContain(t, datastore, "/", "Successfully stored")
}

func TestStorage(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	limit.Wait(context.Background())

	storage := &aeintegrate.App{
		Name:      "st",
		Dir:       tc.Path("appengine_flexible", "storage"),
		ProjectID: tc.ProjectID,
		Env: map[string]string{
			"GCLOUD_STORAGE_BUCKET": tc.ProjectID,
		},
	}
	defer storage.Cleanup()

	if deployed := bodyShouldContain(t, storage, "/", "<form method"); !deployed {
		return
	}

	// Requests may still not be routed correctly. Wait a little while.
	time.Sleep(10 * time.Second)

	url, _ := storage.URL("/upload")
	var body bytes.Buffer
	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	const filename = "flexible-storage-e2e"
	w := multipart.NewWriter(&body)
	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	fw.Write([]byte("hello"))
	w.Close()
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read body: %v", err)
	}
	if shouldContain := "Successful! URL"; !strings.Contains(string(b), shouldContain) {
		t.Errorf("Want resp to contain %q, got %q", shouldContain, string(b))
	}
	if shouldContain := filename; !strings.Contains(string(b), shouldContain) {
		t.Errorf("Want resp to contain %q, got %q", shouldContain, string(b))
	}
}

func TestMemcache(t *testing.T) {
	t.Skip("Doesn't work on Flex.")

	tc := testutil.EndToEndTest(t)
	limit.Wait(context.Background())

	memcache := &aeintegrate.App{
		Name:      "mem",
		Dir:       tc.Path("appengine_flexible", "memcache"),
		ProjectID: tc.ProjectID,
	}
	defer memcache.Cleanup()

	bodyShouldContain(t, memcache, "/", "Count")
}

func bodyShouldContain(t *testing.T, p *aeintegrate.App, path, shouldContain string) bool {
	if p.Deployed() {
		t.Fatalf("[%s] expected non-deployed app", p.Name)
	}

	if err := p.Deploy(); err != nil {
		t.Fatalf("could not deploy %s: %v", p.Name, err)
	}

	url, _ := p.URL("")
	log.Printf("(%s) Deployed to %s", p.Name, url)

	timeout := time.Now().Add(4 * time.Minute)

	for ; ; time.Sleep(time.Second) {
		resp, err := p.Get(path)
		if err != nil {
			t.Logf("Get: %v", err)
			continue
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Logf("could not read body: %v", err)
			continue
		}
		if strings.Contains(string(b), shouldContain) {
			return true
		}
		if time.Now().After(timeout) {
			t.Errorf("wanted to contain %q, but got body: %q", shouldContain, string(b))
			return false
		}
	}
}
