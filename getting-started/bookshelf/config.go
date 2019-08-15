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
	"io"
	"os"

	"cloud.google.com/go/firestore"
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

	// Close closes the database, freeing up any available resources.
	Close(ctx context.Context) error
}

// Bookshelf holds a BookDatabase and storage info.
type Bookshelf struct {
	DB BookDatabase

	StorageBucket     *storage.BucketHandle
	StorageBucketName string

	// logWriter is used for request logging and can be overridden for tests.
	logWriter io.Writer
}

// NewBookshelf creates a new Bookshelf.
func NewBookshelf(projectID string) (*Bookshelf, error) {
	b := &Bookshelf{
		logWriter: os.Stderr,
	}

	var err error
	b.DB, err = configureFirestoreDB(projectID)
	if err != nil {
		return nil, err
	}

	// This Cloud Storage bucket must exist to be able to upload book pictures.
	// You can create it and make it public by running:
	//     gsutil mb my-project_bucket
	//     gsutil defacl set public-read gs://my-project_bucket
	// replacing my-project with your project ID.
	b.StorageBucketName = projectID + "_bucket"
	b.StorageBucket, err = configureStorage(b.StorageBucketName)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func configureFirestoreDB(projectID string) (*firestoreDB, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return newFirestoreDB(client)
}

func configureStorage(bucketID string) (*storage.BucketHandle, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.Bucket(bucketID), nil
}
