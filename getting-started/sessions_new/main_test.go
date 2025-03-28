// Copyright 2025 Google LLC
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

package main

import (
	"context"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"
)

// TestIndex checks if simulating the request twice by reusing the first request increases the counter.
func TestIndex(t *testing.T) {

	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	collectionID := "test-hello-views"

	// Create new app
	a, err := newApp(projectID, collectionID)
	if err != nil {
		t.Fatalf("newApp: %v", err)
	}

	// Simulate HTTP GET request
	r := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	a.index(rr, r)

	// ResponseWriter body should contain 1 view
	if got, want := rr.Body.String(), "1 view"; !strings.Contains(got, want) {
		t.Errorf("index first visit got:\n----\n%v\n----\nWant to contain %q", got, want)
	}

	// Subsequent requests include the cookie from first visit, so it is assigned to the new request
	r = httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Cookie", rr.Header().Get("Set-Cookie"))

	rr = httptest.NewRecorder()

	// Simulate another HTTP GET request
	a.index(rr, r)

	if got, want := rr.Body.String(), "2 views"; !strings.Contains(got, want) {
		t.Errorf("index second visit got:\n----\n%v\n----\nWant to contain %q", got, want)
	}

	cleanup(t, projectID, collectionID)
}

// TestIndexCorrupted checks if changing the cookie's value to an invalid one resets the counter.
func TestIndexCorrupted(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_FIRESTORE_PROJECT")
	collectionID := "test-hello-views"

	// Create new app
	a, err := newApp(projectID, collectionID)
	if err != nil {
		t.Fatalf("newApp: %v", err)
	}

	// Simulate HTTP GET request
	r := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	a.index(rr, r)

	// ResponseWriter body should contain 1 view
	if got, want := rr.Body.String(), "1 view"; !strings.Contains(got, want) {
		t.Errorf("index first visit got:\n----\n%v\n----\nWant to contain %q", got, want)
	}

	// Simulate HTTP Get request but removing the assigned session ID
	r = httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Cookie", "this is not a valid session ID")

	rr = httptest.NewRecorder()

	a.index(rr, r)

	// As the current session ID is not valid, it should contain 1
	if got, want := rr.Body.String(), "1 view"; !strings.Contains(got, want) {
		t.Errorf("index first visit got:\n----\n%v\n----\nWant to contain %q", got, want)
	}

	cleanup(t, projectID, collectionID)
}

// cleanup function deletes all documents inside a collection
func cleanup(t *testing.T, projectID, collectionID string) {

	t.Helper()

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}

	iter := client.Collection(collectionID).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err != nil {
			// Handle the case where the collection might not exist or other errors
			if err.Error() == "iterator ended" {
				log.Printf("Collection %s cleaned up or did not exist.", collectionID)
				return
			}
			log.Printf("Error iterating documents in %s: %v", collectionID, err)
			return
		}
		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			log.Printf("Error deleting document %s in %s: %v", doc.Ref.ID, collectionID, err)
		}
	}
}
