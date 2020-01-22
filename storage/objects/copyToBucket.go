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
package objects

// [START copy_file]
import (
	"context"

	"cloud.google.com/go/storage"
)

// copyToBucket copies an object into specified bucket.
func copyToBucket(dstBucket, srcBucket, srcObject string) error {
	// dstBucket := "bucket-1"
	// srcBucket := "bucket-2"
	// srcObject := "object"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	dstObject := srcObject + "-copy"
	src := client.Bucket(srcBucket).Object(srcObject)
	dst := client.Bucket(dstBucket).Object(dstObject)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}
	return nil
}

// [END copy_file]
