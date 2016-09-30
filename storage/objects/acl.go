// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/net/context"

	"cloud.google.com/go/storage"
)

// TODO(jbd): Add START and END tags once the names are finalized.

func addBucketACL(client *storage.Client, bucket string) error {
	ctx := context.Background()

	acl := client.Bucket(bucket).ACL()
	if err := acl.Set(ctx, storage.AllAuthenticatedUsers, storage.RoleReader); err != nil {
		return err
	}
	return nil
}

func addDefaultBucketACL(client *storage.Client, bucket string) error {
	ctx := context.Background()

	acl := client.Bucket(bucket).DefaultObjectACL()
	if err := acl.Set(ctx, storage.AllAuthenticatedUsers, storage.RoleReader); err != nil {
		return err
	}
	return nil
}

func bucketACL(client *storage.Client, bucket string) error {
	ctx := context.Background()

	rules, err := client.Bucket(bucket).ACL().List(ctx)
	if err != nil {
		return err
	}
	for _, rule := range rules {
		fmt.Printf("ACL rule: %v\n", rule)
	}
	return nil
}

func bucketACLFiltered(client *storage.Client, bucket string, entity storage.ACLEntity) error {
	ctx := context.Background()

	rules, err := client.Bucket(bucket).ACL().List(ctx)
	if err != nil {
		return err
	}
	for _, r := range rules {
		if r.Entity == entity {
			fmt.Printf("ACL rule role: %v\n", r.Role)
		}
	}
	return nil
}

func addObjectACL(client *storage.Client, bucket, object string) error {
	ctx := context.Background()

	acl := client.Bucket(bucket).Object(object).ACL()
	if err := acl.Set(ctx, storage.AllAuthenticatedUsers, storage.RoleReader); err != nil {
		return err
	}
	return nil
}

func objectACL(client *storage.Client, bucket, object string) error {
	ctx := context.Background()

	rules, err := client.Bucket(bucket).Object(object).ACL().List(ctx)
	if err != nil {
		return err
	}
	for _, rule := range rules {
		fmt.Printf("ACL rule: %v\n", rule)
	}
	return nil
}

func objectACLFiltered(client *storage.Client, bucket, object string, entity storage.ACLEntity) error {
	ctx := context.Background()

	rules, err := client.Bucket(bucket).ACL().List(ctx)
	if err != nil {
		return err
	}
	for _, r := range rules {
		if r.Entity == entity {
			fmt.Printf("ACL rule role: %v\n", r.Role)
		}
	}
	return nil
}

func deleteBucketACL(client *storage.Client, bucket string) error {
	ctx := context.Background()

	acl := client.Bucket(bucket).ACL()
	if err := acl.Delete(ctx, storage.AllAuthenticatedUsers); err != nil {
		return err
	}
	return nil
}

func deleteDefaultBucketACL(client *storage.Client, bucket string) error {
	ctx := context.Background()

	acl := client.Bucket(bucket).DefaultObjectACL()
	if err := acl.Delete(ctx, storage.AllAuthenticatedUsers); err != nil {
		return err
	}
	return nil
}

func deleteObjectACL(client *storage.Client, bucket, object string) error {
	ctx := context.Background()

	acl := client.Bucket(bucket).Object(object).ACL()
	if err := acl.Delete(ctx, storage.AllAuthenticatedUsers); err != nil {
		return err
	}
	return nil
}

func signedURL(client *storage.Client, bucket, object string) error {
	// Download a p12 service account private key from the Google Developers Console.
	// And convert it to PEM by running the command below:
	//	$ openssl pkcs12 -in key.p12 -passin pass:notasecret -out my-private-key.pem -nodes
	pkey, err := ioutil.ReadFile("my-private-key.pem")
	if err != nil {
		return err
	}
	url, err := storage.SignedURL(bucket, object, &storage.SignedURLOptions{
		GoogleAccessID: "xxx@developer.gserviceaccount.com",
		PrivateKey:     pkey,
		Method:         "GET",
		Expires:        time.Now().Add(48 * time.Hour),
	})
	if err != nil {
		return err
	}
	fmt.Println(url)
	return nil
}
