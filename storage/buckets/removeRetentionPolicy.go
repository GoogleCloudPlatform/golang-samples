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
package buckets

// [START storage_remove_retention_policy]
import (
	"context"
	"errors"

	"cloud.google.com/go/storage"
)

func removeRetentionPolicy(c *storage.Client, bucketName string) error {
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
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return err
	}
	return nil
}

// [END storage_remove_retention_policy]
