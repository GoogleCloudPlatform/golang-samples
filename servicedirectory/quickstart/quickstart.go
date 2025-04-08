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

// [START servicedirectory_quickstart]

// Sample quickstart is a program that uses Cloud Service Directory
// create. delete, and resolve functionality.
package main

import (
	"context"
	"fmt"
	"log"

	servicedirectory "cloud.google.com/go/servicedirectory/apiv1"
	sdpb "cloud.google.com/go/servicedirectory/apiv1/servicedirectorypb"
)

func main() {
	projectID := "your-project-id"
	location := "us-west1"
	serviceID := "golang-quickstart-service"
	namespaceID := "golang-quickstart-namespace"
	endpointID := "golang-quickstart-endpoint"

	ctx := context.Background()
	// Create a registration client.
	registry, err := servicedirectory.NewRegistrationClient(ctx)
	if err != nil {
		log.Fatalf("servicedirectory.NewRegistrationClient: %v", err)
	}
	defer registry.Close()

	// Create a lookup client.
	resolver, err := servicedirectory.NewLookupClient(ctx)
	if err != nil {
		log.Fatalf("servicedirectory.NewLookupClient: %v", err)
	}
	defer resolver.Close()

	// Create a Namespace.
	createNsReq := &sdpb.CreateNamespaceRequest{
		Parent:      fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		NamespaceId: namespaceID,
	}
	namespace, err := registry.CreateNamespace(ctx, createNsReq)
	if err != nil {
		log.Fatalf("servicedirectory.CreateNamespace: %v", err)
	}

	// Create a Service.
	createServiceReq := &sdpb.CreateServiceRequest{
		Parent:    namespace.Name,
		ServiceId: serviceID,
		Service: &sdpb.Service{
			Annotations: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}
	service, err := registry.CreateService(ctx, createServiceReq)
	if err != nil {
		log.Fatalf("servicedirectory.CreateService: %v", err)
	}

	// Create an Endpoint.
	createEndpointReq := &sdpb.CreateEndpointRequest{
		Parent:     service.Name,
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
	_, err = registry.CreateEndpoint(ctx, createEndpointReq)
	if err != nil {
		log.Fatalf("servicedirectory.CreateEndpoint: %v", err)
	}

	// Now Resolve the service.
	lookupRequest := &sdpb.ResolveServiceRequest{
		Name: service.Name,
	}
	result, err := resolver.ResolveService(ctx, lookupRequest)
	if err != nil {
		log.Fatalf("servicedirectory.ResolveService: %v", err)
		return
	}

	fmt.Printf("Successfully Resolved Service %v", result)

	// Delete the namespace.
	deleteNsReq := &sdpb.DeleteNamespaceRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/namespaces/%s",
			projectID, location, namespaceID),
	}
	registry.DeleteNamespace(ctx, deleteNsReq) // Ignore results.
}

// [END servicedirectory_quickstart]
