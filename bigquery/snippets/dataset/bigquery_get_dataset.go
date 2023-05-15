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

package dataset

// [START bigquery_get_dataset]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// printDatasetInfo demonstrates fetching dataset metadata and printing some of it to an io.Writer.
func printDatasetInfo(w io.Writer, projectID, datasetID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	meta, err := client.Dataset(datasetID).Metadata(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Dataset ID: %s\n", datasetID)
	fmt.Fprintf(w, "Description: %s\n", meta.Description)
	fmt.Fprintln(w, "Labels:")
	for k, v := range meta.Labels {
		fmt.Fprintf(w, "\t%s: %s", k, v)
	}
	fmt.Fprintln(w, "Tables:")
	it := client.Dataset(datasetID).Tables(ctx)

	cnt := 0
	for {
		t, err := it.Next()
		if err == iterator.Done {
			break
		}
		cnt++
		fmt.Fprintf(w, "\t%s\n", t.TableID)
	}
	if cnt == 0 {
		fmt.Fprintln(w, "\tThis dataset does not contain any tables.")
	}
	return nil
}

// [END bigquery_get_dataset]
