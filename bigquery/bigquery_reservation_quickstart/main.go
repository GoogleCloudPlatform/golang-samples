// Copyright 2020 Google LLC
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

// [START bigqueryreservation_quickstart]

// The bigquery_reservation_quickstart application demonstrates usage of the
// BigQuery reservation API by enumerating some of the resources that can be
// associated with a cloud project.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"

	reservation "cloud.google.com/go/bigquery/reservation/apiv1"
	"cloud.google.com/go/bigquery/reservation/apiv1/reservationpb"
	"google.golang.org/api/iterator"
)

func main() {

	// Define two command line flags for controlling the behavior of this quickstart.
	var (
		projectID = flag.String("project_id", "", "Cloud Project ID, used for session creation.")
		location  = flag.String("location", "US", "BigQuery location used for interactions.")
	)
	// Parse flags and do some minimal validation.
	flag.Parse()
	if *projectID == "" {
		log.Fatal("empty --project_id specified, please provide a valid project ID")
	}
	if *location == "" {
		log.Fatal("empty --location specified, please provide a valid location")
	}

	ctx := context.Background()
	bqResClient, err := reservation.NewClient(ctx)
	if err != nil {
		log.Fatalf("NewClient: %v", err)
	}
	defer bqResClient.Close()

	s, err := reportCapacityCommitments(ctx, bqResClient, *projectID, *location)
	if err != nil {
		log.Fatalf("printCapacityCommitments: %v", err)
	}
	fmt.Println(s)

	s, err = reportReservations(ctx, bqResClient, *projectID, *location)
	if err != nil {
		log.Fatalf("printReservations: %v", err)
	}
	fmt.Println(s)
}

// printCapacityCommitments iterates through the capacity commitments and returns a byte buffer with details.
func reportCapacityCommitments(ctx context.Context, client *reservation.Client, projectID, location string) (string, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Capacity commitments in project %s in location %s:\n", projectID, location)

	req := &reservationpb.ListCapacityCommitmentsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}
	totalCommitments := 0
	it := client.ListCapacityCommitments(ctx, req)
	for {
		commitment, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&buf, "\tCommitment %s in state %s\n", commitment.GetName(), commitment.GetState().String())
		totalCommitments++
	}
	fmt.Fprintf(&buf, "\n%d commitments processed.\n", totalCommitments)
	return buf.String(), nil
}

// printReservations iterates through reservations defined in an admin project.
func reportReservations(ctx context.Context, client *reservation.Client, projectID, location string) (string, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Reservations in project %s in location %s:\n", projectID, location)

	req := &reservationpb.ListReservationsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}
	totalReservations := 0
	it := client.ListReservations(ctx, req)
	for {
		reservation, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&buf, "\tReservation %s has %d slot capacity.\n", reservation.GetName(), reservation.GetSlotCapacity())
		totalReservations++
	}
	fmt.Fprintf(&buf, "\n%d reservations processed.\n", totalReservations)
	return buf.String(), nil
}

// [END bigqueryreservation_quickstart]
