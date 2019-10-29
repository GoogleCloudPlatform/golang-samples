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
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/golang-samples/getting-started/bookshelf/internal/webtest"
)

var (
	wt *webtest.W
	b  *Bookshelf

	testDBs = map[string]BookDatabase{}
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if projectID == "" {
		log.Println("GOLANG_SAMPLES_PROJECT_ID not set. Skipping.")
		return
	}

	memoryDB := newMemoryDB()
	testDBs["memory"] = memoryDB

	if firestoreProjectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT"); firestoreProjectID != "" {
		projectID = firestoreProjectID

		client, err := firestore.NewClient(ctx, projectID)
		if err != nil {
			log.Fatalf("firestore.NewClient: %v", err)
		}

		// Delete all docs first to start with a clean slate.
		docs, err := client.Collection("books").DocumentRefs(ctx).GetAll()
		if err == nil {
			for _, d := range docs {
				if _, err := d.Delete(ctx); err != nil {
					log.Fatalf("Delete: %v", err)
				}
			}
		}

		db, err := newFirestoreDB(client)
		if err != nil {
			log.Fatalf("newFirestoreDB: %v", err)
		}
		testDBs["firestore"] = db
	} else {
		log.Println("GOLANG_SAMPES_FIRESTORE_PROJECT not set. Skipping Firestore database tests.")
	}

	var err error
	b, err = NewBookshelf(projectID, memoryDB)
	if err != nil {
		log.Fatalf("NewBookshelf: %v", err)
	}

	// Don't log anything during testing.
	log.SetOutput(ioutil.Discard)
	b.logWriter = ioutil.Discard

	serv := httptest.NewServer(nil)
	wt = webtest.New(nil, serv.Listener.Addr().String())

	b.registerHandlers()

	os.Exit(m.Run())
}

func TestNoBooks(t *testing.T) {
	for name, db := range testDBs {
		t.Run(name, func(t *testing.T) {
			b.DB = db
			bodyContains(t, wt, "/", "No books found")
		})
	}
}

func TestBookDetail(t *testing.T) {
	for name, db := range testDBs {
		t.Run(name, func(t *testing.T) {
			b.DB = db
			ctx := context.Background()
			const title = "book mcbook"
			id, err := b.DB.AddBook(ctx, &Book{
				Title: title,
			})
			if err != nil {
				t.Fatal(err)
			}

			bodyContains(t, wt, "/", title)

			bookPath := fmt.Sprintf("/books/%s", id)
			bodyContains(t, wt, bookPath, title)

			if err := b.DB.DeleteBook(ctx, id); err != nil {
				t.Fatal(err)
			}

			bodyContains(t, wt, "/", "No books found")
		})
	}

}

func TestEditBook(t *testing.T) {
	for name, db := range testDBs {
		t.Run(name, func(t *testing.T) {
			b.DB = db
			ctx := context.Background()
			const title = "book mcbook"
			id, err := b.DB.AddBook(ctx, &Book{
				Title: title,
			})
			if err != nil {
				t.Fatal(err)
			}

			bookPath := fmt.Sprintf("/books/%s", id)
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

			if err := b.DB.DeleteBook(ctx, id); err != nil {
				t.Fatalf("got err %v, want nil", err)
			}
		})
	}
}

func TestAddAndDelete(t *testing.T) {
	for name, db := range testDBs {
		t.Run(name, func(t *testing.T) {
			b.DB = db
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
		})
	}
}

func TestSendLog(t *testing.T) {
	buf := &bytes.Buffer{}
	oldLogger := b.logWriter
	b.logWriter = buf

	bodyContains(t, wt, "/logs", "Log sent!")

	b.logWriter = oldLogger

	if got, want := buf.String(), "Good job!"; !strings.Contains(got, want) {
		t.Errorf("/logs logged\n----\n%v\n----\nWant to contain:\n----\n%v", got, want)
	}
}

func TestSendError(t *testing.T) {
	buf := &bytes.Buffer{}
	oldLogger := b.logWriter
	b.logWriter = buf

	bodyContains(t, wt, "/errors", "Error Reporting")

	b.logWriter = oldLogger

	if got, want := buf.String(), "uh oh"; !strings.Contains(got, want) {
		t.Errorf("/errors logged\n----\n%v\n----\nWant to contain:\n----\n%v", got, want)
	}
}

func bodyContains(t *testing.T, wt *webtest.W, path, contains string) (ok bool) {
	t.Helper()

	body, _, err := wt.GetBody(path)
	if err != nil {
		t.Error(err)
		return false
	}
	if !strings.Contains(body, contains) {
		t.Errorf("got:\n----\n%s\nWant to contain:\n----\n%s", body, contains)
		return false
	}
	return true
}
