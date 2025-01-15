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

// [START firestore_store_vectors]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
)

type CoffeeBean struct {
	Name           string             `firestore:"name,omitempty"`
	Description    string             `firestore:"description,omitempty"`
	EmbeddingField firestore.Vector32 `firestore:"embedding_field,omitempty"`
	Color          string             `firestore:"color,omitempty"`
}

func storeVectors(w io.Writer, projectID string) error {
	ctx := context.Background()

	// Create client
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient: %w", err)
	}
	defer client.Close()

	// Vector can be represented by Vector32 or Vector64
	doc := CoffeeBean{
		Name:           "Kahawa coffee beans",
		Description:    "Information about the Kahawa coffee beans.",
		EmbeddingField: []float32{1.0, 2.0, 3.0},
		Color:          "red",
	}
	ref := client.Collection("coffee-beans").NewDoc()
	if _, err = ref.Set(ctx, doc); err != nil {
		fmt.Fprintf(w, "failed to upsert: %v", err)
		return err
	}

	return nil
}

// [END firestore_store_vectors]
