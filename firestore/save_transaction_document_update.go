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

// [START firestore_transaction_document_update]

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

func runSimpleTransaction(ctx context.Context, client *firestore.Client) error {
	// [START_EXCLUDE]
	// Initialize data as baseline for the multi-step transaction below.
	_, preErr := client.Collection("cities").Doc("SF").Set(ctx, map[string]interface{}{
		"population": 860000,
	})

	if preErr != nil {
		return preErr
	}
	// [END_EXCLUDE]

	ref := client.Collection("cities").Doc("SF")
	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref) // tx.Get, NOT ref.Get!
		if err != nil {
			return err
		}
		pop, err := doc.DataAt("population")
		if err != nil {
			return err
		}
		return tx.Set(ref, map[string]interface{}{
			"population": pop.(int64) + 1,
		}, firestore.MergeAll)
	})
	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
	}

	return err
}

// [END firestore_transaction_document_update]
