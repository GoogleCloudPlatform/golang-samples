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

// [START bigtable_configure_connection_pool]

// Connectionpool is a sample program demonstrating how to configure the number
// of connection pools the Cloud Bigtable client should use.
package connectionpool

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigtable"
	"google.golang.org/api/option"
)

func configureConnectionPool(w io.Writer, projectID, instanceID string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance-id"
	ctx := context.Background()

	// Set up Bigtable data operations client.
	poolSize := 250
	client, err := bigtable.NewClient(ctx, projectID, instanceID,
		option.WithGRPCConnectionPool(poolSize))

	if err != nil {
		return fmt.Errorf("Could not create data operations client: %v", err)
	}

	fmt.Fprintf(w, "Connected with pool size of %d", poolSize)

	if err = client.Close(); err != nil {
		return fmt.Errorf("Could not close data operations client: %v", err)
	}

	return nil
}

// [END bigtable_configure_connection_pool]
