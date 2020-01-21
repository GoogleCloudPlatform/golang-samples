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

// Sample addDefaultBucketACL demonstrates adding default ACL to a single bucket.
package acl

// [START bucket_add_default_acl]
import (
	"context"

	"cloud.google.com/go/storage"
)

// addDefaultBucketACL adds default ACL to a bucket.
func addDefaultBucketACL(bucket string) error {
	// bucket := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	acl := client.Bucket(bucket).DefaultObjectACL()
	if err := acl.Set(ctx, storage.AllAuthenticatedUsers, storage.RoleReader); err != nil {
		return err
	}
	return nil
}

// [END bucket_add_default_acl]
