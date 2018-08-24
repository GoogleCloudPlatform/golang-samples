// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Samples for the Container Analysis golang libnraries: https://cloud.google.com/container-registry/docs/container-analysis
package main

import (
	"fmt"
	"sync"
	"time"

	containeranalysis "cloud.google.com/go/devtools/containeranalysis/apiv1beta1"
	pubsub "cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	grafeaspb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/grafeas"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/vulnerability"
)

// [START create_note]

//createNote creates and returns a new vulnerability Note.
func createNote(ctx context.Context, client *containeranalysis.GrafeasV1Beta1Client, noteID, projectID string) (*grafeaspb.Note, error) {
	parent := "projects/" + projectID

	vul := vulnerability.Vulnerability{}
	// Details about the your vulnerability can be added here.
	// For Example: vul.Severity = vulnerability.Severity_CRITICAL
	noteVul := grafeaspb.Note_Vulnerability{Vulnerability: &vul}
	note := grafeaspb.Note{Type: &noteVul}
	req := &grafeaspb.CreateNoteRequest{Parent: parent, NoteId: noteID, Note: &note}
	return client.CreateNote(ctx, req)
}

// [END create_note]

// [START create_occurrence]

//createsOccurrence creates and returns a new Occurrence of a previously created vulnerability Note.
func createOccurrence(ctx context.Context, client *containeranalysis.GrafeasV1Beta1Client, imageUrl, parentNoteID, projectID string) (*grafeaspb.Occurrence, error) {
	noteName := "projects/" + projectID + "/notes/" + parentNoteID
	projectName := "projects/" + projectID

	details := vulnerability.Details{}
	// Details about the vulnerability instance can be added here.
	vulDetails := grafeaspb.Occurrence_Vulnerability{Vulnerability: &details}
	// Attach the occurrence to the associated image uri.
	resource := grafeaspb.Resource{Uri: imageUrl}
	occurrence := grafeaspb.Occurrence{NoteName: noteName, Resource: &resource, Details: &vulDetails}
	req := &grafeaspb.CreateOccurrenceRequest{Parent: projectName, Occurrence: &occurrence}
	return client.CreateOccurrence(ctx, req)
}

// [END create_occurrence]

// [START update_note]

// updateNote pushes an update to a Note that already exists on the server.
func updateNote(ctx context.Context, client *containeranalysis.GrafeasV1Beta1Client, updated *grafeaspb.Note, noteID, projectID string) error {
	noteName := "projects/" + projectID + "/notes/" + noteID

	req := &grafeaspb.UpdateNoteRequest{Name: noteName, Note: updated}
	_, err := client.UpdateNote(ctx, req)
	return err
}

// [END update_note]

// [START update_occurrence]

// updateOccurrences pushes an update to an Occurrence that already exists on the server.
func updateOccurrence(ctx context.Context, client *containeranalysis.GrafeasV1Beta1Client, updated *grafeaspb.Occurrence, occurrenceName string) error {
	req := &grafeaspb.UpdateOccurrenceRequest{Name: occurrenceName, Occurrence: updated}
	_, err := client.UpdateOccurrence(ctx, req)
	return err
}

// [END update_occurrence]

// [START delete_note]

// deleteNote removes an existing Note from the server.
func deleteNote(ctx context.Context, client *containeranalysis.GrafeasV1Beta1Client, noteID, projectID string) error {
	noteName := "projects/" + projectID + "/notes/" + noteID

	req := &grafeaspb.DeleteNoteRequest{Name: noteName}
	return client.DeleteNote(ctx, req)
}

// [END delete_note]

// [START delete_occurrence]

// deleteOccurrence removes an existing Occurrence from the server.
func deleteOccurrence(ctx context.Context, client *containeranalysis.GrafeasV1Beta1Client, occurrenceName string) error {
	req := &grafeaspb.DeleteOccurrenceRequest{Name: occurrenceName}
	return client.DeleteOccurrence(ctx, req)
}

