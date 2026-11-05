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

package regional_parametermanager

// [START parametermanager_render_regional_param_template_version]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"google.golang.org/api/option"
)

// renderRegionalParamTemplateVersion renders a template version.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// locationID: The region where the resources are located.
// templateID: The ID of the template.
// versionID: The ID of the template version.
// parameterVersionName: The resource name of the parameter version.
//
// The function returns an error if the operation fails.
func renderRegionalParamTemplateVersion(w io.Writer, projectID, locationID, templateID, versionID, parameterVersionName string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationID)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the template version.
	name := fmt.Sprintf("projects/%s/locations/%s/templates/%s/versions/%s", projectID, locationID, templateID, versionID)

	// Build the request to render the template version.
	req := &parametermanagerpb.RenderTemplateVersionRequest{
		Name:             name,
		ParameterVersion: parameterVersionName,
	}

	// Call the API to render the template version.
	rendered, err := client.RenderTemplateVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to render template version: %w", err)
	}

	fmt.Fprintf(w, "Rendered regional parameter template version: %s\n", rendered.TemplateVersion)
	fmt.Fprintf(w, "Template payload: %s\n", rendered.Payload.Data)

	fmt.Fprintf(w, "Rendered payload: %s\n", rendered.RenderedPayload)
	return nil
}

// [END parametermanager_render_regional_param_template_version]
