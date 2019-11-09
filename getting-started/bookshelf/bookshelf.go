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
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/storage"
)

// Book holds metadata about a book.
type Book struct {
	ID            string
	Title         string
	Author        string
	PublishedDate string
	ImageURL      string
	Description   string
}

// BookDatabase provides thread-safe access to a database of books.
type BookDatabase interface {
	// ListBooks returns a list of books, ordered by title.
	ListBooks(context.Context) ([]*Book, error)

	// GetBook retrieves a book by its ID.
	GetBook(ctx context.Context, id string) (*Book, error)

	// AddBook saves a given book, assigning it a new ID.
	AddBook(ctx context.Context, b *Book) (id string, err error)

	// DeleteBook removes a given book by its ID.
	DeleteBook(ctx context.Context, id string) error

	// UpdateBook updates the entry for a given book.
	UpdateBook(ctx context.Context, b *Book) error
}

// Bookshelf holds a BookDatabase and storage info.
type Bookshelf struct {
	DB BookDatabase

	StorageBucket     *storage.BucketHandle
	StorageBucketName string

	// logWriter is used for request logging and can be overridden for tests.
	//
	// See https://cloud.google.com/logging/docs/setup/go for how to use the
	// Stackdriver logging client. Output to stdout and stderr is automaticaly
	// sent to Stackdriver when running on App Engine.
	logWriter io.Writer

	errorClient *errorreporting.Client
}

// NewBookshelf creates a new Bookshelf.
func NewBookshelf(projectID string, db BookDatabase) (*Bookshelf, error) {
	ctx := context.Background()

	// This Cloud Storage bucket must exist to be able to upload book pictures.
	// You can create it and make it public by running:
	//     gsutil mb my-project_bucket
	//     gsutil defacl set public-read gs://my-project_bucket
	// replacing my-project with your project ID.
	bucketName := projectID + "_bucket"
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}

	errorClient, err := errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: "bookshelf",
		OnError: func(err error) {
			fmt.Fprintf(os.Stderr, "Could not log error: %v", err)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("errorreporting.NewClient: %v", err)
	}

	b := &Bookshelf{
		logWriter:         os.Stderr,
		errorClient:       errorClient,
		DB:                db,
		StorageBucketName: bucketName,
		StorageBucket:     storageClient.Bucket(bucketName),
	}
	return b, nil
}
