// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

// [START iam_delete_service_account]
import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
)

// deleteServiceAccount deletes a service account.
func deleteServiceAccount(w io.Writer, email string) error {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return fmt.Errorf("iam.New: %v", err)
	}

	_, err = service.Projects.ServiceAccounts.Delete("projects/-/serviceAccounts/" + email).Do()
	if err != nil {
		return fmt.Errorf("Projects.ServiceAccounts.Delete: %v", err)
	}
	fmt.Fprintf(w, "Deleted service account: %v", email)
	return nil
}

// [END iam_delete_service_account]
