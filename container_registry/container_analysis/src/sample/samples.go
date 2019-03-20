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

// Samples for the Container Analysis golang libraries: https://cloud.google.com/container-registry/docs/container-analysis
package sample

import (
	"context"
	"fmt"
	"sync"
	"time"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1beta1"
	pubsub "cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
	grafeaspb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/grafeas"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/vulnerability"
)

// [START containeranalysis_create_note]

// createNote creates and returns a new vulnerability Note.
func createNote(projectID string) (*grafeaspb.Note, error) {
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return nil, err
	}
	
	projectName := fmt.Sprintf("projects/%s", projectID)

	req := &grafeaspb.CreateNoteRequest{
		Parent: projectName,
		NoteId: noteID,
		Note: &grafeaspb.Note{
			Type: &grafeaspb.Note_Vulnerability{
				// The 'Vulnerability' field can be modified to contain information about your vulnerability.
				Vulnerability: &vulnerability.Vulnerability{},
			},
		},
	}

	return client.CreateNote(ctx, req)
}

// [END containeranalysis_create_note]

// [START containeranalysis_create_occurrence]

// createsOccurrence creates and returns a new Occurrence of a previously created vulnerability Note.
func createOccurrence(resourceURL, noteID, occProjectID, noteProjectID string) (*grafeaspb.Occurrence, error) {
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return nil, err
	}
	
	req := &grafeaspb.CreateOccurrenceRequest{
		Parent: fmt.Sprintf("projects/%s", occProjectID),
		Occurrence: &grafeaspb.Occurrence{
			NoteName: fmt.Sprintf("projects/%s/notes/%s", noteProjectID, noteID),
			// Attach the occurrence to the associated resource uri.
			Resource: &grafeaspb.Resource{
				Uri: resourceURL,
			},
			// Details about the vulnerability instance can be added here.
			Details: &grafeaspb.Occurrence_Vulnerability{
				Vulnerability: &vulnerability.Details{},
			},
		},
	}
	return client.CreateOccurrence(ctx, req)
}

// [END containeranalysis_create_occurrence]

// [START containeranalysis_update_note]

// updateNote pushes an update to a Note that already exists on the server.
func updateNote(updated *grafeaspb.Note, noteID, projectID string) (*grafeaspb.Note, error) {
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return nil, err
	}
	
	req := &grafeaspb.UpdateNoteRequest{
		Name: fmt.Sprintf("projects/%s/notes/%s", projectID, noteID),
		Note: updated,
	}
	return client.UpdateNote(ctx, req)
}

// [END containeranalysis_update_note]

// [START containeranalysis_update_occurrence]

// updateOccurrences pushes an update to an Occurrence that already exists on the server.
func updateOccurrence(updated *grafeaspb.Occurrence, occurrenceID, projectID string) (*grafeaspb.Occurrence, error) {
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return nil, err
	}
	
	req := &grafeaspb.UpdateOccurrenceRequest{
		Name: fmt.Sprintf("projects/%s/occurrences/%s", projectID, occurrenceID),
		Occurrence: updated,
	}
	return client.UpdateOccurrence(ctx, req)
}

// [END containeranalysis_update_occurrence]

// [START containeranalysis_delete_note]

// deleteNote removes an existing Note from the server.
func deleteNote(noteID, projectID string) error {
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return err
	}

	req := &grafeaspb.DeleteNoteRequest{
		Name: fmt.Sprintf("projects/%s/notes/%s", projectID, noteID),
	}
	return client.DeleteNote(ctx, req)
}

// [END containeranalysis_delete_note]

// [START containeranalysis_delete_occurrence]

// deleteOccurrence removes an existing Occurrence from the server.
func deleteOccurrence(occurrenceID, projectID string) error {
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return err
	}
	
	req := &grafeaspb.DeleteOccurrenceRequest{
		Name: fmt.Sprintf("projects/%s/occurrences/%s", projectID, occurrenceID),
	}
	return client.DeleteOccurrence(ctx, req)
}

// [END containeranalysis_delete_occurrence]

// [START containeranalysis_get_note]

