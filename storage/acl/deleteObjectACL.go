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

// [START delete_object_default_acl]
import (
	"context"

	"cloud.google.com/go/storage"
)

// deleteObjectACL removes default ACL from object.
func deleteObjectACL(bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	acl := client.Bucket(bucket).Object(object).ACL()
	if err := acl.Delete(ctx, storage.AllAuthenticatedUsers); err != nil {
		return err
	}
	return nil
}

// [END delete_object_default_acl]
