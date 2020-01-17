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

// [START storage_get_retention_policy]
import (
	"context"
	"log"

	"cloud.google.com/go/storage"
)

func getRetentionPolicy(c *storage.Client, bucketName string) (*storage.BucketAttrs, error) {
	ctx := context.Background()

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
	return attrs, nil
}

// [END storage_get_retention_policy]
