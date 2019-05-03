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

package sample

// [START containeranalysis_poll_discovery_occurrence_finished]

import (
	"context"
	"fmt"
	"log"
	"time"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1beta1"
	discovery "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/discovery"
	grafeaspb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/grafeas"
)

// pollDiscoveryOccurrenceFinished returns the discovery occurrence for a resource once it reaches a finished state.
func pollDiscoveryOccurrenceFinished(resourceURL, projectID string, timeout time.Duration) (*grafeaspb.Occurrence, error) {
	// resourceURL := fmt.Sprintf("https://gcr.io/my-project/my-image")
	// timeout := time.Duration(5) * time.Second
	deadline := time.Now().Add(timeout).Unix()
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewGrafeasV1Beta1Client: %v", err)
	}
	defer client.Close()

	// Find the discovery occurrence using a filter string.
	var discoveryOccurrence *grafeaspb.Occurrence
	for discoveryOccurrence == nil {
		log.Printf("Querying for discovery occurrence")
		req := &grafeaspb.ListOccurrencesRequest{
			Parent: fmt.Sprintf("projects/%s", projectID),
			Filter: fmt.Sprintf(`kind="DISCOVERY" AND resourceUrl=%q`, resourceURL),
		}
		it := client.ListOccurrences(ctx, req)
		// Only one should ever be returned by ListOccurrences and the given filter.
		result, err := it.Next()
		if result != nil && result.GetDiscovered() != nil {
			discoveryOccurrence = result;
		} else if time.Now().Unix() > deadline {
			return nil, fmt.Errorf("timeout while retrieving discovery occurrence: %v", err)
		} else {
			time.Sleep(time.Second)
		}
	}

	// Wait for the discovery occurrence to enter a terminal state.
	status := discovery.Discovered_PENDING
	for status != discovery.Discovered_FINISHED_SUCCESS &&
		status != discovery.Discovered_FINISHED_FAILED &&
		status != discovery.Discovered_FINISHED_UNSUPPORTED {
		// Update the occurrence.
		req := &grafeaspb.GetOccurrenceRequest{Name: discoveryOccurrence.GetName()}
		updated, err := client.GetOccurrence(ctx, req)
		if err == nil {
			// Update the analysis status object.
			status = updated.GetDiscovered().GetDiscovered().GetAnalysisStatus()
		} else if time.Now().Unix() > deadline {
			return nil, fmt.Errorf("timeout while waiting for terminal state: %v", err)
		} else {
			time.Sleep(time.Second)
		}
	}
	return discoveryOccurrence, nil
}

// [END containeranalysis_poll_discovery_occurrence_finished]
