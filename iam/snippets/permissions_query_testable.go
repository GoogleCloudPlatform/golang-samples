// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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

	request := iam.QueryTestablePermissionsRequest{FullResourceName: fullResourceName}
	response, err := service.Permissions.QueryTestablePermissions(&request).Do()
	if err != nil {
		return nil, fmt.Errorf("Permissions.QueryTestablePermissions: %v", err)
	}
	for _, p := range response.Permissions {
		fmt.Fprintf(w, "Found permissions: %v", p.Name)
	}
	return response.Permissions, nil
}

// [END iam_query_testable_permissions]
