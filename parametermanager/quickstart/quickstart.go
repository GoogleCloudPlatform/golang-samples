// Copyright 2025 Google LLC
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

// [START parametermanager_quickstart]

// Sample quickstart is a basic program that uses Parameter Manager.
import (
	"context"
	"fmt"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
)

func main() {
	// GCP project in which to store parameters in Parameter Manager.
	projectID := "test-project-id"
	// Id of the parameter which you want to create
	parameterID := "test-parameter-id"
	// Id of the parameter version which you want to create
	versionID := "test-version-id"
	payload := `{"username": "test-user", "host": "localhost"}`

	// Create a context.
	ctx := context.Background()

	// Create a Parameter Manager client.
	client, err := parametermanager.NewClient(ctx)
	if err != nil {
		fmt.Println("Failed to create Parameter Manager client: %v", err)
		return
	}
	defer client.Close()

	// Construct the name of the parent resource to create the parameter.
	parent := fmt.Sprintf("projects/%s/locations/global", projectID)

	// Build a parameter creation request.
	paramCreateReq := &parametermanagerpb.CreateParameterRequest{
		Parent:      parent,
		ParameterId: parameterID,
		Parameter: &parametermanagerpb.Parameter{
			Format: parametermanagerpb.ParameterFormat_JSON,
		},
	}

	// Create a new parameter.
	parameter, err := client.CreateParameter(ctx, paramCreateReq)
	if err != nil {
		fmt.Printf("Failed to create parameter: %v\n", err)
		return
	}
	fmt.Printf("Created parameter %s with format %s\n", parameter.Name, parameter.Format.String())

	// Build a parameter version creation request.
	paramVersionCreateReq := &parametermanagerpb.CreateParameterVersionRequest{
		Parent:             parameter.Name,
		ParameterVersionId: versionID,
		ParameterVersion: &parametermanagerpb.ParameterVersion{
			Payload: &parametermanagerpb.ParameterVersionPayload{
				Data: []byte(payload),
			},
		},
	}

	parameterVersion, err := client.CreateParameterVersion(ctx, paramVersionCreateReq)
	if err != nil {
		fmt.Printf("Failed to create parameter version: %v\n", err)
		return
	}
	fmt.Printf("Created parameter version: %s\n", parameterVersion.Name)

	// Build a parameter version get request.
	paramVersionGetReq := &parametermanagerpb.GetParameterVersionRequest{
		Name: parameterVersion.Name,
	}

	// Get the parameter version.
	result, err := client.GetParameterVersion(ctx, paramVersionGetReq)
	if err != nil {
		fmt.Printf("Failed to get parameter version: %v\n", err)
		return
	}

	fmt.Printf("Retrieved parameter version: %s", result.Name)
	fmt.Printf("Payload: %s\n", result.Payload.Data)
}

// [END parametermanager_quickstart]
