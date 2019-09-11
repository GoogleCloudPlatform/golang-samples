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

package bookshelf

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// Ensure memoryDB conforms to the BookDatabase interface.
var _ BookDatabase = &memoryDB{}

// memoryDB is a simple in-memory persistence layer for books.
type memoryDB struct {
	mu     sync.Mutex
	nextID int64           // next ID to assign to a book.
	books  map[int64]*Book // maps from Book ID to Book.
}

func newMemoryDB() *memoryDB {
	return &memoryDB{
		books:  make(map[int64]*Book),
		nextID: 1,
	}
}

// Close closes the database.
func (db *memoryDB) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.books = nil
}

// GetBook retrieves a book by its ID.
func (db *memoryDB) GetBook(id int64) (*Book, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	book, ok := db.books[id]
	if !ok {
		return nil, fmt.Errorf("memorydb: book not found with ID %d", id)
	}
	return book, nil
}

// AddBook saves a given book, assigning it a new ID.
func (db *memoryDB) AddBook(b *Book) (id int64, err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	b.ID = db.nextID
	db.books[b.ID] = b

	db.nextID++

	return b.ID, nil
}

// DeleteBook removes a given book by its ID.
func (db *memoryDB) DeleteBook(id int64) error {
	if id == 0 {
		return errors.New("memorydb: book with unassigned ID passed into deleteBook")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.books[id]; !ok {
		return fmt.Errorf("memorydb: could not delete book with ID %d, does not exist", id)
	}
	delete(db.books, id)
	return nil
}

// UpdateBook updates the entry for a given book.
func (db *memoryDB) UpdateBook(b *Book) error {
	if b.ID == 0 {
		return errors.New("memorydb: book with unassigned ID passed into updateBook")
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	db.books[b.ID] = b
	return nil
}

// booksByTitle implements sort.Interface, ordering books by Title.
// https://golang.org/pkg/sort/#example__sortWrapper
type booksByTitle []*Book

func (s booksByTitle) Less(i, j int) bool { return s[i].Title < s[j].Title }
func (s booksByTitle) Len() int           { return len(s) }
func (s booksByTitle) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// ListBooks returns a list of books, ordered by title.
func (db *memoryDB) ListBooks() ([]*Book, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	var books []*Book
	for _, b := range db.books {
		books = append(books, b)
	}

	sort.Sort(booksByTitle(books))
	return books, nil
}

// ListBooksCreatedBy returns a list of books, ordered by title, filtered by
// the user who created the book entry.
func (db *memoryDB) ListBooksCreatedBy(userID string) ([]*Book, error) {
	if userID == "" {
		return db.ListBooks()
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	var books []*Book
	for _, b := range db.books {
		if b.CreatedByID == userID {
			books = append(books, b)
		}
	}

	sort.Sort(booksByTitle(books))
	return books, nil
}
