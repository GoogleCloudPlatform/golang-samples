// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

// [START iam_get_role]
import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
)

// getRole gets role metadata.
func getRole(w io.Writer, name string) (*iam.Role, error) {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("iam.New: %v", err)
	}

	role, err := service.Roles.Get(name).Do()
	if err != nil {
		return nil, fmt.Errorf("Roles.Get: %v", err)
	}
	fmt.Fprintf(w, "Got role: %v\n", role.Name)
	for _, permission := range role.IncludedPermissions {
		fmt.Fprintf(w, "Got permission: %v\n", permission)
	}
	return role, nil
}

// [END iam_get_role]
