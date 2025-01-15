// Copyright 2022 Google LLC
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

package client

// [START bigquerydatatransfer_create_scheduled_query]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	datatransfer "cloud.google.com/go/bigquery/datatransfer/apiv1"
	"cloud.google.com/go/bigquery/datatransfer/apiv1/datatransferpb"
	"google.golang.org/protobuf/types/known/structpb"
)

// createScheduledQuery schedule a query to run every
// 24 hours with a destination table identifier based on the run date.
func createScheduledQuery(projectID, datasetID, query string) error {
	// projectID := "my-project-id"
	// datasetID := "my-dataset-id"
	// query = `SELECT CURRENT_TIMESTAMP() as current_time, @run_time as intended_run_time,
	// @run_date as intended_run_date, 17 as some_integer`
	ctx := context.Background()
	dtc, err := datatransfer.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("datatransfer.NewClient: %w", err)
	}
	defer dtc.Close()

	paramsMap := map[string]interface{}{}
	paramsMap["destination_table_name_template"] = "my_destination_table_{run_date}"
	paramsMap["write_disposition"] = string(bigquery.WriteTruncate)
	paramsMap["partitioning_field"] = ""
	paramsMap["query"] = query
	params, err := structpb.NewStruct(paramsMap)
	if err != nil {
		return fmt.Errorf("structpb.NewStruct: %w", err)
	}

	req := &datatransferpb.CreateTransferConfigRequest{
		Parent: fmt.Sprintf("projects/%s", projectID),
		TransferConfig: &datatransferpb.TransferConfig{
			Destination: &datatransferpb.TransferConfig_DestinationDatasetId{
				DestinationDatasetId: datasetID,
			},
			DisplayName:  "Your Scheduled Query Name",
			DataSourceId: "scheduled_query",
			Params:       params,
			Schedule:     "every 24 hours",
		},
	}

	_, err = dtc.CreateTransferConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("dtc.CreateTransferConfig: %w", err)
	}

	return nil
}

// [END bigquerydatatransfer_create_scheduled_query]
