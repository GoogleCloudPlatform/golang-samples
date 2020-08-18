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

package main

// [START fs_listen_document]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
)

// listenDocument listens to a single document.
func listenDocument(w io.Writer, projectID string) error {
	// projectID := "project-id"
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	dsnap := client.Collection("cities").Doc("SF").Snapshots(ctx)

	snap, err := dsnap.Next()
	if !snap.Exists() {
		return fmt.Errorf("current data: null")
	}
	if err != nil {
		return fmt.Errorf("listen failed: %v", err)
	}
	fmt.Fprintf(w, "Received document snapshot: %v\n", snap.Data())
	return nil
}

// [END fs_listen_document]
