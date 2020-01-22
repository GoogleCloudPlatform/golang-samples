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

// [START storage_disable_uniform_bucket_level_access]
import (
	"context"

	"cloud.google.com/go/storage"
)

// disableUniformBucketLevelAccess sets uniform bucket-level access to false.
func disableUniformBucketLevelAccess(bucketName string) error {
	// bucketName := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	disableUniformBucketLevelAccess := storage.BucketAttrsToUpdate{
		UniformBucketLevelAccess: &storage.UniformBucketLevelAccess{
			Enabled: false,
		},
	}
	if _, err := bucket.Update(ctx, disableUniformBucketLevelAccess); err != nil {
		return err
	}
	return nil
}

// [END storage_disable_uniform_bucket_level_access]