// getNote retrieves and prints a specified Note from the server.
func getNote(noteID, projectID string) (*grafeaspb.Note, error) {
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return nil, err
	}

	req := &grafeaspb.GetNoteRequest{
		Name: fmt.Sprintf("projects/%s/notes/%s", projectID, noteID),
	}
	note, err := client.GetNote(ctx, req)
	fmt.Println(note)
	return note, err
}

// [END containeranalysis_get_note]

// [START containeranalysis_get_occurrence]

// getOccurrence retrieves and prints a specified Occurrence from the server.
func getOccurrence(occurrenceID, projectID c  string) (*grafeaspb.Occurrence, error) {
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return nil, err
	}

	req := &grafeaspb.GetOccurrenceRequest{
		Name: fmt.Sprintf("projects/%s/occurrences/%s", projectID, occurrenceID),
	}
	occ, err := client.GetOccurrence(ctx, req)
	fmt.Println(occ)
	return occ, err
}

// [END containeranalysis_get_occurrence]

// [START containeranalysis_discovery_info]

// getDiscoveryInfo retrieves and prints the Discovery Occurrence created for a specified image.
// The Discovery Occurrence contains information about the initial scan on the image.
func getDiscoveryInfo(resourceURL, projectID string) error {
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return err
	}

	req := &grafeaspb.ListOccurrencesRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		Filter: fmt.Sprintf(`kind="DISCOVERY" AND resourceUrl=%q`, resourceURL),
	}
	it := client.ListOccurrences(ctx, req)
	for {
		occ, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(occ)
	}
	return nil
}

// [END containeranalysis_discovery_info]

// [START containeranalysis_occurrences_for_note]

// getOccurrencesForNote retrieves all the Occurrences associated with a specified Note.
// Here, all Occurrences are printed and counted.
func getOccurrencesForNote( noteID, projectID string) (int, error) {
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return -1, err
	}

	req := &grafeaspb.ListNoteOccurrencesRequest{
		Name: fmt.Sprintf("projects/%s/notes/%s", projectID, noteID),
	}
	it := client.ListNoteOccurrences(ctx, req)
	count := 0
	for {
		occ, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return -1, err
		}
		// Write custom code to process each Occurrence here.
		fmt.Println(occ)
		count = count + 1
	}
	return count, nil
}

// [END containeranalysis_occurrences_for_note]

// [START containeranalysis_occurrences_for_image]

// getOccurrencesForImage retrieves all the Occurrences associated with a specified image.
// Here, all Occurrences are simply printed and counted.
func getOccurrencesForImage(resourceURL, projectID string) (int, error) {
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return -1, err
	}
	
	req := &grafeaspb.ListOccurrencesRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		Filter: fmt.Sprintf("resourceUrl=%q", resourceURL),
	}
	it := client.ListOccurrences(ctx, req)
	count := 0
	for {
		occ, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return -1, err
		}
		// Write custom code to process each Occurrence here.
		fmt.Println(occ)
		count = count + 1
	}
	return count, nil
}

// [END containeranalysis_occurrences_for_image]

// [START containeranalysis_pubsub]

// occurrencePubsub handles incoming Occurrences using a Cloud Pub/Sub subscription.
func occurrencePubsub(subscriptionID string, timeout int, projectID string) (int, error) {
	ctx := context.Background()
	
	var mu sync.Mutex
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return -1, err
	}
	// Subscribe to the requested Pub/Sub channel.
	sub := client.Subscription(subscriptionID)
	count := 0

	// Listen to messages for 'timeout' seconds.
	toctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	err = sub.Receive(toctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		count = count + 1
		fmt.Printf("Message %d: %q\n", count, string(msg.Data))
		msg.Ack()
		mu.Unlock()
	})
	if err != nil {
		return -1, err
	}
	// Print and return the number of Pub/Sub messages received.
	fmt.Println(count)
	return count, nil
}

// createOccurrenceSubscription creates and returns a Pub/Sub subscription object listening to the Occurrence topic.
func createOccurrenceSubscription(subscriptionID, projectID string) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	// This topic id will automatically receive messages when Occurrences are added or modified
	topicID := "container-analysis-occurrences-v1beta1"
	topic := client.Topic(topicID)
	config := pubsub.SubscriptionConfig{Topic: topic}
	_, err = client.CreateSubscription(ctx, subscriptionID, config)
	return err
}

// [END containeranalysis_pubsub]
