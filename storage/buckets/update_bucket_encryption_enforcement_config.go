// Copyright 2026 Google LLC
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

// [START storage_update_bucket_encryption_enforcement_config]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// updateBucketEncryptionEnforcementConfig updates a bucket's encryption enforcement configuration.
func updateBucketEncryptionEnforcementConfig(w io.Writer, bucketName string) error {
	// bucketName := "bucket-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	bucket := client.Bucket(bucketName)
	if _, err := bucket.Update(ctx, storage.BucketAttrsToUpdate{
		Encryption: &storage.BucketEncryption{
			GoogleManagedEncryptionEnforcementConfig: &storage.EncryptionEnforcementConfig{
				RestrictionMode: storage.NotRestricted,
			},
			CustomerManagedEncryptionEnforcementConfig: &storage.EncryptionEnforcementConfig{
				RestrictionMode: storage.FullyRestricted,
			},
			CustomerSuppliedEncryptionEnforcementConfig: &storage.EncryptionEnforcementConfig{
				RestrictionMode: storage.FullyRestricted,
			},
		},
	}); err != nil {
		return fmt.Errorf("Bucket(%q).Update: %w", bucketName, err)
	}
	fmt.Fprintf(w, "Bucket %v encryption enforcement policies updated.\n", bucketName)
	return nil
}

// [END storage_update_bucket_encryption_enforcement_config]
