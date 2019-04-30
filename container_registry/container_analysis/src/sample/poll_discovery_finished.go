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

// Samples for the Container Analysis golang libraries: https://cloud.google.com/container-registry/docs/container-analysis
package sample

// [START containeranalysis_poll_discovery_occurrence_finished]

import (
	"context"
	"fmt"
	wait "k8s.io/apimachinery/pkg/util/wait"
	"log"
	"time"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1beta1"
	"google.golang.org/api/iterator"
	discovery "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/discovery"
	grafeaspb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/grafeas"
	vulnerability "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/vulnerability"
)

// pollDiscoveryOccurrenceFinished returns the discovery occurrence for a resource once it reaches a finished state.
func pollDiscoveryOccurrenceFinished(resourceURL, projectID string, timeout time.Duration) (*grafeaspb.Occurrence, error) {
	// resourceURL := fmt.Sprintf("https://gcr.io/my-project/my-image")
	// timeout := time.Duration(5) * time.Second
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewGrafeasV1Beta1Client: %v", err)
	}
	defer client.Close()

	// Find the discovery occurrence using a filter string.
	var discoveryOccurrence *grafeaspb.Occurrence
	err = wait.Poll(time.Second, timeout, func() (bool, error) {
		log.Printf("Querying for discovery occurrence")
		req := &grafeaspb.ListOccurrencesRequest{
			Parent: fmt.Sprintf("projects/%s", projectID),
			Filter: fmt.Sprintf(`kind="DISCOVERY" AND resourceUrl=%q`, resourceURL),
		}
		it := client.ListOccurrences(ctx, req)
		// Only one should ever be returned by ListOccurrences and the given filter.
		result, err := it.Next()
		if err != nil || result == nil || result.GetDiscovered() == nil {
			return false, nil
		}
		discoveryOccurrence = result
		return true, nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not find dicovery occurrence: %v", err)
	}

	// Wait for the discovery occurrence to enter a terminal state.
	err = wait.Poll(time.Second, timeout, func() (bool, error) {
		// Update the occurrence
		req := &grafeaspb.GetOccurrenceRequest{Name: discoveryOccurrence.GetName()}
		newOccurrence, err := client.GetOccurrence(ctx, req)
		if err != nil {
			return false, fmt.Errorf("GetOccurrence: %v", err)
		}
		discoveryOccurrence = newOccurrence
		// Check if the discovery occurrence is in a ternimal state.
		state := discoveryOccurrence.GetDiscovered().GetDiscovered().GetAnalysisStatus()
		isTerminal := (state == discovery.Discovered_FINISHED_SUCCESS ||
			state == discovery.Discovered_FINISHED_FAILED ||
			state == discovery.Discovered_FINISHED_UNSUPPORTED)
		return isTerminal, nil
	})
	if err != nil {
		return nil, fmt.Errorf("occurrence never reached terminal state: %v", err)
	}
	return discoveryOccurrence, nil
}

// [END containeranalysis_poll_discovery_occurrence_finished]
