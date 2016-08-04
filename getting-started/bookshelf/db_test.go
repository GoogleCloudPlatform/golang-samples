// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package bookshelf

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"cloud.google.com/go/datastore"

	"golang.org/x/net/context"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func testDB(t *testing.T, db BookDatabase) {
	defer db.Close()

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
	if got, want := gotBook.Description, b.Description; got != want {
		t.Errorf("Update description: got %q, want %q", got, want)
	}

	if err := db.DeleteBook(id); err != nil {
		t.Error(err)
	}

	if _, err := db.GetBook(id); err == nil {
		t.Error("want non-nil err")
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

func TestMySQLDB(t *testing.T) {
	t.Parallel()

	host := os.Getenv("GOLANG_SAMPLES_MYSQL_HOST")
	port := os.Getenv("GOLANG_SAMPLES_MYSQL_PORT")

	if host == "" {
		t.Skip("GOLANG_SAMPLES_MYSQL_HOST not set.")
	}
	if port == "" {
		port = "3306"
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		t.Fatalf("Could not parse port: %v", err)
	}

	db, err := newMySQLDB(MySQLConfig{
		Username: "root",
		Host:     host,
		Port:     p,
	})
	if err != nil {
		t.Fatal(err)
	}
	testDB(t, db)
}
