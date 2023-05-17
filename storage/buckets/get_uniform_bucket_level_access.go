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

// [START storage_get_uniform_bucket_level_access]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// getUniformBucketLevelAccess gets uniform bucket-level access.
func getUniformBucketLevelAccess(w io.Writer, bucketName string) (*storage.BucketAttrs, error) {
	// bucketName := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("Bucket(%q).Attrs: %w", bucketName, err)
	}
	uniformBucketLevelAccess := attrs.UniformBucketLevelAccess
	if uniformBucketLevelAccess.Enabled {
		fmt.Fprintf(w, "Uniform bucket-level access is enabled for %q.\n", attrs.Name)
		fmt.Fprintf(w, "Bucket will be locked on %q.\n", uniformBucketLevelAccess.LockedTime)
	} else {
		fmt.Fprintf(w, "Uniform bucket-level access is not enabled for %q.\n", attrs.Name)
	}
	return attrs, nil
}

// [END storage_get_uniform_bucket_level_access]
