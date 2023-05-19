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

// [START firestore_data_set_numeric_increment]
import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

// updateDocumentIncrement increments the population of the city document in the
// cities collection by 50.
func updateDocumentIncrement(projectID, city string) error {
	// projectID := "my-project"

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient: %w", err)
	}
	defer client.Close()

	dc := client.Collection("cities").Doc(city)
	_, err = dc.Update(ctx, []firestore.Update{
		{Path: "population", Value: firestore.Increment(50)},
	})
	if err != nil {
		return fmt.Errorf("Update: %w", err)
	}

	return nil
}

// [END firestore_data_set_numeric_increment]
