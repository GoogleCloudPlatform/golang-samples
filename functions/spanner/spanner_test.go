// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package spanner

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

func TestHelloSpanner(t *testing.T) {
	instance := os.Getenv("GOLANG_SAMPLES_SPANNER")
	if instance == "" {
		t.Skip("GOLANG_SAMPLES_SPANNER not set")
	}
	// TODO: use testutil
	db = fmt.Sprintf("%s/databases/functions-%s", instance, "golang-samples-tests")

	adminClient, err := database.NewDatabaseAdminClient(context.Background())
	if err != nil {
		t.Fatalf("NewDatabaseAdminClient: %v", err)
	}

	createDatabase(context.Background(), adminClient, db)

	defer func() {
		// testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		err := adminClient.DropDatabase(context.Background(), &adminpb.DropDatabaseRequest{Database: db})
		if err != nil {
			t.Errorf("DropDatabase(%q): %v", db, err)
		}
		// })
	}()

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	HelloSpanner(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("HelloSpanner got code %v, want %v", rr.Code, http.StatusOK)
	}
}

func createDatabase(ctx context.Context, adminClient *database.DatabaseAdminClient, db string) error {
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("Invalid database id %s", db)
	}
	op, err := adminClient.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          matches[1],
		CreateStatement: "CREATE DATABASE `" + matches[2] + "`",
		ExtraStatements: []string{
			`CREATE TABLE Singers (
				SingerId   INT64 NOT NULL,
				FirstName  STRING(1024),
				LastName   STRING(1024),
				SingerInfo BYTES(MAX)
			) PRIMARY KEY (SingerId)`,
			`CREATE TABLE Albums (
				SingerId     INT64 NOT NULL,
				AlbumId      INT64 NOT NULL,
				AlbumTitle   STRING(MAX)
			) PRIMARY KEY (SingerId, AlbumId),
			INTERLEAVE IN PARENT Singers ON DELETE CASCADE`,
		},
	})
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}
	return nil
}
