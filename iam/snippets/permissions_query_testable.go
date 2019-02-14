// Copyright 2019 Google LLC
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

package snippets

// [START iam_query_testable_permissions]
import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
)

// queryTestablePermissions lists testable permissions on a resource.
func queryTestablePermissions(w io.Writer, fullResourceName string) ([]*iam.Permission, error) {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("iam.New: %v", err)
	}

	request := &iam.QueryTestablePermissionsRequest{
		FullResourceName: fullResourceName,
	}
	response, err := service.Permissions.QueryTestablePermissions(request).Do()
	if err != nil {
		return nil, fmt.Errorf("Permissions.QueryTestablePermissions: %v", err)
	}
	for _, p := range response.Permissions {
		fmt.Fprintf(w, "Found permissions: %v", p.Name)
	}
	return response.Permissions, nil
}

// [END iam_query_testable_permissions]
