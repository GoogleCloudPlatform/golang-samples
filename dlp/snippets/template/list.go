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

// [START dlp_list_inspect_templates]
import (
	"context"
	"fmt"
	"io"
	"time"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/dlp/apiv2/dlppb"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/api/iterator"
)

// listInspectTemplates lists the inspect templates in the project.
func listInspectTemplates(w io.Writer, projectID string) error {
	// projectID := "my-project-id"

	ctx := context.Background()

	client, err := dlp.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("dlp.NewClient: %w", err)
	}
	defer client.Close()

	// Create a configured request.
	req := &dlppb.ListInspectTemplatesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
	}

	// Send the request and iterate over the results.
	it := client.ListInspectTemplates(ctx, req)
	for {
		t, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Next: %w", err)
		}
		fmt.Fprintf(w, "Inspect template %v\n", t.GetName())
		c, err := ptypes.Timestamp(t.GetCreateTime())
		if err != nil {
			return fmt.Errorf("CreateTime Timestamp: %w", err)
		}
		fmt.Fprintf(w, "  Created: %v\n", c.Format(time.RFC1123))
		u, err := ptypes.Timestamp(t.GetUpdateTime())
		if err != nil {
			return fmt.Errorf("UpdateTime Timestamp: %w", err)
		}
		fmt.Fprintf(w, "  Updated: %v\n", u.Format(time.RFC1123))
		fmt.Fprintf(w, "  Display Name: %q\n", t.GetDisplayName())
		fmt.Fprintf(w, "  Description: %q\n", t.GetDescription())
	}

	return nil
}

// [END dlp_list_inspect_templates]
