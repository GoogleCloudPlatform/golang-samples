// Copyright 2025 Google LLC
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

package dataset

// [START bigquery_list_datasets_preview]
import (
	"context"
	"fmt"
	"io"

	bigquery "cloud.google.com/go/bigquery/apiv2"
	"cloud.google.com/go/bigquery/apiv2/bigquerypb"

	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// listDatasets demonstrates iterating through datasets.
func listDatasets(w io.Writer, projectID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	ctx := context.Background()

	// Construct a gRPC-based client.
	// To construct a REST-based client, use NewDatasetRESTClient instead.
	dsClient, err := bigquery.NewDatasetClient(ctx)
	if err != nil {
		return fmt.Errorf("bigquery.NewDatasetClient: %w", err)
	}
	defer dsClient.Close()

	req := &bigquerypb.ListDatasetsRequest{
		ProjectId: projectID,
		// MaxResults is the per-page threshold (aka page size).  Generally you should only
		// worry about setting this if you're executing code in a memory constrained environment
		// and don't want to process large pages of results.
		MaxResults: &wrapperspb.UInt32Value{Value: 100},
	}

	// ListDatasets returns an iterator so users don't have to manage pagination when processing
	// the results.
	it := dsClient.ListDatasets(ctx, req)

	// Process data from the iterator one result at a time.  The internal implementation of the iterator
	// is fetching pages at a time.
	for {
		dataset, err := it.Next()
		if err == iterator.Done {
			// We're reached the end of the iteration, break the loop.
			break
		}
		if err != nil {
			return fmt.Errorf("iterator errored: %w", err)
		}
		// Print basic information to the provided writer.
		fmt.Fprintf(w, "dataset %q in location %q\n", dataset.GetDatasetReference().GetDatasetId(), dataset.GetLocation())
	}
	return nil
}

// [END bigquery_list_datasets_preview]
