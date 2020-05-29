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

// [START servicedirectory_create_namespace]
import (
	"context"
	"fmt"
	"io"

	servicedirectory "cloud.google.com/go/servicedirectory/apiv1beta1"
	sdpb "google.golang.org/genproto/googleapis/cloud/servicedirectory/v1beta1"
)

func createNamespace(w io.Writer, projectID string) error {
	// projectID := "my-project"
	location := "us-east4"
	namespaceID := "golang-test-namespace"

	ctx := context.Background()
	// Create a registration client.
	client, err := servicedirectory.NewRegistrationClient(ctx)
	if err != nil {
		return fmt.Errorf("ServiceDirectory.NewRegistrationClient: %v", err)
	}

	defer client.Close()
	// Create a Namespace.
	req := &sdpb.CreateNamespaceRequest{
		Parent:      fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		NamespaceId: namespaceID,
	}
	resp, err := client.CreateNamespace(ctx, req)
	if err != nil {
		return fmt.Errorf("CreateNamespace: %v", err)
	}
	fmt.Fprintf(w, "servicedirectory.CreateNamespace result: %s\n", resp.Name)
	return nil
}

// [END servicedirectory_create_namespace]
