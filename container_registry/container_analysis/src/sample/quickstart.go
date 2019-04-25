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

// Sample quickstart for getting vulnerabilities from the Container Analysis API: https://cloud.google.com/container-registry/docs/vulnerability-scan-go
package sample

// [START containeranalysis_imports_quickstart]

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

// [END containeranalysis_imports_quickstart]

// [START containeranalysis_poll_discovery_occurrence_finished]

// pollDiscoveryOccurrenceFinished returns the discovery occurrence for a resource once it reaches a finished state.
func pollDiscoveryOccurrenceFinished(resourceURL, projectID string, timeout time.Duration) (*grafeaspb.Occurrence, error) {
	// resourceURL := fmt.Sprintf("https://gcr.io/my-project/my-image")
	// timeout := time.Duration(5) * time.Second
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return nil, err
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
			return false, err
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

// [START containeranalysis_vulnerability_occurrences_for_image]

// findVulnerabilityOccurrencesForImage retrieves all vulnerability Occurrences associated with a resource.
func findVulnerabilityOccurrencesForImage(resourceURL, projectID string) ([]*grafeaspb.Occurrence, error) {
	// resourceURL := fmt.Sprintf("https://gcr.io/my-project/my-image")
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &grafeaspb.ListOccurrencesRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		Filter: fmt.Sprintf("resourceUrl = %q kind = %q", resourceURL, "VULNERABILITY"),
	}

	var occurrenceList []*grafeaspb.Occurrence
	it := client.ListOccurrences(ctx, req)
	for {
		occ, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		occurrenceList = append(occurrenceList, occ)
	}

	return occurrenceList, nil
}

// [END containeranalysis_vulnerability_occurrences_for_image]

// [START containeranalysis_filter_vulnerability_occurrences]

// findHighSeverityVulnerabilitiesForImage retrieves a list of only high vulnerability occurrences associated with a resource.
func findHighSeverityVulnerabilitiesForImage(resourceURL, projectID string) ([]*grafeaspb.Occurrence, error) {
	// resourceURL := fmt.Sprintf("https://gcr.io/my-project/my-image")
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	req := &grafeaspb.ListOccurrencesRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		Filter: fmt.Sprintf("resourceUrl = %q kind = %q", resourceURL, "VULNERABILITY"),
	}

	var occurrenceList []*grafeaspb.Occurrence
	it := client.ListOccurrences(ctx, req)
	for {
		occ, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		severityLevel := occ.GetVulnerability().GetSeverity()
		if severityLevel == vulnerability.Severity_HIGH || severityLevel == vulnerability.Severity_CRITICAL {
			occurrenceList = append(occurrenceList, occ)
		}
	}

	return occurrenceList, nil
}

// [END containeranalysis_filter_vulnerability_occurrences]
