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
	// [START uploading_a_blob_2]
	var rootTemplate = template.Must(template.New("root").Parse(rootTemplateHTML))

	const rootTemplateHTML = `
<html><body>
<form action="{{.}}" method="POST" enctype="multipart/form-data">
Upload File: <input type="file" name="file"><br>
<input type="submit" name="submit" value="Submit">
</form></body></html>
`
	// [END uploading_a_blob_2]

	// [START uploading_a_blob_1]
	ctx := appengine.NewContext(r)
	uploadURL, err := blobstore.UploadURL(ctx, "/upload", nil)
	if err != nil {
		serveError(ctx, w, err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = rootTemplate.Execute(w, uploadURL)
	if err != nil {
		log.Errorf(ctx, "%v", err)
	}
	// [END uploading_a_blob_1]
}

func sampleHandler2(w http.ResponseWriter, r *http.Request) {
	// [START uploading_a_blob_3]
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
	// [END uploading_a_blob_3]

	// [START serving_a_blob]
	blobstore.Send(w, appengine.BlobKey(r.FormValue("blobKey")))
	// [END serving_a_blob]
}

/* Requires old package (import "appengine/blobstore")

// [START writing_files_to_the_Blobstore]
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
// [END writing_files_to_the_Blobstore]

*/
