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

// Package automl contains samples for Google Cloud AutoML API v1.
package automl

// [START automl_list_operation_status]
import (
	"context"
	"fmt"
	"io"

	automl "cloud.google.com/go/automl/apiv1"
	longrunning "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"google.golang.org/api/iterator"
)

// listOperationStatus lists existing operations' status.
func listOperationStatus(w io.Writer, projectID string, location string) error {
	// projectID := "my-project-id"
	// location := "us-central1"

	ctx := context.Background()
	client, err := automl.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &longrunning.ListOperationsRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}

	it := client.LROClient.ListOperations(ctx, req)

	// Iterate over all results
	for {
		op, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListOperations.Next: %w", err)
		}

		fmt.Fprintf(w, "Name: %v\n", op.GetName())
		fmt.Fprintf(w, "Operation details:\n")
		fmt.Fprintf(w, "%v", op)
	}

	return nil
}

// [END automl_list_operation_status]
