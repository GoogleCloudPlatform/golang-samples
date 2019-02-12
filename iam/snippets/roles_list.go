// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

// [START iam_list_roles]
import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
)

// listRoles lists a project's roles.
func listRoles(w io.Writer, projectID string) ([]*iam.Role, error) {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("iam.New: %v", err)
	}

	response, err := service.Projects.Roles.List("projects/" + projectID).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.Roles.List: %v", err)
	}
	for _, role := range response.Roles {
		fmt.Fprintf(w, "Listing role: %v\n", role.Name)
	}
	return response.Roles, nil
}

// [END iam_list_roles]
