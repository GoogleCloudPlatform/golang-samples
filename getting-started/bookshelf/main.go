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

// The bookshelf command starts the bookshelf server, a sample app
// demonstrating several Google Cloud APIs, including App Engine, Firestore, and
// Cloud Storage.
// See https://cloud.google.com/go/getting-started/tutorial-app.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/gofrs/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	// See template.go.
	listTmpl   = parseTemplate("list.html")
	editTmpl   = parseTemplate("edit.html")
	detailTmpl = parseTemplate("detail.html")
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal("GOOGLE_CLOUD_PROJECT must be set")
	}

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	db, err := newFirestoreDB(client)
	if err != nil {
		log.Fatalf("newFirestoreDB: %v", err)
	}
	b, err := NewBookshelf(projectID, db)
	if err != nil {
		log.Fatalf("NewBookshelf: %v", err)
	}

	b.registerHandlers()

	log.Printf("Listening on localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func (b *Bookshelf) registerHandlers() {
	// Use gorilla/mux for rich routing.
	// See https://www.gorillatoolkit.org/pkg/mux.
	r := mux.NewRouter()

	r.Handle("/", http.RedirectHandler("/books", http.StatusFound))

	r.Methods("GET").Path("/books").
		Handler(appHandler(b.listHandler))
	r.Methods("GET").Path("/books/add").
		Handler(appHandler(b.addFormHandler))
	r.Methods("GET").Path("/books/{id:[0-9a-zA-Z_\\-]+}").
		Handler(appHandler(b.detailHandler))
	r.Methods("GET").Path("/books/{id:[0-9a-zA-Z_\\-]+}/edit").
		Handler(appHandler(b.editFormHandler))

	r.Methods("POST").Path("/books").
		Handler(appHandler(b.createHandler))
	r.Methods("POST", "PUT").Path("/books/{id:[0-9a-zA-Z_\\-]+}").
		Handler(appHandler(b.updateHandler))
	r.Methods("POST").Path("/books/{id:[0-9a-zA-Z_\\-]+}:delete").
		Handler(appHandler(b.deleteHandler)).Name("delete")

	// Respond to App Engine and Compute Engine health checks.
	// Indicate the server is healthy.
	r.Methods("GET").Path("/_ah/health").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})

	r.Methods("GET").Path("/logs").Handler(appHandler(b.sendLog))
	r.Methods("GET").Path("/errors").Handler(appHandler(b.sendError))

	// Delegate all of the HTTP routing and serving to the gorilla/mux router.
	// Log all requests using the standard Apache format.
	http.Handle("/", handlers.CombinedLoggingHandler(b.logWriter, r))
}

// listHandler displays a list with summaries of books in the database.
func (b *Bookshelf) listHandler(w http.ResponseWriter, r *http.Request) *appError {
	ctx := r.Context()
	books, err := b.DB.ListBooks(ctx)
	if err != nil {
		return b.appErrorf(r, err, "could not list books: %v", err)
	}

	return listTmpl.Execute(b, w, r, books)
}

// bookFromRequest retrieves a book from the database given a book ID in the
// URL's path.
func (b *Bookshelf) bookFromRequest(r *http.Request) (*Book, error) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	if id == "" {
		return nil, errors.New("no book with empty ID")
	}
	book, err := b.DB.GetBook(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not find book: %v", err)
	}
	return book, nil
}

// detailHandler displays the details of a given book.
func (b *Bookshelf) detailHandler(w http.ResponseWriter, r *http.Request) *appError {
	book, err := b.bookFromRequest(r)
	if err != nil {
		return b.appErrorf(r, err, "%v", err)
	}

	return detailTmpl.Execute(b, w, r, book)
}

// addFormHandler displays a form that captures details of a new book to add to
// the database.
func (b *Bookshelf) addFormHandler(w http.ResponseWriter, r *http.Request) *appError {
	return editTmpl.Execute(b, w, r, nil)
}

// editFormHandler displays a form that allows the user to edit the details of
// a given book.
func (b *Bookshelf) editFormHandler(w http.ResponseWriter, r *http.Request) *appError {
	book, err := b.bookFromRequest(r)
	if err != nil {
		return b.appErrorf(r, err, "%v", err)
	}

	return editTmpl.Execute(b, w, r, book)
}

// bookFromForm populates the fields of a Book from form values
// (see templates/edit.html).
func (b *Bookshelf) bookFromForm(r *http.Request) (*Book, error) {
	ctx := r.Context()
	imageURL, err := b.uploadFileFromForm(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("could not upload file: %v", err)
	}
	if imageURL == "" {
		imageURL = r.FormValue("imageURL")
	}

	book := &Book{
		Title:         r.FormValue("title"),
		Author:        r.FormValue("author"),
		PublishedDate: r.FormValue("publishedDate"),
		ImageURL:      imageURL,
		Description:   r.FormValue("description"),
	}

	return book, nil
}

// [START getting_started_bookshelf_storage]

