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

package sample

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1beta1"
	pubsub "cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	grafeaspb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/grafeas"
	vulnerability "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/vulnerability"
)

type TestVariables struct {
	ctx       context.Context
	client    *containeranalysis.GrafeasV1Beta1Client
	noteID    string
	subID     string
	imageUrl  string
	projectID string
	noteObj   *grafeaspb.Note
	tryLimit  int
}

// Run before each test. Creates a set of useful variables
func setup(t *testing.T) TestVariables {
	tc := testutil.SystemTest(t)
	// Create client and context
	ctx := context.Background()
	client, _ := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	// Get current timestamp
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	// Make a random portion so each test is unique
	rand := strconv.Itoa(rand.Int())
	// Set how many times to retry network tasks
	tryLimit := 20

	// Create variables used by tests
	projectID := tc.ProjectID
	noteID := "note-" + timestamp + "-" + rand
	subID := "CA-Occurrences-" + timestamp + "-" + rand
	imageUrl := "www." + timestamp + "-" + rand + ".com"
	noteObj, err := createNote(ctx, client, noteID, projectID)
	if err != nil {
		t.Fatalf("createNote(%s): %v", noteID, err)
	}
	v := TestVariables{ctx, client, noteID, subID, imageUrl, projectID, noteObj, tryLimit}
	return v
}

// Run after each test
// Removes any unneeded resources allocated for test
func teardown(t *testing.T, v TestVariables) {
	if err := deleteNote(v.ctx, v.client, v.noteID, v.projectID); err != nil {
		t.Log(err)
	}
}

func TestCreateNote(t *testing.T) {
	v := setup(t)

	newNote, err := getNote(v.ctx, v.client, v.noteID, v.projectID)
	if err != nil {
		t.Errorf("getNote(%s): %v", v.noteID, err)
	} else if newNote == nil {
		t.Error("created note is nil")
	} else if newNote.Name != v.noteObj.Name {
		t.Errorf("created note has wrong name: %s; want: %s", newNote.Name, v.noteObj.Name)
	}

	teardown(t, v)
}

func TestDeleteNote(t *testing.T) {
	v := setup(t)

	if err := deleteNote(v.ctx, v.client, v.noteID, v.projectID); err != nil {
		t.Errorf("deleteNote(%s): %v", v.noteID, err)
	}
	deleted, err := getNote(v.ctx, v.client, v.noteID, v.projectID)
	if err == nil {
		t.Error("expected error from getNote; got nil")
	}
	if deleted != nil {
		t.Errorf("expected nil note; got %v", deleted)
	}

	teardown(t, v)
}

func TestUpdateNote(t *testing.T) {
	v := setup(t)

	description := "updated"
	v.noteObj.ShortDescription = description
	returned, err := updateNote(v.ctx, v.client, v.noteObj, v.noteID, v.projectID)
	if err != nil {
		t.Errorf("updateNote(%s): %v", v.noteID, err)
	} else if returned.ShortDescription != description {
		t.Errorf("returned note doesn't contain requested description text: %s; want: %s", returned.ShortDescription, description)
	}
	updated, err := getNote(v.ctx, v.client, v.noteID, v.projectID)
	if err != nil {
		t.Errorf("getNote(%s): %v", v.noteID, err)
	} else if updated == nil {
		t.Error("GetNote after update returns nil Note object")
	} else if updated.ShortDescription != description {
		t.Errorf("updated note doesn't contain requested description text: %s; want: %s", updated.ShortDescription, description)
	}

	teardown(t, v)
}

func TestCreateOccurrence(t *testing.T) {
	v := setup(t)

	created, err := createOccurrence(v.ctx, v.client, v.imageUrl, v.noteID, v.projectID, v.projectID)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageUrl, v.noteID, err)
	} else if created == nil {
		t.Error("returned occurrence is nil")
	} else {
		retrieved, err := getOccurrence(v.ctx, v.client, created.Name)
		if err != nil {
			t.Errorf("getOccurrence(%s): %v", created.Name, err)
		} else if retrieved == nil {
			t.Error("getOccurrence returns nil Occurrence object")
		} else if retrieved.Name != created.Name {
			t.Errorf("created occurrence has wrong name: %s; want: %s", retrieved.Name, created.Name)
		}
	}

	teardown(t, v)
}

func TestDeleteOccurrence(t *testing.T) {
	v := setup(t)

	created, err := createOccurrence(v.ctx, v.client, v.imageUrl, v.noteID, v.projectID, v.projectID)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageUrl, v.noteID, err)
	} else if created == nil {
		t.Error("createOccurrence returns nil Occurrence object")
	} else {
		err = deleteOccurrence(v.ctx, v.client, created.Name)
		if err != nil {
			t.Errorf("deleteOccurrence(%s): %v", created.Name, err)
		}
		deleted, err := getOccurrence(v.ctx, v.client, created.Name)
		if err == nil {
			t.Error("getOccurrence returned nil error after DeleteOccurrence. expected error")
		}
		if deleted != nil {
			t.Errorf("getOccurrence returned occurrence after deletion: %v; expected nil", deleted)
		}
	}

	teardown(t, v)
}

