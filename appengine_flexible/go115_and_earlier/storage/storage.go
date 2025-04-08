// Copyright 2023 Google LLC
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

// [START gae_flex_storage_app]

// Sample storage demonstrates use of the cloud.google.com/go/storage package from App Engine flexible environment.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/appengine"
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

// [END gae_flex_storage_app]
