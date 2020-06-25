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

// [START servicedirectory_resolve_service]
import (
	"context"
	"fmt"
	"io"

	servicedirectory "cloud.google.com/go/servicedirectory/apiv1beta1"
	sdpb "google.golang.org/genproto/googleapis/cloud/servicedirectory/v1beta1"
)

func resolveService(w io.Writer, projectID string) error {
	// projectID := "my-project"
	location := "us-east4"
	namespaceID := "golang-test-namespace"
	serviceID := "golang-test-service"

	ctx := context.Background()
	// Create a lookup client.
	resolver, err := servicedirectory.NewLookupClient(ctx)
	if err != nil {
		return fmt.Errorf("ServiceDirectory.NewRegistrationClient: %v", err)
	}

	defer resolver.Close()
	// Now Resolve the service.
	req := &sdpb.ResolveServiceRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/namespaces/%s/services/%s", projectID, location, namespaceID, serviceID),
	}
	result, err := resolver.ResolveService(ctx, req)
	if err != nil {
		return fmt.Errorf("ResolveService: %v", err)
	}

	fmt.Fprintf(w, "Successfully Resolved Service %v\n", result)
	return nil
}

// [END servicedirectory_resolve_service]