func TestUpdateOccurrence(t *testing.T) {
	v := setup(t)

	created, err := createOccurrence(v.ctx, v.client, v.imageUrl, v.noteID, v.projectID, v.projectID)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageUrl, v.noteID, err)
	} else if created == nil {
		t.Error("createOccurrence returns nil Occurrence object")
	} else {
		newType := "updated"

		details := vulnerability.Details{Type: newType}
		vulDetails := grafeaspb.Occurrence_Vulnerability{Vulnerability: &details}
		resource := grafeaspb.Resource{Uri: created.Resource.Uri}
		occurrence := grafeaspb.Occurrence{NoteName: created.NoteName, Resource: &resource, Details: &vulDetails}

		returned, err := updateOccurrence(v.ctx, v.client, &occurrence, created.Name)
		if err != nil {
			t.Errorf("updateOccurrence(%s): %v", created.Name, err)
		} else if returned.GetVulnerability().Type != newType {
			t.Errorf("returned occurrence doesn't contain requested vulnerability type: %s; want: %s", returned.GetVulnerability().Type, newType)
		}
		retrieved, err := getOccurrence(v.ctx, v.client, created.Name)
		if err != nil {
			t.Errorf("getOccurrence(%s): %v", created.Name, err)
		} else if retrieved == nil {
			t.Error("GetOccurrence returned nil Occurrence object")
		} else if retrieved.GetVulnerability().Type != newType {
			t.Errorf("updated occurrence doesn't contain requested vulnerability type: %s; want: %s", retrieved.GetVulnerability().Type, newType)
		}
	}
	teardown(t, v)
}

func TestOccurrencesForImage(t *testing.T) {
	v := setup(t)

	origCount, err := getOccurrencesForImage(v.ctx, v.client, v.imageUrl, v.projectID)
	if err != nil {
		t.Errorf("getOccurrenceForImage(%s): %v", v.imageUrl, err)
	}
	if origCount != 0 {
		t.Errorf("unexpected initial number of occurrences: %d; want: %d", origCount, 0)
	}
	created, err := createOccurrence(v.ctx, v.client, v.imageUrl, v.noteID, v.projectID, v.projectID)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageUrl, v.noteID, err)
	} else if created == nil {
		t.Error("createOccurrence returns nil Occurrence object")
	}
	testutil.Retry(t, v.tryLimit, time.Second, func(r *testutil.R) {
		newCount, err := getOccurrencesForImage(v.ctx, v.client, v.imageUrl, v.projectID)
		if err != nil {
			r.Errorf("getOccurrencesForImage(%s): %v", v.imageUrl, err)
		}
		if newCount != 1 {
			r.Errorf("unexpected updated number of occurrences: %d; want: %d", newCount, 1)
		}
	})

	// Clean up
	deleteOccurrence(v.ctx, v.client, created.Name)
	teardown(t, v)
}

func TestOccurrencesForNote(t *testing.T) {
	v := setup(t)

	origCount, err := getOccurrencesForNote(v.ctx, v.client, v.noteID, v.projectID)
	if err != nil {
		t.Errorf("getOccurrenceForNote(%s): %v", v.noteID, err)
	}
	if origCount != 0 {
		t.Errorf("unexpected initial number of occurrences: %d; want: %d", origCount, 0)
	}
	created, err := createOccurrence(v.ctx, v.client, v.imageUrl, v.noteID, v.projectID, v.projectID)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageUrl, v.noteID, err)
	} else if created == nil {
		t.Error("createOccurrence returns nil Occurrence object")
	}

	testutil.Retry(t, v.tryLimit, time.Second, func(r *testutil.R) {
		newCount, err := getOccurrencesForNote(v.ctx, v.client, v.noteID, v.projectID)
		if err != nil {
			r.Errorf("getOccurrencesForNote(%s): %v", v.noteID, err)
		}
		if newCount != 1 {
			r.Errorf("unexpected updated number of occurrences: %d; want: %d", newCount, 1)
		}
	})

	// Clean up
	deleteOccurrence(v.ctx, v.client, created.Name)
	teardown(t, v)
}

func TestPubSub(t *testing.T) {
	v := setup(t)
	// Create a new subscription if it doesn't exist.
	createOccurrenceSubscription(v.ctx, v.subID, v.projectID)

	testutil.Retry(t, v.tryLimit, time.Second, func(r *testutil.R) {
		// Use a channel and a goroutine to count incomming messages.
		c := make(chan int)
		go func() {
			count, err := occurrencePubsub(v.ctx, v.subID, 20, v.projectID)
			if err != nil {
				t.Errorf("occurrencePubsub(%s): %v", v.subID, err)
			}
			c <- count
		}()

		// Create some Occurrences.
		totalCreated := 3
		for i := 0; i < totalCreated; i++ {
			created, _ := createOccurrence(v.ctx, v.client, v.imageUrl, v.noteID, v.projectID, v.projectID)
			time.Sleep(time.Second)
			if err := deleteOccurrence(v.ctx, v.client, created.Name); err != nil {
				t.Errorf("deleteOccurrence(%s): %v", created.Name, err)
			}
			time.Sleep(time.Second)
		}
		result := <-c
		if result != totalCreated {
			r.Errorf("invalid occurrence count: %d; want: %d", result, totalCreated)
		}
	})

	// Clean up
	client, _ := pubsub.NewClient(v.ctx, v.projectID)
	sub := client.Subscription(v.subID)
	sub.Delete(v.ctx)
	teardown(t, v)
}
