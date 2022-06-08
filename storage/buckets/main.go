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

// Sample buckets creates a bucket, lists buckets and deletes a bucket
// using the Google Storage API. More documentation is available at
// https://cloud.google.com/storage/docs/json_api/v1/.
//go:build ignore
// +build ignore

package buckets

import (
	"context"
	"errors"
	"log"
	"time"

	"cloud.google.com/go/storage"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/genproto/googleapis/type/expr"
)

// TODO: Move remaining region tags in this file to separate stand-alone files,
// then delete this file.
func main() {
	log.Fatalf("Running main is not supported.")
}

func addUser(c *storage.Client, bucketName string) error {
	// [START add_bucket_iam_member]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	bucket := c.Bucket(bucketName)
	policy, err := bucket.IAM().V3().Policy(ctx)
	if err != nil {
		return err
	}
	// Other valid prefixes are "serviceAccount:", "user:"
	// See the documentation for more values.
	// https://cloud.google.com/storage/docs/access-control/iam
	policy.Bindings = append(policy.Bindings, &iampb.Binding{
		Role:    "roles/storage.objectViewer",
		Members: []string{"group:cloud-logs@google.com"},
	})
	if err := bucket.IAM().V3().SetPolicy(ctx, policy); err != nil {
		return err
	}
	// NOTE: It may be necessary to retry this operation if IAM policies are
	// being modified concurrently. SetPolicy will return an error if the policy
	// was modified since it was retrieved.
	// [END add_bucket_iam_member]
	return nil
}

func removeUser(c *storage.Client, bucketName string) error {
	// [START remove_bucket_iam_member]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	bucket := c.Bucket(bucketName)
	policy, err := bucket.IAM().V3().Policy(ctx)
	if err != nil {
		return err
	}
	// Other valid prefixes are "serviceAccount:", "user:"
	// See the documentation for more values.
	// https://cloud.google.com/storage/docs/access-control/iam
	for _, binding := range policy.Bindings {
		// Only remove unconditional bindings matching role
		if binding.Role == "roles/storage.objectViewer" && binding.Condition == nil {
			// Filter out member.
			i := -1
			for j, member := range binding.Members {
				if member == "group:cloud-logs@google.com" {
					i = j
				}
			}

			if i == -1 {
				return errors.New("No matching binding group found.")
			} else {
				binding.Members = append(binding.Members[:i], binding.Members[i+1:]...)
			}
		}
	}
	if err := bucket.IAM().V3().SetPolicy(ctx, policy); err != nil {
		return err
	}
	// NOTE: It may be necessary to retry this operation if IAM policies are
	// being modified concurrently. SetPolicy will return an error if the policy
	// was modified since it was retrieved.
	// [END remove_bucket_iam_member]
	return nil
}

func addBucketConditionalIAMBinding(c *storage.Client, bucketName string, role string, member string, title string, description string, expression string) error {
	// [START storage_add_bucket_conditional_iam_binding]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	bucket := c.Bucket(bucketName)
	policy, err := bucket.IAM().V3().Policy(ctx)
	if err != nil {
		return err
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
		return err
	}
	// NOTE: It may be necessary to retry this operation if IAM policies are
	// being modified concurrently. SetPolicy will return an error if the policy
	// was modified since it was retrieved.
	// [END storage_add_bucket_conditional_iam_binding]
	return nil
}

func removeBucketConditionalIAMBinding(c *storage.Client, bucketName string, role string, title string, description string, expression string) error {
	// [START storage_remove_bucket_conditional_iam_binding]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	bucket := c.Bucket(bucketName)
	policy, err := bucket.IAM().V3().Policy(ctx)
	if err != nil {
		return err
	}

	// Find the index of the binding matching inputs
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
		return errors.New("No matching binding group found.")
	}

	// Get a slice of the bindings, removing the binding at index i
	policy.Bindings = append(policy.Bindings[:i], policy.Bindings[i+1:]...)

	if err := bucket.IAM().V3().SetPolicy(ctx, policy); err != nil {
		return err
	}
	// NOTE: It may be necessary to retry this operation if IAM policies are
	// being modified concurrently. SetPolicy will return an error if the policy
	// was modified since it was retrieved.
	// [END storage_remove_bucket_conditional_iam_binding]
	return nil
}

func setRetentionPolicy(c *storage.Client, bucketName string, retentionPeriod time.Duration) error {
	// [START storage_set_retention_policy]
	ctx := context.Background()

	bucket := c.Bucket(bucketName)
	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		RetentionPolicy: &storage.RetentionPolicy{
			RetentionPeriod: retentionPeriod,
		},
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return err
	}
	// [END storage_set_retention_policy]
	return nil
}

func removeRetentionPolicy(c *storage.Client, bucketName string) error {
	// [START storage_remove_retention_policy]
	ctx := context.Background()
	bucket := c.Bucket(bucketName)

	attrs, err := c.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return err
	}
	if attrs.RetentionPolicy.IsLocked {
		return errors.New("retention policy is locked")
	}

	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		RetentionPolicy: &storage.RetentionPolicy{},
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return err
	}
	// [END storage_remove_retention_policy]
	return nil
}

