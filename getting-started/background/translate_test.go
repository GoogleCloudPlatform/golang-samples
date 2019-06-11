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

package background

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
)

func TestTranslate(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("Skipping Firestore test. Set GOLANG_SAMPLES_FIRESTORE_PROJECT.")
	}
	os.Setenv("GOOGLE_CLOUD_PROJECT", projectID)

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("firestore.NewClient: %v", err)
	}

	// Remove any old translations.
	if err := deleteAll(ctx, client, projectID); err != nil {
		t.Fatalf("deleteAll: %v", err)
	}

	msg, err := json.Marshal(Translation{
		Original: "Hello",
		Language: "fr",
	})
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	m := PubSubMessage{Data: msg}
	if err := Translate(ctx, m); err != nil {
		t.Fatalf("Translate: %v", err)
	}
	translations, err := getAll(ctx, client, projectID)
	if err != nil {
		t.Fatalf("getAll: %v", err)
	}
	if len(translations) != 1 {
		t.Fatalf("Translate got %d translations, want 1", len(translations))
	}
	want := Translation{
		Original:         "Hello",
		OriginalLanguage: "en",
		Language:         "fr",
		Translated:       "Bonjour",
	}
	if translations[0] != want {
		t.Fatalf("Translate got:\n%+v\nWant:\n%+v", translations[0], want)
	}

	// Ensure duplicate requests are only processed once.
	if err := Translate(ctx, m); err != nil {
		t.Fatalf("Translate: %v", err)
	}
	translations, err = getAll(ctx, client, projectID)
	if err != nil {
		t.Fatalf("getAll: %v", err)
	}
	if len(translations) != 1 {
		t.Fatalf("Translate got %d translations, want 1", len(translations))
	}
}

func deleteAll(ctx context.Context, client *firestore.Client, projectID string) error {
	docs, err := client.Collection("translations").Documents(ctx).GetAll()
	if err != nil {
		return err
	}
	for _, doc := range docs {
		if _, err := doc.Ref.Delete(ctx); err != nil {
			return err
		}
	}
	return nil
}

func getAll(ctx context.Context, client *firestore.Client, projectID string) ([]Translation, error) {
	docs, err := client.Collection("translations").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	var translations []Translation
	for _, doc := range docs {
		t := Translation{}
		if err := doc.DataTo(&t); err != nil {
			return nil, fmt.Errorf("DataTo: %v", err)
		}
		translations = append(translations, t)
	}
	return translations, nil
}
