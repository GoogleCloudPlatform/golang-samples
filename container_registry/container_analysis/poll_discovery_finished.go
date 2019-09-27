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

// [START containeranalysis_poll_discovery_occurrence_finished]

import (
	"context"
	"fmt"
	"time"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1"
	"google.golang.org/api/iterator"
	grafeaspb "google.golang.org/genproto/googleapis/grafeas/v1"
)

// pollDiscoveryOccurrenceFinished returns the discovery occurrence for a resource once it reaches a finished state.
func pollDiscoveryOccurrenceFinished(resourceURL, projectID string, timeout time.Duration) (*grafeaspb.Occurrence, error) {
	// resourceURL := fmt.Sprintf("https://gcr.io/my-project/my-image")
	// timeout := time.Duration(5) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := containeranalysis.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	// ticker is used to poll once per second.
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Find the discovery occurrence using a filter string.
	var discoveryOccurrence *grafeaspb.Occurrence
	for discoveryOccurrence == nil {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout while retrieving discovery occurrence")
		case <-ticker.C:
			req := &grafeaspb.ListOccurrencesRequest{
				Parent: fmt.Sprintf("projects/%s", projectID),
				// Vulnerability discovery occurrences are always associated with the
				// PACKAGE_VULNERABILITY note in the "goog-analysis" GCP project.
				Filter: fmt.Sprintf(`resourceUrl=%q AND noteProjectId="goog-analysis" AND noteId="PACKAGE_VULNERABILITY"`, resourceURL),
			}
			// [END containeranalysis_poll_discovery_occurrence_finished]
			// The above filter isn't testable, since it looks for occurrences in a locked down project.
			// Fall back to a more permissive filter for testing.
			req = &grafeaspb.ListOccurrencesRequest{
				Parent: fmt.Sprintf("projects/%s", projectID),
				Filter: fmt.Sprintf(`kind="DISCOVERY" AND resourceUrl=%q`, resourceURL),
			}
			// [START containeranalysis_poll_discovery_occurrence_finished]
			it := client.GetGrafeasClient().ListOccurrences(ctx, req)
			// Only one occurrence should ever be returned by ListOccurrences
			// and the given filter.
			result, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("it.Next: %v", err)
			}
			if result.GetDiscovery() != nil {
				discoveryOccurrence = result
			}
		}
	}

	// Wait for the discovery occurrence to enter a terminal state.
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for terminal state")
		case <-ticker.C:
			// Update the occurrence.
			req := &grafeaspb.GetOccurrenceRequest{Name: discoveryOccurrence.GetName()}
			updated, err := client.GetGrafeasClient().GetOccurrence(ctx, req)
			if err != nil {
				return nil, fmt.Errorf("GetOccurrence: %v", err)
			}
			switch updated.GetDiscovery().GetAnalysisStatus() {
			case grafeaspb.DiscoveryOccurrence_FINISHED_SUCCESS,
				grafeaspb.DiscoveryOccurrence_FINISHED_FAILED,
				grafeaspb.DiscoveryOccurrence_FINISHED_UNSUPPORTED:
				return discoveryOccurrence, nil
			}
		}
	}
}

// [END containeranalysis_poll_discovery_occurrence_finished]
