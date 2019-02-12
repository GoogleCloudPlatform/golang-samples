// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

// [START iam_delete_key]
import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
)

// deleteKey deletes a service account key.
func deleteKey(w io.Writer, fullKeyName string) error {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return fmt.Errorf("iam.New: %v", err)
	}

	_, err = service.Projects.ServiceAccounts.Keys.Delete(fullKeyName).Do()
	if err != nil {
		return fmt.Errorf("Projects.ServiceAccounts.Keys.Delete: %v", err)
	}
	fmt.Fprintf(w, "Deleted key: %v", fullKeyName)
	return nil
}

// [END iam_delete_key]
