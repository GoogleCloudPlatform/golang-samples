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

// Package sample provides code samples for the Container Analysis libraries: https://cloud.google.com/container-registry/docs/container-analysis
package main

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"testing"
	"time"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1"
	pubsub "cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
	grafeaspb "google.golang.org/genproto/googleapis/grafeas/v1"
)

type TestVariables struct {
	ctx       context.Context
	client    *containeranalysis.Client
	noteID    string
	subID     string
	imageURL  string
	projectID string
	noteObj   *grafeaspb.Note
	tryLimit  int
	uuid      string
}

// Run before each test. Creates a set of useful variables
func setup(t *testing.T) TestVariables {
	tc := testutil.SystemTest(t)
	// Create client and context
	ctx := context.Background()
	client, _ := containeranalysis.NewClient(ctx)
	// Get unique id
	uuid, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("Could not generate uuid: %v", err)
	}
	uuidStr := uuid.String()
	// Set how many times to retry network tasks
	tryLimit := 20

	// Create variables used by tests
	projectID := tc.ProjectID
	noteID := "note-" + uuidStr
	subID := "occurrences-" + uuidStr
	imageURL := "https://gcr.io/" + uuidStr
	noteObj, err := createNote(noteID, projectID)
	if err != nil {
		t.Fatalf("createNote(%s): %v", noteID, err)
	}
	v := TestVariables{ctx, client, noteID, subID, imageURL, projectID, noteObj, tryLimit, uuidStr}
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

	created, err := createOccurrence(v.imageURL, v.noteID, v.projectID, v.projectID)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageURL, v.noteID, err)
	} else if created == nil {
		t.Error("returned occurrence is nil")
	}
	retrieved, err := getOccurrence(path.Base(created.Name), v.projectID)
	if err != nil {
		t.Errorf("getOccurrence(%s): %v", created.Name, err)
	} else if retrieved == nil {
		t.Error("getOccurrence returns nil Occurrence object")
	} else if retrieved.Name != created.Name {
		t.Errorf("created occurrence has wrong name: %s; want: %s", retrieved.Name, created.Name)
	}

	teardown(t, v)
}

func TestDeleteOccurrence(t *testing.T) {
	v := setup(t)

	created, err := createOccurrence(v.imageURL, v.noteID, v.projectID, v.projectID)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageURL, v.noteID, err)
	} else if created == nil {
		t.Error("createOccurrence returns nil Occurrence object")
	}
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

	teardown(t, v)
}

