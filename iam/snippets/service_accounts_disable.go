// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package snippets

// [START iam_disable_service_account]
import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	iam "google.golang.org/api/iam/v1"
)

// disableServiceAccount disables a service account.
func disableServiceAccount(w io.Writer, email string) error {
	// email:= service-account@your-project.iam.gserviceaccount.com
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return fmt.Errorf("google.DefaultClient: %v", err)
	}
	service, err := iam.New(client)
	if err != nil {
		return fmt.Errorf("iam.New: %v", err)
	}

	request := &iam.DisableServiceAccountRequest{}
	_, err = service.Projects.ServiceAccounts.Disable("projects/-/serviceAccounts/"+email, request).Do()
	if err != nil {
		return fmt.Errorf("Projects.ServiceAccounts.Disable: %v", err)
	}
	fmt.Fprintf(w, "Disabled service account: %v", email)
	return nil
}

// [END iam_disable_service_account]
