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

// [START storage_move_file]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// moveFile moves an object into another location.
func moveFile(w io.Writer, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	dstName := object + "-rename"
	src := client.Bucket(bucket).Object(object)
	dst := client.Bucket(bucket).Object(dstName)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to copy the file is aborted
	// if the object's generation number does not match your precondition.
	// For a dst object that does not yet exist, set the DoesNotExist precondition.
	dst = dst.If(storage.Conditions{DoesNotExist: true})
	// If the destination object already exists in your bucket, set instead a
	// generation-match precondition using its generation number.
	// attrs, err := dst.Attrs(ctx)
	// if err != nil {
	// 	return fmt.Errorf("object.Attrs: %w", err)
	// }
	// dst = dst.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return fmt.Errorf("Object(%q).CopierFrom(%q).Run: %w", dstName, object, err)
	}
	if err := src.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %w", object, err)
	}
	fmt.Fprintf(w, "Blob %v moved to %v.\n", object, dstName)
	return nil
}

// [END storage_move_file]