// [END delete_occurrence]

// [START get_note]

// getNote retrieves a specified Note from the server.
func getNote(ctx context.Context, client *containeranalysis.GrafeasV1Beta1Client, noteID, projectID string) (*grafeaspb.Note, error) {
	noteName := "projects/" + projectID + "/notes/" + noteID
	req := &grafeaspb.GetNoteRequest{Name: noteName}
	return client.GetNote(ctx, req)
}

// [END get_note]

// [START get_occurrence]

// getOccurrence retrieves an Occurrence from the server.
func getOccurrence(ctx context.Context, client *containeranalysis.GrafeasV1Beta1Client, occurrenceName string) (*grafeaspb.Occurrence, error) {
	req := &grafeaspb.GetOccurrenceRequest{Name: occurrenceName}
	return client.GetOccurrence(ctx, req)
}

// [END get_occurrence]

// [START discovery_info]

// getDiscoveryInfo retrieves and prints the Discovery Occurrence created for a specified image.
// The Discovery Occurrence contains information about the initial scan on the image.
func getDiscoveryInfo(ctx context.Context, client *containeranalysis.GrafeasV1Beta1Client, imageUrl, projectID string) error {
	filterStr := `kind="DISCOVERY" AND resourceUrl="` + imageUrl + `"`
	projectName := "projects/" + projectID

	req := &grafeaspb.ListOccurrencesRequest{Parent: projectName, Filter: filterStr}
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

// [END discovery_info]

// [START occurrences_for_note]

// getOccurrencesForNote retrieves all the Occurrences associated with a specified Note.
// Here, all Occurrences are printed and counted.
func getOccurrencesForNote(ctx context.Context, client *containeranalysis.GrafeasV1Beta1Client, noteID, projectID string) (int, error) {
	noteName := "projects/" + projectID + "/notes/" + noteID

	req := &grafeaspb.ListNoteOccurrencesRequest{Name: noteName}
	it := client.ListNoteOccurrences(ctx, req)
	count := 0
	for {
		_, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return -1, err
		}
		// Write custom code to process each Occurrence here.
		count = count + 1
	}
	return count, nil
}

// [END occurrences_for_note]

// [START occurrences_for_image]

// getOccurrencesForImage retrieves all the Occurrences associated with a specified image.
// Here, all Occurrences are simply printed and counted.
func getOccurrencesForImage(ctx context.Context, client *containeranalysis.GrafeasV1Beta1Client, imageUrl, projectID string) (int, error) {
	filterStr := `resourceUrl="` + imageUrl + `"`
	project := "projects/" + projectID

	req := &grafeaspb.ListOccurrencesRequest{Parent: project, Filter: filterStr}
	it := client.ListOccurrences(ctx, req)
	count := 0
	for {
		_, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return -1, err
		}
		// Write custom code to process each Occurrence here.
		count = count + 1
	}
	return count, nil
}

// [END occurrences_for_image]

// [START pubsub]

// occurrencePubsub handles incoming Occurrences using a Cloud Pub/Sub subscription.
func occurrencePubsub(ctx context.Context, subscriptionID string, timeout int, projectID string) (int, error) {
	toctx, _ := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	var mu sync.Mutex
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return -1, err
	}
	// Subscribe to the requested Pub/Sub channel.
	sub := client.Subscription(subscriptionID)
	count := 0

	// Listen to messages for 'timeout' seconds.
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
func createOccurrenceSubscription(ctx context.Context, subscriptionID, projectID string) error {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	// This topic id will automatically receive messages when Occurrences are added or modified
	topicID := "resource-notes-occurrences-v1alpha1"
	topic := client.Topic(topicID)
	config := pubsub.SubscriptionConfig{Topic: topic}
	_, err = client.CreateSubscription(ctx, subscriptionID, config)
	return err
}

// [END pubsub]
