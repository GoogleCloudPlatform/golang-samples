// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package bookshelf

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/cloud/datastore"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func testDB(t *testing.T, db BookDatabase) {
	b := &Book{
		Author:      "testy mc testface",
		Title:       fmt.Sprintf("t-%d", time.Now().Unix()),
		Description: "desc",
	}

	id, err := db.AddBook(b)
	if err != nil {
		t.Fatal(err)
	}

	b.ID = id
	b.Description = "newdesc"
	if err := db.UpdateBook(b); err != nil {
		t.Error(err)
	}

	gotBook, err := db.GetBook(id)
	if err != nil {
		t.Error(err)
	}
	if want, got := b.Description, gotBook.Description; want != got {
		t.Errorf("Update description: want %q, got %q", want, got)
	}

	if err := db.DeleteBook(id); err != nil {
		t.Error(err)
	}

	if _, err := db.GetBook(id); err == nil {
		t.Error("want err, got nil")
	}
}

func TestMemoryDB(t *testing.T) {
	testDB(t, newMemoryDB())
}

func TestDatastoreDB(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	db, err := newDatastoreDB(client)
	if err != nil {
		t.Fatal(err)
	}
	testDB(t, db)
}
