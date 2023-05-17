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

// [START firestore_data_set_from_map_nested]

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
)

func addDocDataTypes(ctx context.Context, client *firestore.Client) error {
	doc := make(map[string]interface{})
	doc["stringExample"] = "Hello world!"
	doc["booleanExample"] = true
	doc["numberExample"] = 3.14159265
	doc["dateExample"] = time.Now()
	doc["arrayExample"] = []interface{}{5, true, "hello"}
	doc["nullExample"] = nil
	doc["objectExample"] = map[string]interface{}{
		"a": 5,
		"b": true,
	}

	_, err := client.Collection("data").Doc("one").Set(ctx, doc)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}

	return err
}

// [END firestore_data_set_from_map_nested]
