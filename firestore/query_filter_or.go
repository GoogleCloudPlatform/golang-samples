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

package firestore

// [START firestore_query_filter_or]
import (
	"context"
	"fmt"
	"io"

	firestore "cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func queryFilterOr(w io.Writer, projectId string) error {
	// Instantiate a client
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		return err
	}
	// always be sure to close the client to release resources
	defer client.Close()

	q1 := firestore.PropertyFilter{
		Path:     "birthYear",
		Operator: "==",
		Value:    1906,
	}

	q2 := firestore.PropertyFilter{
		Path:     "birthYear",
		Operator: "==",
		Value:    1815,
	}

	orFilter := firestore.OrFilter{
		Filters: []firestore.EntityFilter{q1, q2},
	}

	orQuery := client.Collection("users").WhereEntity(orFilter)
	it, err := orQuery.Documents(ctx)
	if err != nil {
		return err
	}

	fmt.Fprint(w, "Individual documents:\n")
	for {
		doc, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("documents iterator: %w", err)
		}
		fmt.Fprintf(w, "%s: %s", doc.Ref.ID, doc.Data()["birthYear"])
	}

	return nil
}

// [END firestore_query_filter_or]
