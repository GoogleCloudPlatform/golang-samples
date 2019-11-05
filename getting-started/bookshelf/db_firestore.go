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
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// firestoreDB persists books to Cloud Firestore.
// See https://cloud.google.com/firestore/docs.
type firestoreDB struct {
	client *firestore.Client
}

// Ensure firestoreDB conforms to the BookDatabase interface.
var _ BookDatabase = &firestoreDB{}

// [START getting_started_bookshelf_firestore]

// newFirestoreDB creates a new BookDatabase backed by Cloud Firestore.
// See the firestore package for details on creating a suitable
// firestore.Client: https://godoc.org/cloud.google.com/go/firestore.
func newFirestoreDB(client *firestore.Client) (*firestoreDB, error) {
	ctx := context.Background()
	// Verify that we can communicate and authenticate with the Firestore
	// service.
	err := client.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("firestoredb: could not connect: %v", err)
	}
	return &firestoreDB{
		client: client,
	}, nil
}

// Close closes the database.
func (db *firestoreDB) Close(context.Context) error {
	return db.client.Close()
}

// Book retrieves a book by its ID.
func (db *firestoreDB) GetBook(ctx context.Context, id string) (*Book, error) {
	ds, err := db.client.Collection("books").Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("firestoredb: Get: %v", err)
	}
	b := &Book{}
	ds.DataTo(b)
	return b, nil
}

// [END getting_started_bookshelf_firestore]

// AddBook saves a given book, assigning it a new ID.
func (db *firestoreDB) AddBook(ctx context.Context, b *Book) (id string, err error) {
	ref := db.client.Collection("books").NewDoc()
	b.ID = ref.ID
	if _, err := ref.Create(ctx, b); err != nil {
		return "", fmt.Errorf("Create: %v", err)
	}
	return ref.ID, nil
}

// DeleteBook removes a given book by its ID.
func (db *firestoreDB) DeleteBook(ctx context.Context, id string) error {
	if _, err := db.client.Collection("books").Doc(id).Delete(ctx); err != nil {
		return fmt.Errorf("firestore: Delete: %v", err)
	}
	return nil
}

// UpdateBook updates the entry for a given book.
func (db *firestoreDB) UpdateBook(ctx context.Context, b *Book) error {
	if _, err := db.client.Collection("books").Doc(b.ID).Set(ctx, b); err != nil {
		return fmt.Errorf("firestsore: Set: %v", err)
	}
	return nil
}

// ListBooks returns a list of books, ordered by title.
func (db *firestoreDB) ListBooks(ctx context.Context) ([]*Book, error) {
	books := make([]*Book, 0)
	iter := db.client.Collection("books").Query.OrderBy("Title", firestore.Asc).Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("firestoredb: could not list books: %v", err)
		}
		b := &Book{}
		doc.DataTo(b)
		log.Printf("Book %q ID: %q", b.Title, b.ID)
		books = append(books, b)
	}

	return books, nil
}
