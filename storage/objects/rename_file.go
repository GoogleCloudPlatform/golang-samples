// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package objects contains samples for working with Google Cloud Storage objects.
package objects

// [START storage_rename_file]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// renameFile renames a file in a Google Cloud Storage bucket.
// This operation is only available for buckets with the Hierarchical Namespace
// feature enabled.
func renameFile(w io.Writer, bucket, srcObject, destObject string) error {
	// bucket := "my-bucket"
	// srcObject := "path/to/source/object"
	// destObject := "path/to/destination/object"

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	src := client.Bucket(bucket).Object(srcObject)
	dst := client.Bucket(bucket).Object(destObject)

	// This operation is only available for buckets with the Hierarchical
	// Namespace feature enabled.
	mover := dst.MoverFrom(src)
	if _, err := mover.Run(ctx); err != nil {
		return fmt.Errorf("Object(%q).MoverFrom(%q).Run: %w", destObject, srcObject, err)
	}
	fmt.Fprintf(w, "Renamed object %s to %s\n", srcObject, destObject)
	return nil
}

// [END storage_rename_file]
