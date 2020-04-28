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

// Sample quickstart is a basic program that uses Cloud Service Directory.
package main

import (
        "context"
        "fmt"

        servicedirectory "cloud.google.com/go/servicedirectory/apiv1beta1"
        sdpb "google.golang.org/genproto/googleapis/cloud/servicedirectory/v1beta1"
)

func main() {
        projectId := "your-project-id"
        location := "us-west1"
        serviceId := "golang-quickstart-service"
        namespaceId := "golang-quickstart-namespace"
        endpointId := "golang-quickstart-endpoint"

        ctx := context.Background()
        // Create a registration client.
        registry, err := servicedirectory.NewRegistrationClient(ctx)
        if err != nil {
                fmt.Printf("servicedirectory.NewRegistrationClient: %v", err)
                return
        }

        // Create a lookup client.
        resolver, err := servicedirectory.NewLookupClient(ctx)
        if err != nil {
                fmt.Printf("servicedirectory.NewLookupClient: %v", err)
                return
        }

	// Create a Namespace.
        createNsReq := &sdpb.CreateNamespaceRequest{
                Parent:      fmt.Sprintf("projects/%s/locations/%s", projectId, location),
                NamespaceId: namespaceId,
        }
        namespace, err := registry.CreateNamespace(ctx, createNsReq)
        if err != nil {
                fmt.Printf("servicedirectory.CreateNamespace: %v", err)
                return
        }

        // Create a Service.
        createServiceReq := &sdpb.CreateServiceRequest{
                Parent:    namespace.Name,
                ServiceId: serviceId,
                Service: &sdpb.Service{
                        Metadata: map[string]string{
                                "key1": "value1",
                                "key2": "value2",
                        },
                },
        }
        service, err := registry.CreateService(ctx, createServiceReq)
        if err != nil {
                fmt.Printf("servicedirectory.CreateService: %v", err)
                return
        }

        // Create an Endpoint.
        createEndpointReq := &sdpb.CreateEndpointRequest{
                Parent:     service.Name,
                EndpointId: endpointId,
                Endpoint: &sdpb.Endpoint{
                        Address: "10.10.10.10",
                        Port:    8080,
                        Metadata: map[string]string{
                                "key1": "value1",
                                "key2": "value2",
                        },
                },
        }
        _, err = registry.CreateEndpoint(ctx, createEndpointReq)
        if err != nil {
                fmt.Printf("servicedirectory.CreateEndpoint: %v", err)
                return
        }

        // Now Resolve the service.
        lookupRequest := &sdpb.ResolveServiceRequest{
                Name: service.Name,
        }
        result, err := resolver.ResolveService(ctx, lookupRequest)
        if err != nil {
                fmt.Printf("servicedirectory.ResolveService: %v", err)
                return
        }

        fmt.Printf("Successfully Resolved Service %v", result)

	// Delete the namespace.
        deleteNsReq := &sdpb.DeleteNamespaceRequest{
                Name: fmt.Sprintf("projects/%s/locations/%s/namespaces/%s",
                        projectId, location, namespaceId),
        }
        registry.DeleteNamespace(ctx, deleteNsReq) // Ignore results.
}

// [END servicedirectory_quickstart]
