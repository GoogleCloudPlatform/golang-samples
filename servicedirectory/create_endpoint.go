// Copyright 2020 Google LLC
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

package servicedirectory

// [START servicedirectory_create_endpoint]
import (
	"context"
	"fmt"
	"io"

	servicedirectory "cloud.google.com/go/servicedirectory/apiv1"
	sdpb "cloud.google.com/go/servicedirectory/apiv1/servicedirectorypb"
)

func createEndpoint(w io.Writer, projectID string) error {
	// projectID := "my-project"
	location := "us-east4"
	namespaceID := "golang-test-namespace"
	serviceID := "golang-test-service"
	endpointID := "golang-test-endpoint"

	ctx := context.Background()
	// Create a registration client.
	client, err := servicedirectory.NewRegistrationClient(ctx)
	if err != nil {
		return fmt.Errorf("ServiceDirectory.NewRegistrationClient: %w", err)
	}

	defer client.Close()
	// Create an Endpoint.
	req := &sdpb.CreateEndpointRequest{
		Parent:     fmt.Sprintf("projects/%s/locations/%s/namespaces/%s/services/%s", projectID, location, namespaceID, serviceID),
		EndpointId: endpointID,
		Endpoint: &sdpb.Endpoint{
			Address: "8.8.8.8",
			Port:    8080,
			Annotations: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}
	endpoint, err := client.CreateEndpoint(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateEndpoint: %w", err)
	}
	fmt.Fprintf(w, "servicedirectory.CreateEndpoint result: %s", endpoint.Name)
	return nil
}

// [END servicedirectory_create_endpoint]
