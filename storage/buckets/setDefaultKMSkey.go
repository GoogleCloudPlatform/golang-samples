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

// [START storage_set_bucket_default_kms_key]
import (
	"context"

	"cloud.google.com/go/storage"
)

func setDefaultKMSkey(c *storage.Client, bucketName string, keyName string) error {
	ctx := context.Background()

	bucket := c.Bucket(bucketName)
	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		Encryption: &storage.BucketEncryption{DefaultKMSKeyName: keyName},
	}
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return err
	}
	return nil
}

// [END storage_set_bucket_default_kms_key]
