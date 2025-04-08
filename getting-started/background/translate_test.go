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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

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

	maxRetries := 5
	sleep := 5 * time.Second
	failureLog := &bytes.Buffer{}
	errorf := func(format string, args ...interface{}) {
		fmt.Fprintf(failureLog, format, args...)
	}
	for retry := 0; retry < maxRetries; retry++ {
		failureLog.Reset()
		pass := func() bool {
			// Remove any old translations.
			if err := deleteAll(ctx, client, projectID); err != nil {
				errorf("deleteAll: %v", err)
				return false
			}

			msg, err := json.Marshal(Translation{
				Original: "Me",
				Language: "fr",
			})
			if err != nil {
				errorf("json.Marshal: %v", err)
				return false
			}
			m := PubSubMessage{Data: msg}
			if err := Translate(ctx, m); err != nil {
				errorf("Translate: %v", err)
				return false
			}
			translations, err := getAll(ctx, client, projectID)
			if err != nil {
				errorf("getAll: %v", err)
				return false
			}
			if len(translations) != 1 {
				errorf("Translate got %d translations, want 1", len(translations))
				return false
			}
			want := Translation{
				Original:         "Me",
				OriginalLanguage: "en",
				Language:         "fr",
				Translated:       "Moi",
			}
			if translations[0] != want {
				errorf("Translate got:\n%+v\nWant:\n%+v", translations[0], want)
				return false
			}

			// Ensure duplicate requests are only processed once.
			if err := Translate(ctx, m); err != nil {
				errorf("Translate: %v", err)
				return false
			}
			translations, err = getAll(ctx, client, projectID)
			if err != nil {
				errorf("getAll: %v", err)
				return false
			}
			if len(translations) != 1 {
				errorf("Translate got %d translations, want 1", len(translations))
				return false
			}
			return true
		}()
		if pass {
			return
		}
		if retry < maxRetries-1 {
			time.Sleep(sleep)
		}
	}
	t.Fatalf("Translate failed after %d attempts: %v", maxRetries, failureLog.String())
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
