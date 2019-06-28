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

// [START containeranalysis_discovery_info]

import (
	"context"
	"fmt"
	"io"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1"
	"google.golang.org/api/iterator"
	grafeaspb "google.golang.org/genproto/googleapis/grafeas/v1"
)

// getDiscoveryInfo retrieves and prints the Discovery Occurrence created for a specified image.
// The Discovery Occurrence contains information about the initial scan on the image.
func getDiscoveryInfo(w io.Writer, resourceURL, projectID string) error {
	// resourceURL := fmt.Sprintf("https://gcr.io/my-project/my-image")
	ctx := context.Background()
	client, err := containeranalysis.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	req := &grafeaspb.ListOccurrencesRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		Filter: fmt.Sprintf(`kind="DISCOVERY" AND resourceUrl=%q`, resourceURL),
	}
	it := client.GetGrafeasClient().ListOccurrences(ctx, req)
	for {
		occ, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("occurrence iteration error: %v", err)
		}
		fmt.Fprintln(w, occ)
	}
	return nil
}

// [END containeranalysis_discovery_info]
