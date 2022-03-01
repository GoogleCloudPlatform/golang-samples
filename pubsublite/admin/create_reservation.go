// Copyright 2022 Google LLC
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

package admin

// [START pubsublite_create_reservation]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/pubsublite"
)

func createReservation(w io.Writer, projectID, region, reservationID string, throughputCapacity int) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// reservationID := "my-reservation"
	// throughputCapacity := 4
	ctx := context.Background()
	client, err := pubsublite.NewAdminClient(ctx, region)
	if err != nil {
		return fmt.Errorf("pubsublite.NewAdminClient: %v", err)
	}
	defer client.Close()

	reservationPath := fmt.Sprintf("projects/%s/locations/%s/reservations/%s", projectID, region, reservationID)
	res, err := client.CreateReservation(ctx, pubsublite.ReservationConfig{
		Name:               reservationPath,
		ThroughputCapacity: throughputCapacity,
	})
	if err != nil {
		return fmt.Errorf("client.CreateReservation got err: %v", err)
	}
	fmt.Fprintf(w, "Created reservation: %s\n", res.Name)
	return nil
}

// [END pubsublite_create_reservation]
