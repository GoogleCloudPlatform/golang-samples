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

package main

// [START parametermanager_regional_templates_quickstart]

// Sample quickstart is a basic program that renders a Parameter Manager
// template using the values of a parameter version.
import (
	"context"
	"fmt"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"google.golang.org/api/option"
)

func main() {
	// GCP project in which to store resources in Parameter Manager.
	projectID := "test-project-id"
	// Location at which you want to store your resources.
	locationID := "us-central1"
	// Id of the template which you want to create.
	templateID := "test-template-id"
	// Id of the template version which you want to create.
	templateVersionID := "test-template-version-id"
	// Id of the parameter which you want to create.
	parameterID := "test-parameter-id"
	// Id of the parameter version which you want to create.
	parameterVersionID := "test-parameter-version-id"

	// The template payload contains {{.variableName}} placeholders that are
	// filled in with the values of the parameter version when rendered.
	templatePayload := `{"username": "{{.username}}", "host": "{{.host}}"}`
	parameterPayload := `{"username": "test-user", "host": "localhost"}`

	// Create a context and a Parameter Manager client.
	ctx := context.Background()
	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationID)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		fmt.Printf("Failed to create Parameter Manager client: %v\n", err)
		return
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)

	// Create a JSON template.
	template, err := client.CreateTemplate(ctx, &parametermanagerpb.CreateTemplateRequest{
		Parent:     parent,
		TemplateId: templateID,
		Template: &parametermanagerpb.Template{
			Format: parametermanagerpb.TemplateFormat_TEMPLATE_FORMAT_JSON,
		},
	})
	if err != nil {
		fmt.Printf("Failed to create template: %v\n", err)
		return
	}
	fmt.Printf("Created template %s with format %s\n", template.Name, template.Format)

	// Create a template version containing the placeholders.
	templateVersion, err := client.CreateTemplateVersion(ctx, &parametermanagerpb.CreateTemplateVersionRequest{
		Parent:            template.Name,
		TemplateVersionId: templateVersionID,
		TemplateVersion: &parametermanagerpb.TemplateVersion{
			Payload: &parametermanagerpb.TemplateVersionPayload{
				Data: []byte(templatePayload),
			},
		},
	})
	if err != nil {
		fmt.Printf("Failed to create template version: %v\n", err)
		return
	}
	fmt.Printf("Created template version: %s\n", templateVersion.Name)

	// Create a JSON parameter and a version holding the values to render with.
	parameter, err := client.CreateParameter(ctx, &parametermanagerpb.CreateParameterRequest{
		Parent:      parent,
		ParameterId: parameterID,
		Parameter: &parametermanagerpb.Parameter{
			Format: parametermanagerpb.ParameterFormat_JSON,
		},
	})
	if err != nil {
		fmt.Printf("Failed to create parameter: %v\n", err)
		return
	}
	fmt.Printf("Created parameter %s with format %s\n", parameter.Name, parameter.Format)

	parameterVersion, err := client.CreateParameterVersion(ctx, &parametermanagerpb.CreateParameterVersionRequest{
		Parent:             parameter.Name,
		ParameterVersionId: parameterVersionID,
		ParameterVersion: &parametermanagerpb.ParameterVersion{
			Payload: &parametermanagerpb.ParameterVersionPayload{
				Data: []byte(parameterPayload),
			},
		},
	})
	if err != nil {
		fmt.Printf("Failed to create parameter version: %v\n", err)
		return
	}
	fmt.Printf("Created parameter version: %s\n", parameterVersion.Name)

	// Render the template version using the parameter version.
	rendered, err := client.RenderTemplateVersion(ctx, &parametermanagerpb.RenderTemplateVersionRequest{
		Name:             templateVersion.Name,
		ParameterVersion: parameterVersion.Name,
	})
	if err != nil {
		fmt.Printf("Failed to render template version: %v\n", err)
		return
	}
	fmt.Printf("Rendered payload: %s\n", rendered.RenderedPayload)
}

// [END parametermanager_regional_templates_quickstart]
