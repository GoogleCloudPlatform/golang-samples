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

// [START bigquery_create_routine_preview]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/bigquery/v2/apiv2/bigquerypb"
	"cloud.google.com/go/bigquery/v2/apiv2_client"

	"github.com/googleapis/gax-go/v2/apierror"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
)

// createRoutine demonstrates creation of a new routine that represents a SQL User Defined Function (UDF).
func createRoutine(client *apiv2_client.Client, w io.Writer, projectID, datasetID, routineID string) error {
	// client can be instantiated per-RPC service, or use cloud.google.com/go/bigquery/v2/apiv2_client to create
	// an aggregate client.
	//
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// routineID := "myroutine"
	ctx := context.Background()

	// Construct a request, populating some of the available configuration
	// settings.
	req := &bigquerypb.InsertRoutineRequest{
		ProjectId: projectID,
		DatasetId: datasetID,
		Routine: &bigquerypb.Routine{
			RoutineReference: &bigquerypb.RoutineReference{
				ProjectId: projectID,
				DatasetId: datasetID,
				RoutineId: routineID,
			},
			RoutineType: bigquerypb.Routine_SCALAR_FUNCTION,
			Arguments: []*bigquerypb.Routine_Argument{
				{
					Name: "x",
					DataType: &bigquerypb.StandardSqlDataType{
						TypeKind: bigquerypb.StandardSqlDataType_INT64,
					},
				},
			},
			ReturnType: &bigquerypb.StandardSqlDataType{
				TypeKind: bigquerypb.StandardSqlDataType_INT64,
			},
			DefinitionBody: "x * 3",
		},
	}
	resp, err := client.InsertRoutine(ctx, req)
	if err != nil {
		// Examine the error structure more deeply.
		if apierr, ok := apierror.FromError(err); ok {
			if status := apierr.GRPCStatus(); status.Code() == codes.AlreadyExists {
				// The error was due to the routine already existing.  For this sample
				// we don't consider that a failure, so return nil.
				return nil
			}
		}
		return fmt.Errorf("InsertRoutine: %w", err)
	}
	// Print the JSON representation of the response to the provided writer.
	fmt.Fprintf(w, "Response from insert: %s", protojson.Format(resp))
	return nil
}

// [END bigquery_create_routine_preview]
