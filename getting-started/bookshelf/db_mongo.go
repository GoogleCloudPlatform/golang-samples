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
	"crypto/rand"
	"fmt"
	"math/big"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoDB struct {
	conn *mgo.Session
	c    *mgo.Collection
}

// Ensure mongoDB conforms to the BookDatabase interface.
var _ BookDatabase = &mongoDB{}

// newMongoDB creates a new BookDatabase backed by a given Mongo server,
// authenticated with given credentials.
func newMongoDB(addr string, cred *mgo.Credential) (BookDatabase, error) {
	conn, err := mgo.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("mongo: could not dial: %v", err)
	}

	if cred != nil {
		if err := conn.Login(cred); err != nil {
			return nil, err
		}
	}

	return &mongoDB{
		conn: conn,
		c:    conn.DB("bookshelf").C("books"),
	}, nil
}

// Close closes the database.
func (db *mongoDB) Close() {
	db.conn.Close()
}

// GetBook retrieves a book by its ID.
func (db *mongoDB) GetBook(id int64) (*Book, error) {
	b := &Book{}
	if err := db.c.Find(bson.D{{Name: "id", Value: id}}).One(b); err != nil {
		return nil, err
	}
	return b, nil
}

var maxRand = big.NewInt(1<<63 - 1)

// randomID returns a positive number that fits within an int64.
func randomID() (int64, error) {
	// Get a random number within the range [0, 1<<63-1)
	n, err := rand.Int(rand.Reader, maxRand)
	if err != nil {
		return 0, err
	}
	// Don't assign 0.
	return n.Int64() + 1, nil
}

// AddBook saves a given book, assigning it a new ID.
func (db *mongoDB) AddBook(b *Book) (id int64, err error) {
	id, err = randomID()
	if err != nil {
		return 0, fmt.Errorf("mongodb: could not assign a new ID: %v", err)
	}

	b.ID = id
	if err := db.c.Insert(b); err != nil {
		return 0, fmt.Errorf("mongodb: could not add book: %v", err)
	}
	return id, nil
}

// DeleteBook removes a given book by its ID.
func (db *mongoDB) DeleteBook(id int64) error {
	return db.c.Remove(bson.D{{Name: "id", Value: id}})
}

// UpdateBook updates the entry for a given book.
func (db *mongoDB) UpdateBook(b *Book) error {
	return db.c.Update(bson.D{{Name: "id", Value: b.ID}}, b)
}

// ListBooks returns a list of books, ordered by title.
func (db *mongoDB) ListBooks() ([]*Book, error) {
	var result []*Book
	if err := db.c.Find(nil).Sort("title").All(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// ListBooksCreatedBy returns a list of books, ordered by title, filtered by
// the user who created the book entry.
func (db *mongoDB) ListBooksCreatedBy(userID string) ([]*Book, error) {
	var result []*Book
	if err := db.c.Find(bson.D{{Name: "createdbyid", Value: userID}}).Sort("title").All(&result); err != nil {
		return nil, err
	}
	return result, nil
}
