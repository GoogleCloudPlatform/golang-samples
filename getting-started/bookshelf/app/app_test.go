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
	"fmt"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/getting-started/bookshelf"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/GoogleCloudPlatform/golang-samples/internal/webtest"
)

var wt *webtest.W

func TestMain(m *testing.M) {
	serv := httptest.NewServer(nil)
	wt = webtest.New(nil, serv.Listener.Addr().String())
	registerHandlers()

	os.Exit(m.Run())
}

// This function verifies compilation occurs without error.
// It may not be possible to run the application without
// satisfying appengine environmental dependencies such as
// the presence of a GCE metadata server.
func TestBuildable(t *testing.T) {
	m := testutil.BuildMain(t)
	defer m.Cleanup()
	if !m.Built() {
		t.Fatal("failed to compile application.")
	}
}

func TestNoBooks(t *testing.T) {
	bodyContains(t, wt, "/", "No books found")
}

func TestBookDetail(t *testing.T) {
	const title = "book mcbook"
	id, err := bookshelf.DB.AddBook(&bookshelf.Book{
		Title: title,
	})
	if err != nil {
		t.Fatal(err)
	}

	bodyContains(t, wt, "/", title)

	bookPath := fmt.Sprintf("/books/%d", id)
	bodyContains(t, wt, bookPath, title)

	if err := bookshelf.DB.DeleteBook(id); err != nil {
		t.Fatal(err)
	}

	bodyContains(t, wt, "/", "No books found")
}

func TestEditBook(t *testing.T) {
	const title = "book mcbook"
	id, err := bookshelf.DB.AddBook(&bookshelf.Book{
		Title: title,
	})
	if err != nil {
		t.Fatal(err)
	}

	bookPath := fmt.Sprintf("/books/%d", id)
	editPath := bookPath + "/edit"
	bodyContains(t, wt, editPath, "Edit book")
	bodyContains(t, wt, editPath, title)

	var body bytes.Buffer
	m := multipart.NewWriter(&body)
	m.WriteField("title", "simpsons")
	m.WriteField("author", "homer")
	m.Close()

	resp, err := wt.Post(bookPath, "multipart/form-data; boundary="+m.Boundary(), &body)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := resp.Request.URL.Path, bookPath; got != want {
		t.Errorf("got %s, want %s", got, want)
	}

	bodyContains(t, wt, bookPath, "simpsons")
	bodyContains(t, wt, bookPath, "homer")

	if err := bookshelf.DB.DeleteBook(id); err != nil {
		t.Fatalf("got err %v, want nil", err)
	}
}

func TestAddAndDelete(t *testing.T) {
	bodyContains(t, wt, "/books/add", "Add book")

	bookPath := fmt.Sprintf("/books")

	var body bytes.Buffer
	m := multipart.NewWriter(&body)
	m.WriteField("title", "simpsons")
	m.WriteField("author", "homer")
	m.Close()

	resp, err := wt.Post(bookPath, "multipart/form-data; boundary="+m.Boundary(), &body)
	if err != nil {
		t.Fatal(err)
	}

	gotPath := resp.Request.URL.Path
	if wantPrefix := "/books"; !strings.HasPrefix(gotPath, wantPrefix) {
		t.Fatalf("redirect: got %q, want prefix %q", gotPath, wantPrefix)
	}

	bodyContains(t, wt, gotPath, "simpsons")
	bodyContains(t, wt, gotPath, "homer")

	_, err = wt.Post(gotPath+":delete", "", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func bodyContains(t *testing.T, wt *webtest.W, path, contains string) (ok bool) {
	body, _, err := wt.GetBody(path)
	if err != nil {
		t.Error(err)
		return false
	}
	if !strings.Contains(body, contains) {
		t.Errorf("want %s to contain %s", body, contains)
		return false
	}
	return true
}
