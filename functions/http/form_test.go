// Copyright 2018 Google LLC. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http/httptest"
	"testing"
)

func TestUploadFile(t *testing.T) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	fWriter, err := w.CreateFormFile("file", "my_file.txt")
	if err != nil {
		t.Errorf("Unable to create file: %v", err)
	}
	fmt.Fprintf(fWriter, "Content of my file")
	req := httptest.NewRequest("POST", "/", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	w.Close()

	rr := httptest.NewRecorder()

	UploadFile(rr, req)

	out, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	want := `File uploaded: "my_file.txt" (18 bytes)`
	if got := string(out); got != want {
		t.Errorf("UploadFile got %q, want to contain %q", got, want)
	}
}
