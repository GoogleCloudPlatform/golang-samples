// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
func undeleteRole(w io.Writer, name string, projectID string) (*iam.Role, error) {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("iam.New: %v", err)
	}

	resource := "projects/" + projectID + "/roles/" + name
	request := iam.UndeleteRoleRequest{}
	role, err := service.Projects.Roles.Undelete(resource, &request).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.Roles.Undelete: %v", err)
	}
	fmt.Fprintf(w, "Undeleted role: %v", role.Name)
	return role, nil
}

// [END iam_undelete_role]
