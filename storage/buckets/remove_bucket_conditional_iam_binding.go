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

// [START storage_remove_bucket_conditional_iam_binding]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// removeBucketConditionalIAMBinding removes bucket conditional IAM binding.
func removeBucketConditionalIAMBinding(w io.Writer, bucketName, role, title, description, expression string) error {
	// bucketName := "bucket-name"
	// role := "bucket-level IAM role"
	// title := "condition title"
	// description := "condition description"
	// expression := "condition expression"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	bucket := client.Bucket(bucketName)
	policy, err := bucket.IAM().V3().Policy(ctx)
	if err != nil {
		return fmt.Errorf("Bucket(%q).IAM().V3().Policy: %w", bucketName, err)
	}

	// Find the index of the binding matching inputs.
	i := -1
	for j, binding := range policy.Bindings {
		if binding.Role == role && binding.Condition != nil {
			condition := binding.Condition
			if condition.Title == title &&
				condition.Description == description &&
				condition.Expression == expression {
				i = j
			}
		}
	}

	if i == -1 {
		return fmt.Errorf("no matching binding group found")
	}

	// Get a slice of the bindings, removing the binding at index i.
	policy.Bindings = append(policy.Bindings[:i], policy.Bindings[i+1:]...)

	if err := bucket.IAM().V3().SetPolicy(ctx, policy); err != nil {
		return fmt.Errorf("Bucket(%q).IAM().V3().SetPolicy: %w", bucketName, err)
	}
	// NOTE: It may be necessary to retry this operation if IAM policies are
	// being modified concurrently. SetPolicy will return an error if the policy
	// was modified since it was retrieved.
	fmt.Fprintln(w, "Conditional binding was removed")
	return nil
}

// [END storage_remove_bucket_conditional_iam_binding]
