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
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"
)

var _ BookDatabase = &memoryDB{}

// memoryDB is a simple in-memory persistence layer for books.
type memoryDB struct {
	mu     sync.Mutex
	nextID int64            // next ID to assign to a book.
	books  map[string]*Book // maps from Book ID to Book.
}

func newMemoryDB() *memoryDB {
	return &memoryDB{
		books:  make(map[string]*Book),
		nextID: 1,
	}
}

// Close closes the database.
func (db *memoryDB) Close(context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.books = nil

	return nil
}

// GetBook retrieves a book by its ID.
func (db *memoryDB) GetBook(_ context.Context, id string) (*Book, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	book, ok := db.books[id]
	if !ok {
		return nil, fmt.Errorf("memorydb: book not found with ID %q", id)
	}
	return book, nil
}

// AddBook saves a given book, assigning it a new ID.
func (db *memoryDB) AddBook(_ context.Context, b *Book) (id string, err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	b.ID = strconv.FormatInt(db.nextID, 10)
	db.books[b.ID] = b

	db.nextID++

	return b.ID, nil
}

// DeleteBook removes a given book by its ID.
func (db *memoryDB) DeleteBook(_ context.Context, id string) error {
	if id == "" {
		return errors.New("memorydb: book with unassigned ID passed into DeleteBook")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.books[id]; !ok {
		return fmt.Errorf("memorydb: could not delete book with ID %q, does not exist", id)
	}
	delete(db.books, id)
	return nil
}

// UpdateBook updates the entry for a given book.
func (db *memoryDB) UpdateBook(_ context.Context, b *Book) error {
	if b.ID == "" {
		return errors.New("memorydb: book with unassigned ID passed into UpdateBook")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	db.books[b.ID] = b
	return nil
}

// ListBooks returns a list of books, ordered by title.
func (db *memoryDB) ListBooks(_ context.Context) ([]*Book, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	var books []*Book
	for _, b := range db.books {
		books = append(books, b)
	}

	sort.Slice(books, func(i, j int) bool {
		return books[i].Title < books[j].Title
	})
	return books, nil
}
