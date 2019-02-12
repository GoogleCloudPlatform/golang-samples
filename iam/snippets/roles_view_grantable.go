// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

// [START iam_view_grantable_roles]
import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
)

// viewGrantableRoles lists roles grantable on a resource.
func viewGrantableRoles(w io.Writer, fullResourceName string) ([]*iam.Role, error) {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("iam.New: %v", err)
	}

	request := iam.QueryGrantableRolesRequest{FullResourceName: fullResourceName}
	response, err := service.Roles.QueryGrantableRoles(&request).Do()
	if err != nil {
		return nil, fmt.Errorf("Roles.QueryGrantableRoles: %v", err)
	}
	for _, role := range response.Roles {
		fmt.Fprintf(w, "Found grantable role: %v\n", role.Name)
	}
	return response.Roles, err
}

// [END iam_view_grantable_roles]
