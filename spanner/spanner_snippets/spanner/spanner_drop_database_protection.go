// Copyright 2023 Google LLC
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

// [START spanner_drop_database_protection]
import (
	"context"
	"fmt"
	"io"
	"regexp"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/api/option"
	"google.golang.org/genproto/protobuf/field_mask"
)

func dropDatabaseProtection(ctx context.Context, w io.Writer, db string) error {
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("Invalid database id %s", db)
	}
	// TODO(sriharshach): Remove endpoint
	endpoint := "staging-wrenchworks.sandbox.googleapis.com:443"
	adminClient, err := database.NewDatabaseAdminClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return err
	}
	defer adminClient.Close()

	op, err := adminClient.UpdateDatabase(ctx, &adminpb.UpdateDatabaseRequest{
		Database: &adminpb.Database{
			Name:                 db,
			EnableDropProtection: true,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"enable_drop_protection"},
		},
	})
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Fprintf(w, "Enabled drop database protection to database [%s]\n", db)
	return nil
}

// [END spanner_drop_database_protection]
