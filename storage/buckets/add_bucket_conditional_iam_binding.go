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

// [START storage_add_bucket_conditional_iam_binding]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/storage"
	"google.golang.org/genproto/googleapis/type/expr"
)

// addBucketConditionalIAMBinding adds bucket conditional IAM binding.
func addBucketConditionalIAMBinding(w io.Writer, bucketName, role, member, title, description, expression string) error {
	// bucketName := "bucket-name"
	// role := "bucket-level IAM role"
	// member := "bucket-level IAM member"
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

	policy.Bindings = append(policy.Bindings, &iampb.Binding{
		Role:    role,
		Members: []string{member},
		Condition: &expr.Expr{
			Title:       title,
			Description: description,
			Expression:  expression,
		},
	})

	if err := bucket.IAM().V3().SetPolicy(ctx, policy); err != nil {
		return fmt.Errorf("Bucket(%q).IAM().V3().SetPolicy: %w", bucketName, err)
	}
	// NOTE: It may be necessary to retry this operation if IAM policies are
	// being modified concurrently. SetPolicy will return an error if the policy
	// was modified since it was retrieved.
	fmt.Fprintf(w, "Added %v with role %v to %v with condition %v %v %v\n", member, role, bucketName, title, description, expression)
	return nil
}

// [END storage_add_bucket_conditional_iam_binding]
