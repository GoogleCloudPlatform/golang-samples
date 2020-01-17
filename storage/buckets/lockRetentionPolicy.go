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

// [START storage_lock_retention_policy]
import (
	"context"
	"log"

	"cloud.google.com/go/storage"
)

func lockRetentionPolicy(c *storage.Client, bucketName string) error {
	ctx := context.Background()

	bucket := c.Bucket(bucketName)
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
	return nil
}

// [END storage_lock_retention_policy]
