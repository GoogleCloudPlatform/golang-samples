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

// [START auth_cloud_implicit_adc]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// authenticateImplicitWithAdc uses Application Default Credentials
// to automatically find credentials and authenticate.
func authenticateImplicitWithAdc(w io.Writer, projectId string) error {
	// projectId := "your_project_id"

	ctx := context.Background()

	// NOTE: Replace the client created below with the client required for your application.
	// Note that the credentials are not specified when constructing the client.
	// The client library finds your credentials using ADC.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	it := client.Buckets(ctx, projectId)
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

// [END auth_cloud_implicit_adc]
