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

// [START storage_get_bucket_policy]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/storage"
)

// getBucketPolicy gets the bucket IAM policy.
func getBucketPolicy(w io.Writer, bucketName string) (*iam.Policy, error) {
	// bucketName := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	policy, err := client.Bucket(bucketName).IAM().Policy(ctx)
	if err != nil {
		return nil, fmt.Errorf("Bucket(%q).IAM().Policy: %v", bucketName, err)
	}
	for _, role := range policy.Roles() {
		fmt.Fprintf(w, "%q: %q", role, policy.Members(role))
	}
	return policy, nil
}

// [END storage_get_bucket_policy]
