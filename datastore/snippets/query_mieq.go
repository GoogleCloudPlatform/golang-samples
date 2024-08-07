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

// Package datastore_snippets contains snippet code for the Cloud Datastore API.
// The code is not runnable.

package datastore_snippets

// [START datastore_query_filter_compound_multi_ineq]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

func queryMultipleInequality(w io.Writer, projectId string) error {
	// Create client
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectId)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	// Create query
	q := datastore.NewQuery("Task").
		FilterField("priority", ">", 4).
		FilterField("days", "<", 3)

		// Run query
	it := client.Run(ctx, q)
	for {
		var dst struct {
			Priority    int
			Days        int
			Done        bool
			Category    string
			Description string
		}
		key, err := it.Next(&dst)
		if err == iterator.Done {
			break
		}

		if err != nil {
			return fmt.Errorf("Next: %w", err)
		}
		fmt.Fprintf(w, "Key: %v. Entity: %v\n", key, dst)
	}

	return nil
}

// [END datastore_query_filter_compound_multi_ineq]
