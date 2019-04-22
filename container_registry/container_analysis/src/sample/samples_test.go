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
	"fmt"
	"math/rand"
	"path"
	"strconv"
	"testing"
	"time"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1beta1"
	pubsub "cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	discovery "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/discovery"
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
	timestamp string
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
	imageUrl := "gcr.io/" + timestamp + "-" + rand
	noteObj, err := createNote(noteID, projectID)
	if err != nil {
		t.Fatalf("createNote(%s): %v", noteID, err)
	}
	v := TestVariables{ctx, client, noteID, subID, imageUrl, projectID, noteObj, tryLimit, timestamp}
	return v
}

// Run after each test
// Removes any unneeded resources allocated for test
func teardown(t *testing.T, v TestVariables) {
	if err := deleteNote(v.noteID, v.projectID); err != nil {
		t.Log(err)
	}
}

func TestCreateNote(t *testing.T) {
	v := setup(t)

	newNote, err := getNote(v.noteID, v.projectID)
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

	if err := deleteNote(v.noteID, v.projectID); err != nil {
		t.Errorf("deleteNote(%s): %v", v.noteID, err)
	}
	deleted, err := getNote(v.noteID, v.projectID)
	if err == nil {
		t.Error("expected error from getNote; got nil")
	}
	if deleted != nil {
		t.Errorf("expected nil note; got %v", deleted)
	}

	teardown(t, v)
}

func TestCreateOccurrence(t *testing.T) {
	v := setup(t)

	created, err := createOccurrence(v.imageUrl, v.noteID, v.projectID, v.projectID)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageUrl, v.noteID, err)
	} else if created == nil {
		t.Error("returned occurrence is nil")
	} else {
		retrieved, err := getOccurrence(path.Base(created.Name), v.projectID)
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

	created, err := createOccurrence(v.imageUrl, v.noteID, v.projectID, v.projectID)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageUrl, v.noteID, err)
	} else if created == nil {
		t.Error("createOccurrence returns nil Occurrence object")
	} else {
		err = deleteOccurrence(path.Base(created.Name), v.projectID)
		if err != nil {
			t.Errorf("deleteOccurrence(%s): %v", created.Name, err)
		}
		deleted, err := getOccurrence(path.Base(created.Name), v.projectID)
		if err == nil {
			t.Error("getOccurrence returned nil error after DeleteOccurrence. expected error")
		}
		if deleted != nil {
			t.Errorf("getOccurrence returned occurrence after deletion: %v; expected nil", deleted)
		}
	}

	teardown(t, v)
}

func TestOccurrencesForImage(t *testing.T) {
	v := setup(t)

	origCount, err := getOccurrencesForImage(v.imageUrl, v.projectID)
	if err != nil {
		t.Errorf("getOccurrenceForImage(%s): %v", v.imageUrl, err)
	}
	if origCount != 0 {
		t.Errorf("unexpected initial number of occurrences: %d; want: %d", origCount, 0)
	}
	created, err := createOccurrence(v.imageUrl, v.noteID, v.projectID, v.projectID)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageUrl, v.noteID, err)
	} else if created == nil {
		t.Error("createOccurrence returns nil Occurrence object")
	}
	testutil.Retry(t, v.tryLimit, time.Second, func(r *testutil.R) {
		newCount, err := getOccurrencesForImage(v.imageUrl, v.projectID)
		if err != nil {
			r.Errorf("getOccurrencesForImage(%s): %v", v.imageUrl, err)
		}
		if newCount != 1 {
			r.Errorf("unexpected updated number of occurrences: %d; want: %d", newCount, 1)
		}
	})

	// Clean up
	deleteOccurrence(path.Base(created.Name), v.projectID)
	teardown(t, v)
}

