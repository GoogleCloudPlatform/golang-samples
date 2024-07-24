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

// [START firestore_transaction_document_update_conditional]

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/firestore"
)

func infoTransaction(ctx context.Context, client *firestore.Client) (int64, error) {
	var updatedPop int64
	ref := client.Collection("cities").Doc("SF")
	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref)
		if err != nil {
			return err
		}
		pop, err := doc.DataAt("population")
		if err != nil {
			return err
		}
		newpop := pop.(int64) + 1
		if newpop <= 1000000 {
			err := tx.Set(ref, map[string]interface{}{
				"population": newpop,
			}, firestore.MergeAll)
			if err == nil {
				updatedPop = newpop
			}
			return err
		}
		return errors.New("population is too big")
	})
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}
	return updatedPop, err
}

// [END firestore_transaction_document_update_conditional]
