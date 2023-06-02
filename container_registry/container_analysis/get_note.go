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

package containeranalysis

// [START containeranalysis_get_note]

import (
	"context"
	"fmt"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1"
	grafeaspb "google.golang.org/genproto/googleapis/grafeas/v1"
)

// getNote retrieves and prints a specified Note from the server.
func getNote(noteID, projectID string) (*grafeaspb.Note, error) {
	// noteID := fmt.Sprintf("my-note")
	ctx := context.Background()
	client, err := containeranalysis.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &grafeaspb.GetNoteRequest{
		Name: fmt.Sprintf("projects/%s/notes/%s", projectID, noteID),
	}
	note, err := client.GetGrafeasClient().GetNote(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("client.GetNote: %w", err)
	}
	return note, nil
}

// [END containeranalysis_get_note]
