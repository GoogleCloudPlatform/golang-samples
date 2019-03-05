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

// Sample bookshelf is a fully-featured app demonstrating several Google Cloud APIs, including Datastore, Cloud SQL, Cloud Storage.
// See https://cloud.google.com/go/getting-started/tutorial-app
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"

	uuid "github.com/gofrs/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/GoogleCloudPlatform/golang-samples/getting-started/bookshelf"
)

var (
	// See template.go
	listTmpl   = parseTemplate("list.html")
	editTmpl   = parseTemplate("edit.html")
	detailTmpl = parseTemplate("detail.html")
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	registerHandlers()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func registerHandlers() {
	// Use gorilla/mux for rich routing.
	// See http://www.gorillatoolkit.org/pkg/mux
	r := mux.NewRouter()

	r.Handle("/", http.RedirectHandler("/books", http.StatusFound))

	r.Methods("GET").Path("/books").
		Handler(appHandler(listHandler))
	r.Methods("GET").Path("/books/mine").
		Handler(appHandler(listMineHandler))
	r.Methods("GET").Path("/books/{id:[0-9]+}").
		Handler(appHandler(detailHandler))
	r.Methods("GET").Path("/books/add").
		Handler(appHandler(addFormHandler))
	r.Methods("GET").Path("/books/{id:[0-9]+}/edit").
		Handler(appHandler(editFormHandler))

	r.Methods("POST").Path("/books").
		Handler(appHandler(createHandler))
	r.Methods("POST", "PUT").Path("/books/{id:[0-9]+}").
		Handler(appHandler(updateHandler))
	r.Methods("POST").Path("/books/{id:[0-9]+}:delete").
		Handler(appHandler(deleteHandler)).Name("delete")

	// The following handlers are defined in auth.go and used in the
	// "Authenticating Users" part of the Getting Started guide.
	r.Methods("GET").Path("/login").
		Handler(appHandler(loginHandler))
	r.Methods("POST").Path("/logout").
		Handler(appHandler(logoutHandler))
	r.Methods("GET").Path("/oauth2callback").
		Handler(appHandler(oauthCallbackHandler))

	// Respond to App Engine and Compute Engine health checks.
	// Indicate the server is healthy.
	r.Methods("GET").Path("/_ah/health").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})

	// [START request_logging]
	// Delegate all of the HTTP routing and serving to the gorilla/mux router.
	// Log all requests using the standard Apache format.
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
	// [END request_logging]
}

// listHandler displays a list with summaries of books in the database.
func listHandler(w http.ResponseWriter, r *http.Request) *appError {
	books, err := bookshelf.DB.ListBooks()
	if err != nil {
		return appErrorf(err, "could not list books: %v", err)
	}

	return listTmpl.Execute(w, r, books)
}

// listMineHandler displays a list of books created by the currently
// authenticated user.
func listMineHandler(w http.ResponseWriter, r *http.Request) *appError {
	user := profileFromSession(r)
	if user == nil {
		http.Redirect(w, r, "/login?redirect=/books/mine", http.StatusFound)
		return nil
	}

	books, err := bookshelf.DB.ListBooksCreatedBy(user.ID)
	if err != nil {
		return appErrorf(err, "could not list books: %v", err)
	}

	return listTmpl.Execute(w, r, books)
}

// bookFromRequest retrieves a book from the database given a book ID in the
// URL's path.
func bookFromRequest(r *http.Request) (*bookshelf.Book, error) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("bad book id: %v", err)
	}
	book, err := bookshelf.DB.GetBook(id)
	if err != nil {
		return nil, fmt.Errorf("could not find book: %v", err)
	}
	return book, nil
}

// detailHandler displays the details of a given book.
func detailHandler(w http.ResponseWriter, r *http.Request) *appError {
	book, err := bookFromRequest(r)
	if err != nil {
		return appErrorf(err, "%v", err)
	}

	return detailTmpl.Execute(w, r, book)
}

// addFormHandler displays a form that captures details of a new book to add to
// the database.
func addFormHandler(w http.ResponseWriter, r *http.Request) *appError {
	return editTmpl.Execute(w, r, nil)
}

// editFormHandler displays a form that allows the user to edit the details of
// a given book.
func editFormHandler(w http.ResponseWriter, r *http.Request) *appError {
	book, err := bookFromRequest(r)
	if err != nil {
		return appErrorf(err, "%v", err)
	}

	return editTmpl.Execute(w, r, book)
}

