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

// [START bigquery_get_routine_preview]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery/v2/apiv2/bigquerypb"
	"cloud.google.com/go/bigquery/v2/apiv2_client"
	"google.golang.org/protobuf/encoding/protojson"
)

// getRoutine demonstrates fetching routine information.
func getRoutine(client *apiv2_client.Client, w io.Writer, projectID, datasetID, routineID string) error {
	// client can be instantiated per-RPC service, or use cloud.google.com/go/bigquery/v2/apiv2_client to create
	// an aggregate client.
	//
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	ctx := context.Background()

	req := &bigquerypb.GetRoutineRequest{
		ProjectId: projectID,
		DatasetId: datasetID,
		RoutineId: routineID,
	}

	resp, err := client.GetRoutine(ctx, req)
	if err != nil {
		return fmt.Errorf("GetDataset: %w", err)
	}

	// Print some of the information about the routine to the provided writer.
	fmt.Fprintf(w, "Routine %q has description %q\n",
		resp.GetRoutineReference().GetRoutineId(),
		resp.GetDescription())
	// Alternately, use the protojson package to print a more complete representation
	// of the routine using a basic JSON mapping:
	fmt.Fprintf(w, "Routine JSON representation:\n%s\n", protojson.Format(resp))
	return nil
}

// [END bigquery_get_routine_preview]
