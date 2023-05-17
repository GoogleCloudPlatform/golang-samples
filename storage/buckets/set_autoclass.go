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
// See https://cloud.google.com/storage/docs/using-autoclass for more information.

// Note: Only update requests that disable Autoclass are currently supported.
// To enable Autoclass, you must set it at bucket creation time.
func setAutoclass(w io.Writer, bucketName string, value bool) error {
	// bucketName := "bucket-name"
	// value := false
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
			Enabled: value,
		},
	}
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return fmt.Errorf("Bucket(%q).Update: %w", bucketName, err)
	}
	fmt.Fprintf(w, "Autoclass enabled was set to %v on bucket %q \n", bucketAttrsToUpdate.Autoclass.Enabled, bucketName)
	return nil
}

// [END storage_set_autoclass]
