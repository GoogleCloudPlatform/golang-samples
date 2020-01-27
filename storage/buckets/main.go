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
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func main() {
	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		fmt.Fprintf(os.Stderr, "GOOGLE_CLOUD_PROJECT environment variable must be set.\n")
		os.Exit(1)
	}

	// [START setup]
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// [END setup]

	// Give the bucket a unique name.
	name := fmt.Sprintf("golang-example-buckets-%d", time.Now().Unix())
	if err := create(client, projectID, name); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created bucket: %v\n", name)

	// list buckets from the project
	buckets, err := list(client, projectID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("buckets: %+v\n", buckets)

	// get IAM policy
	if _, err := getPolicy(client, name); err != nil {
		log.Fatal(err)
	}

	// add user to IAM policy
	if err := addUser(client, name); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("added user to bucket %s", name)

	// get IAM policy
	if _, err := getPolicy(client, name); err != nil {
		log.Fatal(err)
	}

	// remove user from IAM policy
	if err := removeUser(client, name); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("removed user from bucket %s", name)

	// get IAM policy
	if _, err := getPolicy(client, name); err != nil {
		log.Fatal(err)
	}

	// delete the bucket
	if err := deleteBucket(client, name); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted bucket: %v\n", name)
}

func create(client *storage.Client, projectID, bucketName string) error {
	// [START create_bucket]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := client.Bucket(bucketName).Create(ctx, projectID, nil); err != nil {
		return err
	}
	// [END create_bucket]
	return nil
}

func createWithAttrs(client *storage.Client, projectID, bucketName string) error {
	// [START create_bucket_with_storageclass_and_location]
	ctx := context.Background()
	bucket := client.Bucket(bucketName)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := bucket.Create(ctx, projectID, &storage.BucketAttrs{
		StorageClass: "COLDLINE",
		Location:     "asia",
	}); err != nil {
		return err
	}
	// [END create_bucket_with_storageclass_and_location]
	return nil
}

func list(client *storage.Client, projectID string) ([]string, error) {
	// [START list_buckets]
	ctx := context.Background()

	var buckets []string
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	it := client.Buckets(ctx, projectID)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, battrs.Name)
	}
	// [END list_buckets]
	return buckets, nil
}

func deleteBucket(client *storage.Client, bucketName string) error {
	// [START delete_bucket]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := client.Bucket(bucketName).Delete(ctx); err != nil {
		return err
	}
	// [END delete_bucket]
	return nil
}

func getPolicy(c *storage.Client, bucketName string) (*iam.Policy, error) {
	// [START storage_get_bucket_policy]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	policy, err := c.Bucket(bucketName).IAM().Policy(ctx)
	if err != nil {
		return nil, err
	}
	for _, role := range policy.Roles() {
		log.Printf("%q: %q", role, policy.Members(role))
	}
	// [END storage_get_bucket_policy]
	return policy, nil
}

func addUser(c *storage.Client, bucketName string) error {
	// [START add_bucket_iam_member]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	bucket := c.Bucket(bucketName)
	policy, err := bucket.IAM().Policy(ctx)
	if err != nil {
		return err
	}
	// Other valid prefixes are "serviceAccount:", "user:"
	// See the documentation for more values.
	// https://cloud.google.com/storage/docs/access-control/iam
	policy.Add("group:cloud-logs@google.com", "roles/storage.objectViewer")
	if err := bucket.IAM().SetPolicy(ctx, policy); err != nil {
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
	policy, err := bucket.IAM().Policy(ctx)
	if err != nil {
		return err
	}
	// Other valid prefixes are "serviceAccount:", "user:"
	// See the documentation for more values.
	// https://cloud.google.com/storage/docs/access-control/iam
	policy.Remove("group:cloud-logs@google.com", "roles/storage.objectViewer")
	if err := bucket.IAM().SetPolicy(ctx, policy); err != nil {
		return err
	}
	// NOTE: It may be necessary to retry this operation if IAM policies are
	// being modified concurrently. SetPolicy will return an error if the policy
	// was modified since it was retrieved.
	// [END remove_bucket_iam_member]
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

func enableRequesterPays(c *storage.Client, bucketName string) error {
	// [START enable_requester_pays]
	ctx := context.Background()

	bucket := c.Bucket(bucketName)
	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		RequesterPays: true,
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return err
	}
	// [END enable_requester_pays]
	return nil
}

func disableRequesterPays(c *storage.Client, bucketName string) error {
	// [START disable_requester_pays]
	ctx := context.Background()

	bucket := c.Bucket(bucketName)
	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		RequesterPays: false,
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return err
	}
	// [END disable_requester_pays]
	return nil
}

func checkRequesterPays(c *storage.Client, bucketName string) error {
	// [START get_requester_pays_status]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	attrs, err := c.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return err
	}
	log.Printf("Is requester pays enabled? %v\n", attrs.RequesterPays)
	// [END get_requester_pays_status]
	return nil
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
