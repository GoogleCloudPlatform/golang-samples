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

// [START iam_edit_role]
import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
)

// editRole modifies a custom role.
func editRole(w io.Writer, projectID, name, newTitle, newDescription, newStage string, newPermissions []string) (*iam.Role, error) {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("iam.New: %v", err)
	}

	resource := "projects/" + projectID + "/roles/" + name
	role, err := service.Projects.Roles.Get(resource).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.Roles.Get: %v", err)
	}
	role.Title = newTitle
	role.Description = newDescription
	role.IncludedPermissions = newPermissions
	role.Stage = newStage
	role, err = service.Projects.Roles.Patch(resource, role).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.Roles.Patch: %v", err)
	}
	fmt.Fprintf(w, "Updated role: %v", role.Name)
	return role, nil
}

// [END iam_edit_role]
