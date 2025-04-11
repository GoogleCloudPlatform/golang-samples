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

// [START iam_create_key]
import (
	"context"
	// "encoding/base64"
	"fmt"
	"io"

	iam "google.golang.org/api/iam/v1"
)

// createKey creates a service account key.
func createKey(w io.Writer, serviceAccountEmail string) (*iam.ServiceAccountKey, error) {
	ctx := context.Background()
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %w", err)
	}

	resource := "projects/-/serviceAccounts/" + serviceAccountEmail
	request := &iam.CreateServiceAccountKeyRequest{}
	key, err := service.Projects.ServiceAccounts.Keys.Create(resource, request).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.Keys.Create: %w", err)
	}
	// The PrivateKeyData field contains the base64-encoded service account key
	// in JSON format.
	// TODO(Developer): Save the below key (jsonKeyFile) to a secure location.
	// You cannot download it later.
	// jsonKeyFile, _ := base64.StdEncoding.DecodeString(key.PrivateKeyData)
	fmt.Fprintf(w, "Key created successfully")
	return key, nil
}

// [END iam_create_key]
