// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
func editRole(w io.Writer, name string, projectID string, newTitle string, newDescription string,
	newPermissions []string, newStage string) (*iam.Role, error) {
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