func TestOccurrencesForNote(t *testing.T) {
	v := setup(t)

	origCount, err := getOccurrencesForNote(v.noteID, v.projectID)
	if err != nil {
		t.Errorf("getOccurrenceForNote(%s): %v", v.noteID, err)
	}
	if origCount != 0 {
		t.Errorf("unexpected initial number of occurrences: %d; want: %d", origCount, 0)
	}
	created, err := createOccurrence(v.imageUrl, v.noteID, v.projectID, v.projectID)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageUrl, v.noteID, err)
	} else if created == nil {
		t.Error("createOccurrence returns nil Occurrence object")
	}

	testutil.Retry(t, v.tryLimit, time.Second, func(r *testutil.R) {
		newCount, err := getOccurrencesForNote(v.noteID, v.projectID)
		if err != nil {
			r.Errorf("getOccurrencesForNote(%s): %v", v.noteID, err)
		}
		if newCount != 1 {
			r.Errorf("unexpected updated number of occurrences: %d; want: %d", newCount, 1)
		}
	})

	// Clean up
	deleteOccurrence(path.Base(created.Name), v.projectID)
	teardown(t, v)
}

func TestPubSub(t *testing.T) {
	t.Skip("Flaky: golang-samples#812")
	v := setup(t)
	// Create a new subscription if it doesn't exist.
	createOccurrenceSubscription(v.subID, v.projectID)

	testutil.Retry(t, v.tryLimit, time.Second, func(r *testutil.R) {
		// Use a channel and a goroutine to count incomming messages.
		c := make(chan int)
		go func() {
			count, err := occurrencePubsub(v.subID, 20, v.projectID)
			if err != nil {
				t.Errorf("occurrencePubsub(%s): %v", v.subID, err)
			}
			c <- count
		}()

		// Create some Occurrences.
		totalCreated := 3
		for i := 0; i < totalCreated; i++ {
			created, _ := createOccurrence(v.imageUrl, v.noteID, v.projectID, v.projectID)
			time.Sleep(time.Second)
			if err := deleteOccurrence(path.Base(created.Name), v.projectID); err != nil {
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

func TestPollDiscoveryOccurrenceFinished(t *testing.T) {
	v := setup(t)

	timeout := time.Duration(1) * time.Second
	discOcc, err := pollDiscoveryOccurrenceFinished(v.imageUrl, v.projectID, timeout)
	if err == nil || discOcc != nil {
		t.Errorf("expected error when resourceUrl has no discovery occurrence")
	}

	// create discovery occurrence
	noteId := "discovery-note-" + v.timestamp
	noteReq := &grafeaspb.CreateNoteRequest{
		Parent: fmt.Sprintf("projects/%s", v.projectID),
		NoteId: noteId,
		Note: &grafeaspb.Note{
			Type: &grafeaspb.Note_Discovery{
				Discovery: &discovery.Discovery{},
			},
		},
	}
	occReq := &grafeaspb.CreateOccurrenceRequest{
		Parent: fmt.Sprintf("projects/%s", v.projectID),
		Occurrence: &grafeaspb.Occurrence{
			NoteName: fmt.Sprintf("projects/%s/notes/%s", v.projectID, noteId),
			Resource: &grafeaspb.Resource{Uri: v.imageUrl},
			Details: &grafeaspb.Occurrence_Discovered{
				Discovered: &discovery.Details{
					Discovered: &discovery.Discovered{
						AnalysisStatus: discovery.Discovered_FINISHED_SUCCESS,
					},
				},
			},
		},
	}
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		t.Errorf("could not create client: %v", err)
	}
	defer client.Close()
	_, err = client.CreateNote(ctx, noteReq)
	created, err := client.CreateOccurrence(ctx, occReq)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageUrl, v.noteID, err)
	}

	// poll again
	timeout = time.Duration(20) * time.Second
	discOcc, err = pollDiscoveryOccurrenceFinished(v.imageUrl, v.projectID, timeout)
	if err != nil {
		t.Fatalf("error getting discovery occurrence: %v", err)
	}
	if discOcc == nil {
		t.Error("discovery occurrence is nil")
	} else {
		analysisStatus := discOcc.GetDiscovered().GetDiscovered().AnalysisStatus
		if analysisStatus != discovery.Discovered_FINISHED_SUCCESS {
			t.Errorf("discovery occurrence reported unexpected state: %sm want: %s", analysisStatus, discovery.Discovered_FINISHED_SUCCESS)
		}
	}

	// Clean up
	deleteOccurrence(path.Base(created.Name), v.projectID)
	deleteNote(noteId, v.projectID)
	teardown(t, v)
}

func TestFindVulnerabilitiesForImage(t *testing.T) {
	v := setup(t)

	occList, err := findVulnerabilityOccurrencesForImage(v.imageUrl, v.projectID)
	if err != nil {
		t.Fatalf("findVulnerabilityOccurrencesForImage(%v): %v", v.imageUrl, err)
	}
	if len(occList) != 0 {
		t.Errorf("unexpected initial number of vulnerabilities: %d; want: %d", len(occList), 0)
	}

	created, err := createOccurrence(v.imageUrl, v.noteID, v.projectID, v.projectID)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageUrl, v.noteID, err)
	} else if created == nil {
		t.Error("createOccurrence returns nil Occurrence object")
	}

	testutil.Retry(t, v.tryLimit, time.Second, func(r *testutil.R) {
		occList, err = findVulnerabilityOccurrencesForImage(v.imageUrl, v.projectID)
		if err != nil {
			r.Errorf("findVulnerabilityOccurrencesForImage(%v): %v", v.imageUrl, err)
		}
		if len(occList) != 1 {
			r.Errorf("unexpected updated number of occurrences: %d; want: %d", len(occList), 1)
		}
	})

	// Clean up
	deleteOccurrence(path.Base(created.Name), v.projectID)
	teardown(t, v)
}

