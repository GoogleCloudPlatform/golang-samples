// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
func createRole(w io.Writer, name string, projectID string, title string, description string,
	permissions []string, stage string) (*iam.Role, error) {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("iam.New: %v", err)
	}

	request := iam.CreateRoleRequest{Role: &iam.Role{
		Title:               title,
		Description:         description,
		IncludedPermissions: permissions,
		Stage:               stage,
	}, RoleId: name}
	role, err := service.Projects.Roles.Create("projects/"+projectID, &request).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.Roles.Create: %v", err)
	}
	fmt.Fprintf(w, "Created role: %v", role.Name)
	return role, nil
}

// [END iam_create_role]
