// Copyright 2020 Google LLC
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

package buckets

// [START storage_set_bucket_public_iam]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/iam"

	"cloud.google.com/go/storage"
)

// setBucketPublicIAM makes all objects in a bucket publicly readable.
func setBucketPublicIAM(w io.Writer, bucketName string, policy iam.Policy) error {
	// bucketName := "bucket-name"
	// policy := client.Bucket(bucketName).IAM().Policy(ctx)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	role := iam.RoleName("roles/storage.objectViewer")
	policy.Add(iam.AllUsers, role)
	if err := client.Bucket(bucketName).IAM().SetPolicy(ctx, &policy); err != nil {
		return fmt.Errorf("Bucket(%q).IAM().SetPolicy: %v", bucketName, err)
	}
	fmt.Fprintf(w, "Bucket %v is now publicly readable\n", bucketName)
	return nil
}

// [END storage_set_bucket_public_iam]
