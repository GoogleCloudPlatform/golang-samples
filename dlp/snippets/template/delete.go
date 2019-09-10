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

package template

// [START dlp_delete_inspect_template]
import (
	"context"
	"fmt"
	"io"

	dlp "cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// deleteInspectTemplate deletes the given template.
func deleteInspectTemplate(w io.Writer, templateID string) error {
	// projectID := "my-project-id"
	// templateID := "my-template"

	ctx := context.Background()

	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %v", err)
	}

	req := &dlppb.DeleteInspectTemplateRequest{
		Name: templateID,
	}

	if err := client.DeleteInspectTemplate(ctx, req); err != nil {
		return fmt.Errorf("DeleteInspectTemplate: %v", err)
	}
	fmt.Fprintf(w, "Successfully deleted inspect template %v", templateID)
	return nil
}

// [END dlp_delete_inspect_template]
