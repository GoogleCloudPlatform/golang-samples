// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample buckets creates a bucket, lists buckets and deletes a bucket
// using the Google Storage API. More documentation is available at
// https://cloud.google.com/storage/docs/json_api/v1/.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"

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
	if err := delete(client, name); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted bucket: %v\n", name)
}

func create(client *storage.Client, projectID, bucketName string) error {
	ctx := context.Background()
	// [START create_bucket]
	if err := client.Bucket(bucketName).Create(ctx, projectID, nil); err != nil {
		return err
	}
	// [END create_bucket]
	return nil
}

func createWithAttrs(client *storage.Client, projectID, bucketName string) error {
	ctx := context.Background()
	// [START create_bucket_with_storageclass_and_location]
	bucket := client.Bucket(bucketName)
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
	ctx := context.Background()
	// [START list_buckets]
	var buckets []string
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

func delete(client *storage.Client, bucketName string) error {
	ctx := context.Background()
	// [START delete_bucket]
	if err := client.Bucket(bucketName).Delete(ctx); err != nil {
		return err
	}
	// [END delete_bucket]
	return nil
}

func getPolicy(c *storage.Client, bucketName string) (*iam.Policy, error) {
	ctx := context.Background()

	// [START storage_get_bucket_policy]
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
	ctx := context.Background()

	// [START add_bucket_iam_member]
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
	ctx := context.Background()

	// [START remove_bucket_iam_member]
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

func enableRequesterPays(c *storage.Client, bucketName string) error {
	ctx := context.Background()

	// [START enable_requester_pays]
	bucket := c.Bucket(bucketName)
	if _, err := bucket.Update(ctx, storage.BucketAttrsToUpdate{
		RequesterPays: true,
	}); err != nil {
		return err
	}
	// [END enable_requester_pays]
	return nil
}

func disableRequesterPays(c *storage.Client, bucketName string) error {
	ctx := context.Background()

	// [START disable_requester_pays]
	bucket := c.Bucket(bucketName)
	if _, err := bucket.Update(ctx, storage.BucketAttrsToUpdate{
		RequesterPays: false,
	}); err != nil {
		return err
	}
	// [END disable_requester_pays]
	return nil
}

func checkRequesterPays(c *storage.Client, bucketName string) error {
	ctx := context.Background()

	// [START get_requester_pays_status]
	attrs, err := c.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Is requester pays enabled? %v\n", attrs.RequesterPays)
	// [END get_requester_pays_status]
	return nil
}