func TestFindHighVulnerabilities(t *testing.T) {
	v := setup(t)

	// check before creation
	occList, err := findHighSeverityVulnerabilitiesForImage(v.imageUrl, v.projectID)
	if err != nil {
		t.Fatalf("findHighSeverityVulnerabilitiesForImage(%v): %v", v.imageUrl, err)
	}
	if len(occList) != 0 {
		t.Errorf("unexpected initial number of vulnerabilities: %d; want: %d", len(occList), 0)
	}

	// create high severity occurrence
	noteId := "severe-note-" + v.timestamp
	noteReq := &grafeaspb.CreateNoteRequest{
		Parent: fmt.Sprintf("projects/%s", v.projectID),
		NoteId: noteId,
		Note: &grafeaspb.Note{
			Type: &grafeaspb.Note_Vulnerability{
				Vulnerability: &vulnerability.Vulnerability{Severity: vulnerability.Severity_CRITICAL},
			},
		},
	}
	occReq := &grafeaspb.CreateOccurrenceRequest{
		Parent: fmt.Sprintf("projects/%s", v.projectID),
		Occurrence: &grafeaspb.Occurrence{
			NoteName: fmt.Sprintf("projects/%s/notes/%s", v.projectID, noteId),
			Resource: &grafeaspb.Resource{Uri: v.imageUrl},
			Details: &grafeaspb.Occurrence_Vulnerability{
				Vulnerability: &vulnerability.Details{Severity: vulnerability.Severity_CRITICAL},
			},
		},
	}
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		t.Errorf("could not create client: %v", err)
	}
	defer client.Close()
	_, err = client.CreateNote(ctx, noteReq)
	created, err := client.CreateOccurrence(ctx, occReq)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageUrl, v.noteID, err)
	} else if created == nil {
		t.Error("createOccurrence returns nil Occurrence object")
	}
	// check after creation
	testutil.Retry(t, v.tryLimit, time.Second, func(r *testutil.R) {
		occList, err = findHighSeverityVulnerabilitiesForImage(v.imageUrl, v.projectID)
		if err != nil {
			r.Errorf("findHighSeverityVulnerabilitiesForImage(%s): %v", v.imageUrl, err)
		}
		if len(occList) != 1 {
			r.Errorf("unexpected updated number of vulnerabilities: %d; want: %d", len(occList), 1)
		}
	})

	// Clean up
	deleteOccurrence(path.Base(created.Name), v.projectID)
	deleteNote(noteId, v.projectID)
	teardown(t, v)
}
