// Copyright 2019 Google LLC
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

// [START firestore_query_collection_group_dataset]
import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

// collectionGroupSetup sets up a collection group to query.
func collectionGroupSetup(projectID, cityCollection string) error {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient: %w", err)
	}
	defer client.Close()

	landmarks := []struct {
		city, name, t string
	}{
		{"SF", "Golden Gate Bridge", "bridge"},
		{"SF", "Legion of Honor", "museum"},
		{"LA", "Griffith Park", "park"},
		{"LA", "The Getty", "museum"},
		{"DC", "Lincoln Memorial", "memorial"},
		{"DC", "National Air and Space Museum", "museum"},
		{"TOK", "Ueno Park", "park"},
		{"TOK", "National Museum of Nature and Science", "museum"},
		{"BJ", "Jingshan Park", "park"},
		{"BJ", "Beijing Ancient Observatory", "museum"},
	}

	cities := client.Collection(cityCollection)
	for _, l := range landmarks {
		if _, err := cities.Doc(l.city).Collection("landmarks").NewDoc().Set(ctx, map[string]string{
			"name": l.name,
			"type": l.t,
		}); err != nil {
			return fmt.Errorf("Set: %w", err)
		}
	}

	return nil
}

// [END firestore_query_collection_group_dataset]
