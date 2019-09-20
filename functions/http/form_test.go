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

package http

import (
	"bytes"
	"fmt"
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

	want := `File uploaded: "my_file.txt" (18 bytes)`
	if got := rr.Body.String(); got != want {
		t.Errorf("UploadFile got %q, want to contain %q", got, want)
	}
}
