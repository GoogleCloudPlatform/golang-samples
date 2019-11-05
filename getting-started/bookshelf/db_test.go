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
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
)

func testDB(t *testing.T, db BookDatabase) {
	t.Helper()

	ctx := context.Background()

	b := &Book{
		Author:      "testy mc testface",
		Title:       fmt.Sprintf("t-%d", time.Now().Unix()),
		Description: "desc",
	}

	id, err := db.AddBook(ctx, b)
	if err != nil {
		t.Fatal(err)
	}

	b.ID = id
	b.Description = "newdesc"
	if err := db.UpdateBook(ctx, b); err != nil {
		t.Error(err)
	}

	gotBook, err := db.GetBook(ctx, id)
	if err != nil {
		t.Error(err)
	}
	if got, want := gotBook.Description, b.Description; got != want {
		t.Errorf("Update description: got %q, want %q", got, want)
	}

	if err := db.DeleteBook(ctx, id); err != nil {
		t.Error(err)
	}

	if _, err := db.GetBook(ctx, id); err == nil {
		t.Error("want non-nil err")
	}
}

func TestMemoryDB(t *testing.T) {
	testDB(t, newMemoryDB())
}

func TestFirestoreDB(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("GOLANG_SAMPLES_FIRESTORE_PROJECT not set")
	}
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	db, err := newFirestoreDB(client)
	if err != nil {
		t.Fatalf("newFirestoreDB: %v", err)
	}

	testDB(t, db)
}
