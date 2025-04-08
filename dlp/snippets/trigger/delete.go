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

package trigger

// [START dlp_delete_trigger]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
)

// deleteTrigger deletes the given trigger.
func deleteTrigger(w io.Writer, triggerID string) error {
	// triggerID := "my-trigger"

	ctx := context.Background()

	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %w", err)
	}
	defer client.Close()

	req := &dlppb.DeleteJobTriggerRequest{
		Name: triggerID,
	}

	if err := client.DeleteJobTrigger(ctx, req); err != nil {
		return fmt.Errorf("DeleteJobTrigger: %w", err)
	}
	fmt.Fprintf(w, "Successfully deleted trigger %v", triggerID)
	return nil
}

// [END dlp_delete_trigger]
