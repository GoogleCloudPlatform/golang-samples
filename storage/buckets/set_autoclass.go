// Copyright 2022 Google LLC
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

// [START storage_set_autoclass]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// setAutoclass sets the Autoclass configuration for a bucket.
func setAutoclass(w io.Writer, bucketName string) error {
	// bucketName := "bucket-name"

	// To update the configuration for Autoclass.TerminalStorageClass,
	// Autoclass.Enabled must also be set to true.
	// To disable autoclass on the bucket, set to an empty &Autoclass{}.
	enabled := true
	terminalStorageClass := "ARCHIVE"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	bucket := client.Bucket(bucketName)
	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		Autoclass: &storage.Autoclass{
			Enabled:              enabled,
			TerminalStorageClass: terminalStorageClass,
		},
	}
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return fmt.Errorf("Bucket(%q).Update: %w", bucketName, err)
	}
	fmt.Fprintf(w, "Autoclass enabled was set to %v on bucket %q \n", bucketAttrsToUpdate.Autoclass.Enabled, bucketName)
	fmt.Fprintf(w, "Autoclass terminal storage class was last updated to %v at %v", bucketAttrsToUpdate.Autoclass.TerminalStorageClass, bucketAttrsToUpdate.Autoclass.TerminalStorageClassUpdateTime)
	return nil
}

// [END storage_set_autoclass]
