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

// [START bigquerydatatransfer_delete_scheduled_query]
import (
	"context"
	"fmt"

	datatransfer "cloud.google.com/go/bigquery/datatransfer/apiv1"
	"cloud.google.com/go/bigquery/datatransfer/apiv1/datatransferpb"
)

// deleteScheduledQuery delete a scheduled query based on
// the config ID, stopping any future runs.
// transferConfigID follows the format:
//
//	`projects/{project_id}/locations/{location_id}/transferConfigs/{config_id}`
//	or `projects/{project_id}/transferConfigs/{config_id}`
func deleteScheduledQuery(transferConfigID string) error {
	// transferConfigID := "projects/{project_id}/transferConfigs/{config_id}"
	ctx := context.Background()
	dtc, err := datatransfer.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("datatransfer.NewClient: %w", err)
	}
	defer dtc.Close()

	req := &datatransferpb.DeleteTransferConfigRequest{
		Name: transferConfigID,
	}
	err = dtc.DeleteTransferConfig(ctx, req)
	if err != nil {
		return fmt.Errorf("dtc.DeleteTransferConfig: %w", err)
	}

	return nil
}

// [END bigquerydatatransfer_delete_scheduled_query]
