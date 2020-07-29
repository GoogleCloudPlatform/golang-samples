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

// [START fs_detach_listener]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
)

// listenStop demonstrates how to detach a listener.
func listenStop(w io.Writer, projectID string) error {
	// projectID := "project-id"
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	qsnap := client.Collection("cities").Where("state", "==", "CA").Snapshots(ctx)
	// Stop listening for changes
	qsnap.Stop()

	return nil
}

// [END fs_detach_listener]
