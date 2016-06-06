// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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

func TestMainFunc(t *testing.T) {
	wt := webtest.New(t, "localhost:8080")
	m := testutil.BuildMain(t)
	defer m.Cleanup()
	m.Run(nil, func() {
		wt.WaitForNet()
		bodyContains(t, wt, "/", "No books found")
	})
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
	m.CreateFormFile("image", "")
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
	m.CreateFormFile("image", "")
	m.Close()

	resp, err := wt.Post(bookPath, "multipart/form-data; boundary="+m.Boundary(), &body)
	if err != nil {
		t.Fatal(err)
	}

	gotPath := resp.Request.URL.Path
	if wantPrefix := "/books/"; !strings.HasPrefix(gotPath, wantPrefix) {
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
