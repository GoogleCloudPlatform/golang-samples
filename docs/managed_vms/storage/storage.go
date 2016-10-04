// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample storage demonstrates use of the cloud.google.com/go/storage package from App Engine flexible environment.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"cloud.google.com/go/storage"

	"google.golang.org/appengine"

	"golang.org/x/net/context"
)

var (
	storageClient *storage.Client

	// Set this in app.yaml when running in production.
	bucket = os.Getenv("GCLOUD_STORAGE_BUCKET")
)

func main() {
	ctx := context.Background()

	var err error
	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", formHandler)
	http.HandleFunc("/upload", uploadHandler)

	appengine.Main()
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	ctx := appengine.NewContext(r)

	f, fh, err := r.FormFile("file")
	if err != nil {
		msg := fmt.Sprintf("Could not get file: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	defer f.Close()

	sw := storageClient.Bucket(bucket).Object(fh.Filename).NewWriter(ctx)
	if _, err := io.Copy(sw, f); err != nil {
		msg := fmt.Sprintf("Could not write file: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if err := sw.Close(); err != nil {
		msg := fmt.Sprintf("Could not put file: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	u, _ := url.Parse("/" + bucket + "/" + sw.Attrs().Name)

	fmt.Fprintf(w, "Successful! URL: https://storage.googleapis.com%s", u.EscapedPath())
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, formHTML)
}

const formHTML = `<!DOCTYPE html>
<html>
  <head>
    <title>Storage</title>
    <meta charset="utf-8">
  </head>
  <body>
    <form method="POST" action="/upload" enctype="multipart/form-data">
      <input type="file" name="file">
      <input type="submit">
    </form>
  </body>
</html>`
