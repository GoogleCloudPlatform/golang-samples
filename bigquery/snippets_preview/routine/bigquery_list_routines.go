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

package routine

// [START bigquery_list_routines_preview]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery/v2/apiv2/bigquerypb"
	"cloud.google.com/go/bigquery/v2/apiv2_client"

	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// listRoutines demonstrates iterating through the routines within a specified dataset.
func listRoutines(client *apiv2_client.Client, w io.Writer, projectID, datasetID string) error {
	// client can be instantiated per-RPC service, or use cloud.google.com/go/bigquery/v2/apiv2_client to create
	// an aggregate client.
	//
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	ctx := context.Background()

	req := &bigquerypb.ListRoutinesRequest{
		ProjectId: projectID,
		DatasetId: datasetID,
		// MaxResults is the per-page threshold (aka page size).  Generally you should only
		// worry about setting this if you're executing code in a memory constrained environment
		// and don't want to process large pages of results.  BigQuery will select a reasonable
		// page size automatically.
		MaxResults: &wrapperspb.UInt32Value{Value: 100},
	}

	// ListRoutines returns an iterator so users don't have to manage pagination when processing
	// the results.
	it := client.ListRoutines(ctx, req)

	// Process data from the iterator one result at a time.  The internal implementation of the iterator
	// is fetching pages at a time.
	for {
		routine, err := it.Next()
		if err == iterator.Done {
			// We're reached the end of the iteration, break the loop.
			break
		}
		if err != nil {
			return fmt.Errorf("iterator errored: %w", err)
		}
		// Print basic information to the provided writer.
		fmt.Fprintf(w, "routine %q reports type %q\n", routine.GetRoutineReference().GetRoutineId(), routine.GetRoutineType().String())
	}
	return nil
}

// [END bigquery_list_routines_preview]