func lockRetentionPolicy(c *storage.Client, bucketName string) error {
	// [START storage_lock_retention_policy]
	ctx := context.Background()
	bucket := c.Bucket(bucketName)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	attrs, err := c.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return err
	}

	conditions := storage.BucketConditions{
		MetagenerationMatch: attrs.MetaGeneration,
	}
	if err := bucket.If(conditions).LockRetentionPolicy(ctx); err != nil {
		return err
	}

	lockedAttrs, err := c.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return err
	}
	log.Printf("Retention policy for %v is now locked\n", bucketName)
	log.Printf("Retention policy effective as of %v\n",
		lockedAttrs.RetentionPolicy.EffectiveTime)
	// [END storage_lock_retention_policy]
	return nil
}

func getRetentionPolicy(c *storage.Client, bucketName string) (*storage.BucketAttrs, error) {
	// [START storage_get_retention_policy]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	attrs, err := c.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return nil, err
	}
	if attrs.RetentionPolicy != nil {
		log.Print("Retention Policy\n")
		log.Printf("period: %v\n", attrs.RetentionPolicy.RetentionPeriod)
		log.Printf("effective time: %v\n", attrs.RetentionPolicy.EffectiveTime)
		log.Printf("policy locked: %v\n", attrs.RetentionPolicy.IsLocked)
	}
	// [END storage_get_retention_policy]
	return attrs, nil
}

func enableDefaultEventBasedHold(c *storage.Client, bucketName string) error {
	// [START storage_enable_default_event_based_hold]
	ctx := context.Background()

	bucket := c.Bucket(bucketName)
	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		DefaultEventBasedHold: true,
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return err
	}
	// [END storage_enable_default_event_based_hold]
	return nil
}

func disableDefaultEventBasedHold(c *storage.Client, bucketName string) error {
	// [START storage_disable_default_event_based_hold]
	ctx := context.Background()

	bucket := c.Bucket(bucketName)
	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		DefaultEventBasedHold: false,
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return err
	}
	// [END storage_disable_default_event_based_hold]
	return nil
}

func getDefaultEventBasedHold(c *storage.Client, bucketName string) (*storage.BucketAttrs, error) {
	// [START storage_get_default_event_based_hold]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	attrs, err := c.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("Default event-based hold enabled? %t\n",
		attrs.DefaultEventBasedHold)
	// [END storage_get_default_event_based_hold]
	return attrs, nil
}

func setDefaultKMSkey(c *storage.Client, bucketName string, keyName string) error {
	// [START storage_set_bucket_default_kms_key]
	ctx := context.Background()

	bucket := c.Bucket(bucketName)
	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		Encryption: &storage.BucketEncryption{DefaultKMSKeyName: keyName},
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return err
	}
	// [END storage_set_bucket_default_kms_key]
	return nil
}

func enableUniformBucketLevelAccess(c *storage.Client, bucketName string) error {
	// [START storage_enable_uniform_bucket_level_access]
	ctx := context.Background()

	bucket := c.Bucket(bucketName)
	enableUniformBucketLevelAccess := storage.BucketAttrsToUpdate{
		UniformBucketLevelAccess: &storage.UniformBucketLevelAccess{
			Enabled: true,
		},
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if _, err := bucket.Update(ctx, enableUniformBucketLevelAccess); err != nil {
		return err
	}
	// [END storage_enable_uniform_bucket_level_access]
	return nil
}

func disableUniformBucketLevelAccess(c *storage.Client, bucketName string) error {
	// [START storage_disable_uniform_bucket_level_access]
	ctx := context.Background()

	bucket := c.Bucket(bucketName)
	disableUniformBucketLevelAccess := storage.BucketAttrsToUpdate{
		UniformBucketLevelAccess: &storage.UniformBucketLevelAccess{
			Enabled: false,
		},
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if _, err := bucket.Update(ctx, disableUniformBucketLevelAccess); err != nil {
		return err
	}
	// [END storage_disable_uniform_bucket_level_access]
	return nil
}

func getUniformBucketLevelAccess(c *storage.Client, bucketName string) (*storage.BucketAttrs, error) {
	// [START storage_get_uniform_bucket_level_access]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	attrs, err := c.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return nil, err
	}
	uniformBucketLevelAccess := attrs.UniformBucketLevelAccess
	if uniformBucketLevelAccess.Enabled {
		log.Printf("Uniform bucket-level access is enabled for %q.\n",
			attrs.Name)
		log.Printf("Bucket will be locked on %q.\n",
			uniformBucketLevelAccess.LockedTime)
	} else {
		log.Printf("Uniform bucket-level access is not enabled for %q.\n",
			attrs.Name)
	}

	// [END storage_get_uniform_bucket_level_access]
	return attrs, nil
}
