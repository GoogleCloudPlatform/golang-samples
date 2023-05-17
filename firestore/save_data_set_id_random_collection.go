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

// [START firestore_data_set_id_random_collection]

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

func addDocWithoutID(ctx context.Context, client *firestore.Client) error {
	_, _, err := client.Collection("cities").Add(ctx, map[string]interface{}{
		"name":    "Tokyo",
		"country": "Japan",
	})
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}

	return err
}

// [END firestore_data_set_id_random_collection]
