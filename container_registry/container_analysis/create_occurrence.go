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

package main

// [START containeranalysis_create_occurrence]

import (
	"context"
	"fmt"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1"
	grafeaspb "google.golang.org/genproto/googleapis/grafeas/v1"
)

// createsOccurrence creates and returns a new Occurrence of a previously created vulnerability Note.
func createOccurrence(resourceURL, noteID, occProjectID, noteProjectID string) (*grafeaspb.Occurrence, error) {
	// resourceURL := fmt.Sprintf("https://gcr.io/my-project/my-image")
	// noteID := fmt.Sprintf("my-note")
	ctx := context.Background()
	client, err := containeranalysis.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	req := &grafeaspb.CreateOccurrenceRequest{
		Parent: fmt.Sprintf("projects/%s", occProjectID),
		Occurrence: &grafeaspb.Occurrence{
			NoteName: fmt.Sprintf("projects/%s/notes/%s", noteProjectID, noteID),
			// Attach the occurrence to the associated resource uri.
			ResourceUri: resourceURL,
			// Details about the vulnerability instance can be added here.
			Details: &grafeaspb.Occurrence_Vulnerability{
				Vulnerability: &grafeaspb.VulnerabilityOccurrence{
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
	return client.GetGrafeasClient().CreateOccurrence(ctx, req)
}

// [END containeranalysis_create_occurrence]