// uploadFileFromForm uploads a file if it's present in the "image" form field.
func (b *Bookshelf) uploadFileFromForm(ctx context.Context, r *http.Request) (url string, err error) {
	f, fh, err := r.FormFile("image")
	if err == http.ErrMissingFile {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	if b.StorageBucket == nil {
		return "", errors.New("storage bucket is missing: check bookshelf.go")
	}
	if _, err := b.StorageBucket.Attrs(ctx); err != nil {
		if err == storage.ErrBucketNotExist {
			return "", fmt.Errorf("bucket %q does not exist: check bookshelf.go", b.StorageBucketName)
		}
		return "", fmt.Errorf("could not get bucket: %v", err)
	}

	// random filename, retaining existing extension.
	name := uuid.Must(uuid.NewV4()).String() + path.Ext(fh.Filename)

	w := b.StorageBucket.Object(name).NewWriter(ctx)

	// Warning: storage.AllUsers gives public read access to anyone.
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.ContentType = fh.Header.Get("Content-Type")

	// Entries are immutable, be aggressive about caching (1 day).
	w.CacheControl = "public, max-age=86400"

	if _, err := io.Copy(w, f); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	const publicURL = "https://storage.googleapis.com/%s/%s"
	return fmt.Sprintf(publicURL, b.StorageBucketName, name), nil
}

// [END getting_started_bookshelf_storage]

// createHandler adds a book to the database.
func (b *Bookshelf) createHandler(w http.ResponseWriter, r *http.Request) *appError {
	ctx := r.Context()
	book, err := b.bookFromForm(r)
	if err != nil {
		return b.appErrorf(r, err, "could not parse book from form: %v", err)
	}
	id, err := b.DB.AddBook(ctx, book)
	if err != nil {
		return b.appErrorf(r, err, "could not save book: %v", err)
	}
	http.Redirect(w, r, fmt.Sprintf("/books/%s", id), http.StatusFound)
	return nil
}

// updateHandler updates the details of a given book.
func (b *Bookshelf) updateHandler(w http.ResponseWriter, r *http.Request) *appError {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	if id == "" {
		return b.appErrorf(r, errors.New("no book with empty ID"), "no book with empty ID")
	}
	book, err := b.bookFromForm(r)
	if err != nil {
		return b.appErrorf(r, err, "could not parse book from form: %v", err)
	}
	book.ID = id

	if err := b.DB.UpdateBook(ctx, book); err != nil {
		return b.appErrorf(r, err, "UpdateBook: %v", err)
	}
	http.Redirect(w, r, fmt.Sprintf("/books/%s", book.ID), http.StatusFound)
	return nil
}

// deleteHandler deletes a given book.
func (b *Bookshelf) deleteHandler(w http.ResponseWriter, r *http.Request) *appError {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	if err := b.DB.DeleteBook(ctx, id); err != nil {
		return b.appErrorf(r, err, "DeleteBook: %v", err)
	}
	http.Redirect(w, r, "/books", http.StatusFound)
	return nil
}

// sendLog logs a message.
//
// See https://cloud.google.com/logging/docs/setup/go for how to use the
// Stackdriver logging client. Output to stdout and stderr is automaticaly
// sent to Stackdriver when running on App Engine.
func (b *Bookshelf) sendLog(w http.ResponseWriter, r *http.Request) *appError {
	fmt.Fprintln(b.logWriter, "Hey, you triggered a custom log entry. Good job!")

	fmt.Fprintln(w, `<html>Log sent! Check the <a href="http://console.cloud.google.com/logs">logging section of the Cloud Console</a>.</html>`)

	return nil
}

// sendError triggers an error that is sent to Error Reporting.
func (b *Bookshelf) sendError(w http.ResponseWriter, r *http.Request) *appError {
	msg := `<html>Logging an error. Check <a href="http://console.cloud.google.com/errors">Error Reporting</a> (it may take a minute or two for the error to appear).</html>`
	err := errors.New("uh oh! an error occurred")
	return b.appErrorf(r, err, msg)
}

// https://blog.golang.org/error-handling-and-go
type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	err     error
	message string
	code    int
	req     *http.Request
	b       *Bookshelf
	stack   []byte
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		fmt.Fprintf(e.b.logWriter, "Handler error (reported to Error Reporting): status code: %d, message: %s, underlying err: %+v\n", e.code, e.message, e.err)
		w.WriteHeader(e.code)
		fmt.Fprint(w, e.message)

		e.b.errorClient.Report(errorreporting.Entry{
			Error: e.err,
			Req:   r,
			Stack: e.stack,
		})
		e.b.errorClient.Flush()
	}
}

func (b *Bookshelf) appErrorf(r *http.Request, err error, format string, v ...interface{}) *appError {
	return &appError{
		err:     err,
		message: fmt.Sprintf(format, v...),
		code:    500,
		req:     r,
		b:       b,
		stack:   debug.Stack(),
	}
}
