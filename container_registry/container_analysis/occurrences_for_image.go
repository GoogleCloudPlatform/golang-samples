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

// [START containeranalysis_occurrences_for_image]

import (
	"context"
	"fmt"
	"io"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1"
	"google.golang.org/api/iterator"
	grafeaspb "google.golang.org/genproto/googleapis/grafeas/v1"
)

// getOccurrencesForImage retrieves all the Occurrences associated with a specified image.
// Here, all Occurrences are simply printed and counted.
func getOccurrencesForImage(w io.Writer, resourceURL, projectID string) (int, error) {
	// Use this style of URL when you use Google Container Registry.
	// resourceURL := "https://gcr.io/my-project/my-repo/my-image"
	// Use this style of URL when you use Google Artifact Registry.
	// resourceURL := "https://LOCATION-docker.pkg.dev/my-project/my-repo/my-image"
	ctx := context.Background()
	client, err := containeranalysis.NewClient(ctx)
	if err != nil {
		return -1, fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &grafeaspb.ListOccurrencesRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		Filter: fmt.Sprintf("resourceUrl=%q", resourceURL),
	}
	it := client.GetGrafeasClient().ListOccurrences(ctx, req)
	count := 0
	for {
		occ, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return -1, fmt.Errorf("occurrence iteration error: %w", err)
		}
		// Write custom code to process each Occurrence here.
		fmt.Fprintln(w, occ)
		count = count + 1
	}
	return count, nil
}

// [END containeranalysis_occurrences_for_image]
