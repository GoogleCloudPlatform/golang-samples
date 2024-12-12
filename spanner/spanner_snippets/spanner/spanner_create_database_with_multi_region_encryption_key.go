// Copyright 2024 Google LLC
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

// [START spanner_create_database_with_MR_CMEK]

import (
	"context"
	"fmt"
	"io"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
)

// createDatabaseWithCustomerManagedMultiRegionEncryptionKey creates a new database with tables using a Customer Managed Multi-Region Encryption Key.
func createDatabaseWithCustomerManagedMultiRegionEncryptionKey(ctx context.Context, w io.Writer, projectID, instanceID, databaseID string, kmsKeyNames []string) error {
	// projectID = `my-project`
	// instanceID = `my-instance`
	// databaseID = `my-database`
	// kmsKeyNames := []string{"projects/my-project/locations/locations/<location1>/keyRings/<keyRing>/cryptoKeys/<keyId>",
	//	 "projects/my-project/locations/locations/<location2>/keyRings/<keyRing>/cryptoKeys/<keyId>",
	//	 "projects/my-project/locations/locations/<location3>/keyRings/<keyRing>/cryptoKeys/<keyId>",
	// }

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return fmt.Errorf("createDatabaseWithCustomerManagedMultiRegionEncryptionKey.NewDatabaseAdminClient: %w", err)
	}
	defer adminClient.Close()

	// Create a database with tables using a Customer Managed Multi-Region Encryption Key
	req := adminpb.CreateDatabaseRequest{
		Parent:          fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID),
		CreateStatement: "CREATE DATABASE `" + databaseID + "`",
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
		EncryptionConfig: &adminpb.EncryptionConfig{KmsKeyNames: kmsKeyNames},
	}
	op, err := adminClient.CreateDatabase(ctx, &req)
	if err != nil {
		return fmt.Errorf("createDatabaseWithCustomerManagedMultiRegionEncryptionKey.CreateDatabase: %w", err)
	}
	dbObj, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("createDatabaseWithCustomerManagedMultiRegionEncryptionKey.Wait: %w", err)
	}
	fmt.Fprintf(w, "Created database [%s] using multi-region encryption keys %q\n", dbObj.Name, dbObj.EncryptionConfig.GetKmsKeyNames())
	return nil
}

// [END spanner_create_database_with_MR_CMEK]
