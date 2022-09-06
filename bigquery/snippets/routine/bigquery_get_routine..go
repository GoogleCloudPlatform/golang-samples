// Copyright 2021 Google LLC
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

// [START bigquery_get_routine]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
)

// getRoutine demonstrates getting a routine's metadata via the API.
func getRoutine(w io.Writer, projectID, datasetID, routineID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydatasetid"
	// routineID := "myroutineid"
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	meta, err := client.Dataset(datasetID).Routine(routineID).Metadata(ctx)
	if err != nil {
		return fmt.Errorf("couldn't retrieve routine metadata: %w", err)
	}
	// Print information about the routine.
	fmt.Fprintf(w, "Routine %s:\n", routineID)
	fmt.Fprintf(w, "\tType %s:\n", meta.Type)
	fmt.Fprintf(w, "\tLanguage %s:\n", meta.Language)
	fmt.Fprintln(w, "\tArguments:")
	for _, v := range meta.Arguments {
		fmt.Fprintf(w, "\t\tName: %s\tType: %v", v.Name, v.DataType)
	}
	return nil
}

// [END bigquery_get_routine]
