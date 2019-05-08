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
	"os"
	"testing"

	"cloud.google.com/go/firestore"
)

func deleteAll(projectID string) error {
	client, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		return err
	}
	docs, err := client.Collection("translations").Documents(context.Background()).GetAll()
	if err != nil {
		return err
	}
	for _, doc := range docs {
		if _, err := doc.Ref.Delete(context.Background()); err != nil {
			return err
		}
	}
	return nil
}

func getAll(projectID string) ([]Translation, error) {
	client, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	docs, err := client.Collection("translations").Documents(context.Background()).GetAll()
	if err != nil {
		return nil, err
	}
	translations := []Translation{}
	for _, doc := range docs {
		t := Translation{}
		doc.DataTo(&t)
		translations = append(translations, t)
	}
	return translations, nil
}

func TestTranslate(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	if projectID == "" {
		t.Skip("Skipping Firestore test. Set GOLANG_SAMPLES_FIRESTORE_PROJECT.")
	}
	os.Setenv("GOOGLE_CLOUD_PROJECT", projectID)

	// Remove any old translations.
	if err := deleteAll(projectID); err != nil {
		t.Fatalf("deleteAll: %v", err)
	}

	ctx := context.Background()
	msg, _ := json.Marshal(Translation{
		Original: "Hello",
		Language: "fr",
	})
	m := PubSubMessage{Data: msg}
	if err := Translate(ctx, m); err != nil {
		t.Fatalf("Translate: %v", err)
	}
	translations, err := getAll(projectID)
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
}
