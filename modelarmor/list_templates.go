// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     [https://www.apache.org/licenses/LICENSE-2.0](https://www.apache.org/licenses/LICENSE-2.0)
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Sample code for getting list of model armor templates.

package modelarmor

// [START modelarmor_list_templates]

import (
	"context"
	"fmt"
	"io"

	modelarmor "cloud.google.com/go/modelarmor/apiv1"
	modelarmorpb "cloud.google.com/go/modelarmor/apiv1/modelarmorpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// listModelArmorTemplates lists Model Armor templates.
//
// This method lists Model Armor templates for a project and location.
//
// Args:
//
//	w io.Writer: The writer to use for logging.
//	projectID string: The ID of the project.
//	location string: The location of the templates.
//
// Returns:
//
//	[]*modelarmorpb.Template: A list of Model Armor templates.
//	error: Any error that occurred during retrieval.
//
// Example:
//
//	templates, err := listModelArmorTemplates(
//	    os.Stdout,
//	    "my-project",
//	    "my-location",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, template := range templates {
//	    fmt.Println(template)
//	}
func listModelArmorTemplates(w io.Writer, projectID, location string) ([]*modelarmorpb.Template, error) {
	ctx := context.Background()

	// Create the Model Armor client.
	client, err := modelarmor.NewClient(ctx,
		option.WithEndpoint(fmt.Sprintf("modelarmor.%s.rep.googleapis.com:443", location)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Initialize request argument(s).
	req := &modelarmorpb.ListTemplatesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}

	// Get list of templates.
	it := client.ListTemplates(ctx, req)
	var templates []*modelarmorpb.Template

	for {
		template, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate templates: %v", err)
		}
		templates = append(templates, template)
		fmt.Fprintf(w, "Template: %s\n", template.Name)
	}

	// [END modelarmor_list_templates]

	return templates, nil
}
