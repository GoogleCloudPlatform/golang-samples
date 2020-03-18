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

package acl

// [START storage_remove_bucket_owner]
import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

// removeBucketOwner removes ACL from a bucket.
func removeBucketOwner(bucket string, entity storage.ACLEntity) error {
	// bucket := "bucket-name"
	// entity := storage.AllUsers
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	acl := client.Bucket(bucket).ACL()
	if err := acl.Delete(ctx, entity); err != nil {
		return fmt.Errorf("ACLHandle.Delete: %v", err)
	}
	return nil
}

// [END storage_remove_bucket_owner]
