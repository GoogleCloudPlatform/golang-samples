// Copyright 2026 Google LLC
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

package parametermanager

// [START parametermanager_disable_param_template_version]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// disableParamTemplateVersion disables a version of a parameter template.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// templateID: The ID of the template.
// versionID: The ID of the template version.
//
// The function returns an error if the operation fails.
func disableParamTemplateVersion(w io.Writer, projectID, templateID, versionID string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the template version.
	name := fmt.Sprintf("projects/%s/locations/global/templates/%s/versions/%s", projectID, templateID, versionID)

	// Build the request to update the disabled state of the template version.
	req := &parametermanagerpb.UpdateTemplateVersionRequest{
		TemplateVersion: &parametermanagerpb.TemplateVersion{
			Name:     name,
			Disabled: true,
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"disabled"}},
	}

	// Call the API to update the template version.
	version, err := client.UpdateTemplateVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update template version: %w", err)
	}

	fmt.Fprintf(w, "Disabled parameter template version: %s\n", version.Name)
	return nil
}

// [END parametermanager_disable_param_template_version]
