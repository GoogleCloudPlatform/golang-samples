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

// [START iam_create_role]
import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
)

// createRole creates a custom role.
func createRole(w io.Writer, projectID, name, title, description, stage string, permissions []string) (*iam.Role, error) {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("iam.New: %v", err)
	}

	request := &iam.CreateRoleRequest{
		Role: &iam.Role{
			Title:               title,
			Description:         description,
			IncludedPermissions: permissions,
			Stage:               stage,
		},
		RoleId: name,
	}
	role, err := service.Projects.Roles.Create("projects/"+projectID, request).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.Roles.Create: %v", err)
	}
	fmt.Fprintf(w, "Created role: %v", role.Name)
	return role, nil
}

// [END iam_create_role]
