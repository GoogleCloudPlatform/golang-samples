// Copyright 2022 Google LLC
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

// [START spanner_create_clients]
// [START spanner_postgresql_create_clients]

import (
	"context"
	"io"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
)

func createClients(w io.Writer, db string) error {
	ctx := context.Background()

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	dataClient, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer dataClient.Close()

	_ = adminClient
	_ = dataClient

	return nil
}

// [END spanner_postgresql_create_clients]
// [END spanner_create_clients]
