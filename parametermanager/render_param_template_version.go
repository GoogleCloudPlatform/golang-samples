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

// [START parametermanager_render_param_template_version]
import (
	"context"
	"fmt"
	"io"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
)

// renderParamTemplateVersion renders a template version using the values of a parameter version.
//
// w: The io.Writer object used to write the output.
// projectID: The ID of the project where the parameter is located.
// templateID: The ID of the template.
// versionID: The ID of the template version.
// parameterVersionName: The full resource name of the parameter version whose values are used for rendering.
//
// The function returns an error if the operation fails.
func renderParamTemplateVersion(w io.Writer, projectID, templateID, versionID, parameterVersionName string) error {
	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Parameter Manager client: %w", err)
	}
	defer client.Close()

	// Construct the name of the template version.
	name := fmt.Sprintf("projects/%s/locations/global/templates/%s/versions/%s", projectID, templateID, versionID)

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

	fmt.Fprintf(w, "Rendered parameter template version: %s\n", rendered.TemplateVersion)
	fmt.Fprintf(w, "Template payload: %s\n", rendered.Payload.Data)

	// If the parameter contains secret references, they will be resolved
	// and the actual secret values will be included in the rendered output.
	// Be cautious with logging or displaying this information.
	fmt.Fprintf(w, "Rendered payload: %s\n", rendered.RenderedPayload)
	return nil
}

// [END parametermanager_render_param_template_version]
