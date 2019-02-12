// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

// [START iam_create_service_account]
import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
)

// createServiceAccount creates a service account.
func createServiceAccount(w io.Writer, projectID string, name string,
	displayName string) (*iam.ServiceAccount, error) {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return nil, fmt.Errorf("iam.New: %v", err)
	}

	request := iam.CreateServiceAccountRequest{
		AccountId:      name,
		ServiceAccount: &iam.ServiceAccount{DisplayName: displayName}}
	account, err := service.Projects.ServiceAccounts.Create("projects/"+projectID, &request).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.Create: %v", err)
	}
	fmt.Fprintf(w, "Created service account: %v", account)
	return account, nil
}

// [END iam_create_service_account]
