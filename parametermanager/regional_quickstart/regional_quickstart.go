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

// [START parametermanager_regional_quickstart]
package main

// Sample quickstart is a basic program that uses Parameter Manager.
import (
	"context"
	"fmt"
	"log"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	parametermanagerpb "cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"google.golang.org/api/option"
)

func main() {
	// GCP project in which to store parameter in Parameter Manager.
	projectID := "test-project-id"
	// Location at which you want to store your parameters
	locationID := "us-central1"
	// Id of the parameter which you want to create
	parameterID := "test-parameter-id"
	// Id of the parameter version which you want to create
	versionID := "test-version-id"
	payload := `{"username": "test-user", "host": "localhost"}`

	// Create a new context.
	ctx := context.Background()

	// Create a Parameter Manager client.
	endpoint := fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", locationID)
	client, err := parametermanager.NewClient(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		log.Fatalf("Failed to create Parameter Manager client: %v\n", err)
	}
	defer client.Close()

	// Construct the name of the parent resource to create the parameter.
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)

	// Create a new parameter.
	parameter, err := client.CreateParameter(ctx, &parametermanagerpb.CreateParameterRequest{
		Parent:      parent,
		ParameterId: parameterID,
		Parameter: &parametermanagerpb.Parameter{
			Format: parametermanagerpb.ParameterFormat_JSON,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create parameter: %v\n", err)
	}
	fmt.Printf("Created regional parameter %s with format %s\n", parameter.Name, parameter.Format.String())

	parameterVersion, err := client.CreateParameterVersion(ctx, &parametermanagerpb.CreateParameterVersionRequest{
		Parent:             parameter.Name,
		ParameterVersionId: versionID,
		ParameterVersion: &parametermanagerpb.ParameterVersion{
			Payload: &parametermanagerpb.ParameterVersionPayload{
				Data: []byte(payload),
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create parameter version: %v\n", err)
	}
	fmt.Printf("Created regional parameter version %s\n", parameterVersion.Name)

	getParameterVersion, err := client.GetParameterVersion(ctx, &parametermanagerpb.GetParameterVersionRequest{
		Name: parameterVersion.Name,
	})
	if err != nil {
		log.Fatalf("Failed to get parameter version: %v\n", err)
	}

	fmt.Printf("Retrieved regional parameter version: %s\n", getParameterVersion.Name)
	fmt.Printf("Payload: %s\n", getParameterVersion.Payload.Data)
}

// [END parametermanager_regional_quickstart]
