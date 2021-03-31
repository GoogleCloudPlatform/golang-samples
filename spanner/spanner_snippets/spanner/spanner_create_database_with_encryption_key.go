// Copyright 2021 Google LLC
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

package spanner

// [START spanner_create_database_with_encryption_key]
import (
	"context"
	"fmt"
	"io"
	"regexp"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

func createDatabaseWithCustomerManagedEncryptionKey(w io.Writer, db, kmsKeyName string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	// kmsKeyName = `projects/<project>/locations/<location>/keyRings/<key_ring>/cryptoKeys/<kms_key_name>`
	matches := regexp.MustCompile("^(.+)/databases/(.+)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("createDatabaseWithCustomerManagedEncryptionKey: invalid database id %q", db)
	}
	instanceName := matches[1]
	databaseId := matches[2]

	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("createDatabaseWithCustomerManagedEncryptionKey.NewDatabaseAdminClient: %v", err)
	}
	defer adminClient.Close()

	// Create a database with tables using a Customer Managed Encryption Key
	req := adminpb.CreateDatabaseRequest{
		Parent:          instanceName,
		CreateStatement: "CREATE DATABASE `" + databaseId + "`",
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
		EncryptionConfig: &adminpb.EncryptionConfig{KmsKeyName: kmsKeyName},
	}
	op, err := adminClient.CreateDatabase(ctx, &req)
	if err != nil {
		return fmt.Errorf("createDatabaseWithCustomerManagedEncryptionKey.CreateDatabase: %v", err)
	}
	dbObj, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("createDatabaseWithCustomerManagedEncryptionKey.Wait: %v", err)
	}
	fmt.Fprintf(w, "Created database [%s] using encryption key %q\n", dbObj.Name, dbObj.EncryptionConfig.KmsKeyName)
	return nil
}

// [END spanner_create_database_with_encryption_key]
