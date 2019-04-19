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

// [START iam_undelete_role]
import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
)

// undeleteRole restores a deleted custom role.
func undeleteRole(w io.Writer, projectID, name string) (*iam.Role, error) {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("iam.New: %v", err)
	}

	resource := "projects/" + projectID + "/roles/" + name
	request := &iam.UndeleteRoleRequest{}
	role, err := service.Projects.Roles.Undelete(resource, request).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.Roles.Undelete: %v", err)
	}
	fmt.Fprintf(w, "Undeleted role: %v", role.Name)
	return role, nil
}

// [END iam_undelete_role]
