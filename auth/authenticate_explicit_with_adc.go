// Copyright 2022 Google LLC
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

// [START auth_cloud_explicit_adc]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/auth/credentials"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// authenticateExplicitWithAdc uses Application Default Credentials
// to print storage buckets.
func authenticateExplicitWithAdc(w io.Writer) error {
	ctx := context.Background()

	// Construct Credentials which obtains the default configuration from your
	// working environment.
	credentials, err := credentials.DetectDefault(&credentials.DetectOptions{
		Scopes: []string{"https://www.googleapis.com/auth/cloud-platform"},
	})
	if err != nil {
		return fmt.Errorf("failed to generate default credentials: %w", err)
	}
	// If you are authenticating to a Cloud API, you can let the library include
	// the default scope, https://www.googleapis.com/auth/cloud-platform,
	// because IAM is used to provide fine-grained permissions for Cloud. For
	// more information on scopes to use, see:
	// https://developers.google.com/identity/protocols/oauth2/scopes

	// Construct the Storage client.
	client, err := storage.NewClient(ctx, option.WithAuthCredentials(credentials))
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	projID, err := credentials.ProjectID(ctx)
	if err != nil {
		return fmt.Errorf("ProjectID: %w", err)
	}

	it := client.Buckets(ctx, projID)
	for {
		bucketAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "Bucket: %v\n", bucketAttrs.Name)
	}

	fmt.Fprintf(w, "Listed all storage buckets.\n")

	return nil
}

// [END auth_cloud_explicit_adc]