func TestOccurrencesForImage(t *testing.T) {
	v := setup(t)

	origCount, err := getOccurrencesForImage(new(bytes.Buffer), v.imageURL, v.projectID)
	if err != nil {
		t.Errorf("getOccurrenceForImage(%s): %v", v.imageURL, err)
	}
	if origCount != 0 {
		t.Errorf("unexpected initial number of occurrences: %d; want: %d", origCount, 0)
	}
	created, err := createOccurrence(v.imageURL, v.noteID, v.projectID, v.projectID)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageURL, v.noteID, err)
	} else if created == nil {
		t.Error("createOccurrence returns nil Occurrence object")
	}
	testutil.Retry(t, v.tryLimit, time.Second, func(r *testutil.R) {
		newCount, err := getOccurrencesForImage(new(bytes.Buffer), v.imageURL, v.projectID)
		if err != nil {
			r.Errorf("getOccurrencesForImage(%s): %v", v.imageURL, err)
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

	origCount, err := getOccurrencesForNote(new(bytes.Buffer), v.noteID, v.projectID)
	if err != nil {
		t.Errorf("getOccurrenceForNote(%s): %v", v.noteID, err)
	}
	if origCount != 0 {
		t.Errorf("unexpected initial number of occurrences: %d; want: %d", origCount, 0)
	}
	created, err := createOccurrence(v.imageURL, v.noteID, v.projectID, v.projectID)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageURL, v.noteID, err)
	} else if created == nil {
		t.Error("createOccurrence returns nil Occurrence object")
	}

	testutil.Retry(t, v.tryLimit, time.Second, func(r *testutil.R) {
		newCount, err := getOccurrencesForNote(new(bytes.Buffer), v.noteID, v.projectID)
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
	v := setup(t)
	// Create a new Topic if needed
	client, _ := pubsub.NewClient(v.ctx, v.projectID)
	topicID := "container-analysis-occurrences-v1"
	client.CreateTopic(v.ctx, topicID)

	// Create a new subscription if it doesn't exist.
	createOccurrenceSubscription(v.subID, v.projectID)

	testutil.Retry(t, v.tryLimit, time.Second, func(r *testutil.R) {
		// Use a channel and a goroutine to count incoming messages.
		c := make(chan int)
		go func() {
			count, err := occurrencePubsub(new(bytes.Buffer), v.subID, time.Duration(20)*time.Second, v.projectID)
			if err != nil {
				t.Errorf("occurrencePubsub(%s): %v", v.subID, err)
			}
			c <- count
		}()

		// Create some Occurrences.
		totalCreated := 3
		for i := 0; i < totalCreated; i++ {
			created, _ := createOccurrence(v.imageURL, v.noteID, v.projectID, v.projectID)
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
	sub := client.Subscription(v.subID)
	sub.Delete(v.ctx)
	teardown(t, v)
}

func TestPollDiscoveryOccurrenceFinished(t *testing.T) {
	v := setup(t)

	timeout := time.Duration(5) * time.Second
	discOcc, err := pollDiscoveryOccurrenceFinished(v.imageURL, v.projectID, timeout)
	if err == nil || discOcc != nil {
		t.Errorf("expected error when resourceURL has no discovery occurrence")
	}

	// create discovery occurrence
	noteID := "discovery-note-" + v.uuid
	noteReq := &grafeaspb.CreateNoteRequest{
		Parent: fmt.Sprintf("projects/%s", v.projectID),
		NoteId: noteID,
		Note: &grafeaspb.Note{
			Type: &grafeaspb.Note_Discovery{
				Discovery: &grafeaspb.DiscoveryNote{
					AnalysisKind: grafeaspb.NoteKind_DISCOVERY,
				},
			},
		},
	}
	occReq := &grafeaspb.CreateOccurrenceRequest{
		Parent: fmt.Sprintf("projects/%s", v.projectID),
		Occurrence: &grafeaspb.Occurrence{
			NoteName:    fmt.Sprintf("projects/%s/notes/%s", v.projectID, noteID),
			ResourceUri: v.imageURL,
			Details: &grafeaspb.Occurrence_Discovery{
				Discovery: &grafeaspb.DiscoveryOccurrence{
					AnalysisStatus: grafeaspb.DiscoveryOccurrence_FINISHED_SUCCESS,
				},
			},
		},
	}
	ctx := context.Background()
	client, err := containeranalysis.NewClient(ctx)
	if err != nil {
		t.Errorf("containeranalysis.NewClient: %v", err)
	}
	defer client.Close()
	_, err = client.GetGrafeasClient().CreateNote(ctx, noteReq)
	if err != nil {
		t.Errorf("createNote(%s): %v", v.noteID, err)
	}
	created, err := client.GetGrafeasClient().CreateOccurrence(ctx, occReq)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageURL, v.noteID, err)
	}

	// poll again
	testutil.Retry(t, v.tryLimit, time.Second, func(r *testutil.R) {
		discOcc, err = pollDiscoveryOccurrenceFinished(v.imageURL, v.projectID, timeout)
		if err != nil {
			r.Errorf("error getting discovery occurrence: %v", err)
		}
		if discOcc == nil {
			r.Errorf("discovery occurrence is nil")
		}
		analysisStatus := discOcc.GetDiscovery().GetAnalysisStatus()
		if analysisStatus != grafeaspb.DiscoveryOccurrence_FINISHED_SUCCESS {
			r.Errorf("discovery occurrence reported unexpected state: %s, want: %s", analysisStatus, grafeaspb.DiscoveryOccurrence_FINISHED_SUCCESS)
		}
	})

	// Clean up
	deleteOccurrence(path.Base(created.Name), v.projectID)
	deleteNote(noteID, v.projectID)
	teardown(t, v)
}

func TestFindVulnerabilitiesForImage(t *testing.T) {
	v := setup(t)

	occList, err := findVulnerabilityOccurrencesForImage(v.imageURL, v.projectID)
	if err != nil {
		t.Errorf("findVulnerabilityOccurrencesForImage(%v): %v", v.imageURL, err)
	}
	if len(occList) != 0 {
		t.Errorf("unexpected initial number of vulnerabilities: %d; want: %d", len(occList), 0)
	}

	created, err := createOccurrence(v.imageURL, v.noteID, v.projectID, v.projectID)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageURL, v.noteID, err)
	} else if created == nil {
		t.Error("createOccurrence returns nil Occurrence object")
	}

	testutil.Retry(t, v.tryLimit, time.Second, func(r *testutil.R) {
		occList, err = findVulnerabilityOccurrencesForImage(v.imageURL, v.projectID)
		if err != nil {
			r.Errorf("findVulnerabilityOccurrencesForImage(%v): %v", v.imageURL, err)
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
	occList, err := findHighSeverityVulnerabilitiesForImage(v.imageURL, v.projectID)
	if err != nil {
		t.Errorf("findHighSeverityVulnerabilitiesForImage(%v): %v", v.imageURL, err)
	}
	if len(occList) != 0 {
		t.Errorf("unexpected initial number of vulnerabilities: %d; want: %d", len(occList), 0)
	}

	// create high severity occurrence
	noteID := "severe-note-" + v.uuid
	noteReq := &grafeaspb.CreateNoteRequest{
		Parent: fmt.Sprintf("projects/%s", v.projectID),
		NoteId: noteID,
		Note: &grafeaspb.Note{
			Type: &grafeaspb.Note_Vulnerability{
				Vulnerability: &grafeaspb.VulnerabilityNote{
					Severity: grafeaspb.Severity_CRITICAL,
					Details: []*grafeaspb.VulnerabilityNote_Detail{
						{
							AffectedCpeUri:  "your-uri-here",
							AffectedPackage: "your-package-here",
							AffectedVersionStart: &grafeaspb.Version{
								Kind: grafeaspb.Version_MINIMUM,
							},
							AffectedVersionEnd: &grafeaspb.Version{
								Kind: grafeaspb.Version_MAXIMUM,
							},
						},
					},
				},
			},
		},
	}
	occReq := &grafeaspb.CreateOccurrenceRequest{
		Parent: fmt.Sprintf("projects/%s", v.projectID),
		Occurrence: &grafeaspb.Occurrence{
			NoteName:    fmt.Sprintf("projects/%s/notes/%s", v.projectID, noteID),
			ResourceUri: v.imageURL,
			Details: &grafeaspb.Occurrence_Vulnerability{
				Vulnerability: &grafeaspb.VulnerabilityOccurrence{
					EffectiveSeverity: grafeaspb.Severity_CRITICAL,
					PackageIssue: []*grafeaspb.VulnerabilityOccurrence_PackageIssue{
						{
							AffectedCpeUri:  "your-uri-here",
							AffectedPackage: "your-package-here",
							AffectedVersion: &grafeaspb.Version{
								Kind: grafeaspb.Version_MINIMUM,
							},
							FixedVersion: &grafeaspb.Version{
								Kind: grafeaspb.Version_MAXIMUM,
							},
						},
					},
				},
			},
		},
	}
	ctx := context.Background()
	client, err := containeranalysis.NewClient(ctx)
	if err != nil {
		t.Errorf("could not create client: %v", err)
	}
	defer client.Close()
	_, err = client.GetGrafeasClient().CreateNote(ctx, noteReq)
	if err != nil {
		t.Errorf("createNote(%s): %v", v.noteID, err)
	}
	created, err := client.GetGrafeasClient().CreateOccurrence(ctx, occReq)
	if err != nil {
		t.Errorf("createOccurrence(%s, %s): %v", v.imageURL, v.noteID, err)
	} else if created == nil {
		t.Error("createOccurrence returns nil Occurrence object")
	}
	// check after creation
	testutil.Retry(t, v.tryLimit, time.Second, func(r *testutil.R) {
		occList, err = findHighSeverityVulnerabilitiesForImage(v.imageURL, v.projectID)
		if err != nil {
			r.Errorf("findHighSeverityVulnerabilitiesForImage(%s): %v", v.imageURL, err)
		}
		if len(occList) != 1 {
			r.Errorf("unexpected updated number of vulnerabilities: %d; want: %d", len(occList), 1)
		}
	})

	// Clean up
	deleteOccurrence(path.Base(created.Name), v.projectID)
	deleteNote(noteID, v.projectID)
	teardown(t, v)
}
