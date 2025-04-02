// Copyright 2025 Google LLC
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

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// restoreSoftDeletedObject restores a soft-deleted object in a bucket.
func restoreSoftDeletedObject(w io.Writer, bucketName, objectName string, generation int64) error {
	// bucketName := "your-bucket-name"
	// objectName := "your-object-name"
	// generation := 1234567890 // The generation of the soft-deleted object.

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// Restore the soft-deleted object.
	objHandle := client.Bucket(bucketName).Object(objectName).Generation(generation)
	attrs, err := objHandle.Restore(ctx, &storage.RestoreOptions{
		CopySourceACL: true,
	})
	if err != nil {
		return fmt.Errorf("Object(%q).Restore: %w", objectName, err)
	}

	fmt.Fprintf(w, "Soft-deleted object %q (generation: %d) has been restored in bucket %q.\n", objectName, generation, bucketName)
	fmt.Fprintf(w, "Updated object attributes: %+v\n", attrs)
	return nil
}
