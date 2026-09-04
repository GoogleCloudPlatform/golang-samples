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

// [START containeranalysis_create_note]

import (
	"context"
	"fmt"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1"
	grafeaspb "google.golang.org/genproto/googleapis/grafeas/v1"
)

// createNote creates and returns a new vulnerability Note.
func createNote(noteID, projectID string) (*grafeaspb.Note, error) {
	ctx := context.Background()
	client, err := containeranalysis.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	projectName := fmt.Sprintf("projects/%s", projectID)

	req := &grafeaspb.CreateNoteRequest{
		Parent: projectName,
		NoteId: noteID,
		Note: &grafeaspb.Note{
			Type: &grafeaspb.Note_Attestation{
			// The "Attestation" field can be modified to contain information about your attestation
				Attestation: &grafeaspb.AttestationNote{
					Hint: &grafeaspb.AttestationNote_Hint{
						HumanReadableName: "my-attestation-authority",
					},
				},
			},
		},
	}

	return client.GetGrafeasClient().CreateNote(ctx, req)
}

// [END containeranalysis_create_note]