// bookFromForm populates the fields of a Book from form values
// (see templates/edit.html).
func bookFromForm(r *http.Request) (*bookshelf.Book, error) {
	imageURL, err := uploadFileFromForm(r)
	if err != nil {
		return nil, fmt.Errorf("could not upload file: %v", err)
	}
	if imageURL == "" {
		imageURL = r.FormValue("imageURL")
	}

	book := &bookshelf.Book{
		Title:         r.FormValue("title"),
		Author:        r.FormValue("author"),
		PublishedDate: r.FormValue("publishedDate"),
		ImageURL:      imageURL,
		Description:   r.FormValue("description"),
		CreatedBy:     r.FormValue("createdBy"),
		CreatedByID:   r.FormValue("createdByID"),
	}

	// If the form didn't carry the user information for the creator, populate it
	// from the currently logged in user (or mark as anonymous).
	if book.CreatedByID == "" {
		user := profileFromSession(r)
		if user != nil {
			// Logged in.
			book.CreatedBy = user.DisplayName
			book.CreatedByID = user.ID
		} else {
			// Not logged in.
			book.SetCreatorAnonymous()
		}
	}

	return book, nil
}

// uploadFileFromForm uploads a file if it's present in the "image" form field.
func uploadFileFromForm(r *http.Request) (url string, err error) {
	f, fh, err := r.FormFile("image")
	if err == http.ErrMissingFile {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	if bookshelf.StorageBucket == nil {
		return "", errors.New("storage bucket is missing - check config.go")
	}

	// random filename, retaining existing extension.
	name := uuid.Must(uuid.NewV4()).String() + path.Ext(fh.Filename)

	ctx := context.Background()
	w := bookshelf.StorageBucket.Object(name).NewWriter(ctx)

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
	return fmt.Sprintf(publicURL, bookshelf.StorageBucketName, name), nil
}

// createHandler adds a book to the database.
func createHandler(w http.ResponseWriter, r *http.Request) *appError {
	book, err := bookFromForm(r)
	if err != nil {
		return appErrorf(err, "could not parse book from form: %v", err)
	}
	id, err := bookshelf.DB.AddBook(book)
	if err != nil {
		return appErrorf(err, "could not save book: %v", err)
	}
	go publishUpdate(id)
	http.Redirect(w, r, fmt.Sprintf("/books/%d", id), http.StatusFound)
	return nil
}

// updateHandler updates the details of a given book.
func updateHandler(w http.ResponseWriter, r *http.Request) *appError {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		return appErrorf(err, "bad book id: %v", err)
	}

	book, err := bookFromForm(r)
	if err != nil {
		return appErrorf(err, "could not parse book from form: %v", err)
	}
	book.ID = id

	err = bookshelf.DB.UpdateBook(book)
	if err != nil {
		return appErrorf(err, "could not save book: %v", err)
	}
	go publishUpdate(book.ID)
	http.Redirect(w, r, fmt.Sprintf("/books/%d", book.ID), http.StatusFound)
	return nil
}

// deleteHandler deletes a given book.
func deleteHandler(w http.ResponseWriter, r *http.Request) *appError {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		return appErrorf(err, "bad book id: %v", err)
	}
	err = bookshelf.DB.DeleteBook(id)
	if err != nil {
		return appErrorf(err, "could not delete book: %v", err)
	}
	http.Redirect(w, r, "/books", http.StatusFound)
	return nil
}

// publishUpdate notifies Pub/Sub subscribers that the book identified with
// the given ID has been added/modified.
func publishUpdate(bookID int64) {
	if bookshelf.PubsubClient == nil {
		return
	}

	ctx := context.Background()

	b, err := json.Marshal(bookID)
	if err != nil {
		return
	}
	topic := bookshelf.PubsubClient.Topic(bookshelf.PubsubTopicID)
	_, err = topic.Publish(ctx, &pubsub.Message{Data: b}).Get(ctx)
	log.Printf("Published update to Pub/Sub for Book ID %d: %v", bookID, err)
}

// http://blog.golang.org/error-handling-and-go
type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	Error   error
	Message string
	Code    int
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		log.Printf("Handler error: status code: %d, message: %s, underlying err: %#v",
			e.Code, e.Message, e.Error)

		http.Error(w, e.Message, e.Code)
	}
}

func appErrorf(err error, format string, v ...interface{}) *appError {
	return &appError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}
