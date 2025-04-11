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

package firestore

// [START firestore_query_filter_compound_multi_ineq]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
)

func multipleInequalitiesQuery(w io.Writer, projectID string) error {
	ctx := context.Background()

	// Create client
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient: %w", err)
	}
	defer client.Close()

	// Create query
	query := client.Collection("cities").
		Where("population", ">", 1000000).
		Where("density", "<", 10000)

	// Get documents
	docSnapshots, err := query.Documents(ctx).GetAll()
	for _, doc := range docSnapshots {
		fmt.Fprintln(w, doc.Data())
	}
	if err != nil {
		return fmt.Errorf("GetAll: %w", err)
	}

	return nil
}

// [END firestore_query_filter_compound_multi_ineq]
