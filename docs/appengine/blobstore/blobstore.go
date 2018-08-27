// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package blobstore_example

import (
	"net/http"
	"text/template"

	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/log"
)

func sampleHandler(w http.ResponseWriter, r *http.Request) {
	// [START gae_blobstore_upload_form]
	var rootTemplate = template.Must(template.New("root").Parse(rootTemplateHTML))

	const rootTemplateHTML = `
<html><body>
<form action="{{.}}" method="POST" enctype="multipart/form-data">
Upload File: <input type="file" name="file"><br>
<input type="submit" name="submit" value="Submit">
</form></body></html>
`
	// [END gae_blobstore_upload_form]

	// [START gae_blobstore_upload_url]
	ctx := appengine.NewContext(r)
	uploadURL, err := blobstore.UploadURL(ctx, "/upload", nil)
	if err != nil {
		serveError(ctx, w, err)
		return
	}
	// [END gae_blobstore_upload_url]

	w.Header().Set("Content-Type", "text/html")
	err = rootTemplate.Execute(w, uploadURL)
	if err != nil {
		log.Errorf(ctx, "%v", err)
	}
}

func sampleHandler2(w http.ResponseWriter, r *http.Request) {
	// [START gae_blobstore_upload_handler]
	ctx := appengine.NewContext(r)
	blobs, _, err := blobstore.ParseUpload(r)
	if err != nil {
		serveError(ctx, w, err)
		return
	}
	file := blobs["file"]
	if len(file) == 0 {
		log.Errorf(ctx, "no file uploaded")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/serve/?blobKey="+string(file[0].BlobKey), http.StatusFound)
	// [END gae_blobstore_upload_handler]

	// [START gae_blobstore_serving]
	blobstore.Send(w, appengine.BlobKey(r.FormValue("blobKey")))
	// [END gae_blobstore_serving]
}

/* Requires old package (import "appengine/blobstore")

var k appengine.BlobKey
bw, err := blobstore.Create(ctx, "application/octet-stream")
if err != nil {
	return k, err
}
_, err = bw.Write([]byte("... some data ..."))
if err != nil {
	return k, err
}
err = bw.Close()
if err != nil {
	return k, err
}
return bw.Key()

*/
