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

import (
	"context"
	"log"
	"fmt"
	"time"
	wait "k8s.io/apimachinery/pkg/util/wait"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1beta1"
	"google.golang.org/api/iterator"
	discovery "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/discovery"
	grafeaspb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/grafeas"
	vulnerability "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/vulnerability"
)

const (
	discoveryProviderID          = "goog-analysis"
	discoveryVulnerabilityNoteID = "PACKAGE_VULNERABILITY"
	backoff                      = time.Second
)

func main() {
	projectID := "my-project-id"                                           // The Google Cloud Platform project ID
	resourceURL := "https://gcr.io/my-project-id/image-name@sha256:digest" // The resource URL to analyze
	
	ctx := context.Background()

	// Get a Container Analysis Client.
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Close client when done.
	defer client.Close()

	// Wait for analysis to complete.
	occ, err := pollDiscoveryOccurrenceFinished(resourceURL, projectID, time.Duration(30)*time.Second)
	if err != nil {
		log.Fatalf("Failed to get occurrence: %v", err)
	}
	analysisStatus := occ.GetDiscovered().GetDiscovered().AnalysisStatus
	log.Printf("Vulnerability analysis complete, status is %s", analysisStatus)


	// Find any vulnerabilities for this resource.
	vulnOccs, err := findVulnerabilityOccurrencesForImage(resourceURL, projectID)
	if err != nil {
		log.Fatalf("Failed to get vulnerability occurrences: %v", err)
	}
	log.Printf("Got %d vulnerabilities", len(vulnOccs))

	// find high severity occurrences
	filteredOccs := findHighSeverityVulnerabilitiesForImage(resourceURL, projectID)

	log.Printf("Got %d severity high or above vulnerabilities:", len(filteredOccs))
	for i, occ := range filteredOccs {
		log.Printf("Vulnerability %d:\n%+v", i+1, occ)
	}

}

// [START containeranalysis_poll_discovery_occurrence_finished]

// pollDiscoveryOccurrenceFinished returns a discovery occurrence for an image once that discovery occurrence is in a finished state.
func pollDiscoveryOccurrenceFinished(resourceURL, projectID string, timeout time.Duration) (*grafeaspb.Occurrence, error) {
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return nil, err
	}

	// Discover Occurrences are stored within a Google-owned project. 
	const providerProjectID = "goog-analysis"
	const providerNoteID = "PACKAGE_VULNERABILITY"

	// find the discovery occurrence
	var discoveryOccurrence *grafeaspb.Occurrence
	err = wait.Poll(time.Second, timeout, func() (bool, error) {
		log.Printf("Querying for discovery occurrence")
		req := &grafeaspb.ListOccurrencesRequest{
			Parent: fmt.Sprintf("projects/%s", projectID),
			Filter: fmt.Sprintf("resourceUrl = %q noteProjectId = %q noteId = %q", resourceURL, providerProjectID, providerNoteID),
		}
		it := client.ListOccurrences(ctx, req)
		// Only one should ever be returned by ListOccurrences and the given filter.
		result, err := it.Next()
		if err != nil || result == nil || result.GetDiscovered() == nil {
			return false, nil
		} else {
			discoveryOccurrence = result
			return true, nil
		}
	})
	if err != nil {
		return nil, fmt.Errorf("could not find dicovery occurrence: %v", err)
	}

	// wait for terminal state
	err = wait.Poll(time.Second, timeout, func() (bool, error) {
		// check for updated occurrence state
		newOccurrence, err := client.GetOccurrence(ctx, &grafeaspb.GetOccurrenceRequest{Name: discoveryOccurrence.GetName()})
		if err != nil {
			return false, err
		} else {
			discoveryOccurrence = newOccurrence
		}
		// check if in ternimal state
		state := discoveryOccurrence.GetDiscovered().GetDiscovered().GetAnalysisStatus()
		isTerminal := (state ==  discovery.Discovered_FINISHED_SUCCESS ||
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

// findVulnerabilityOccurrencesForImage retrieves all vulnerability Occurrences associated with an image.
func findVulnerabilityOccurrencesForImage(resourceURL, projectID string) ([]*grafeaspb.Occurrence, error) {
	ctx := context.Background()
	client, err := containeranalysis.NewGrafeasV1Beta1Client(ctx)
	if err != nil {
		return nil, err
	}
	
	var occs []*grafeaspb.Occurrence

	req := &grafeaspb.ListOccurrencesRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		Filter: fmt.Sprintf("resourceUrl = %q kind = %q", resourceURL, "VULNERABILITY"),
	}

	it := client.ListOccurrences(ctx, req)
	for {
		occ, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		occs = append(occs, occ)
	}

	return occs, nil
}

// [END containeranalysis_vulnerability_occurrences_for_image]


// [START containeranalysis_filter_vulnerability_occurrences]

func findHighSeverityVulnerabilitiesForImage(resourceUrl, projectID string) error {
	// retrieve a list of all vulnerabilities using the function defined above
	vulnOccs, err := findVulnerabilityOccurrencesForImage(resourceURL, projectID)
	if err != nil {
		return fmt.Errorf("Failed to get vulnerability occurrences: %v", err)
	}
	// add high severity occurrences to a new filtered list
	var filteredOccs []*grafeaspb.Occurrence
	for _, occ := range vulnOccs {
		severityLevel = occ.GetVulnerability().GetSeverity()
		if severityLevel == vulnerability.Severity_HIGH || severityLevel == vulnerability.Severity_CRITICAL {
			filteredOccs = append(filteredOccs, occ)
		}
	}
	return filteredOccs
}

// [END containeranalysis_filter_vulnerability_occurrences]
