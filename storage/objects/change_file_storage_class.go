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

// [START storage_change_file_storage_class]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// changeObjectStorageClass changes the storage class of a single object.
func changeObjectStorageClass(w io.Writer, bucket, object string) error {
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

	o := client.Bucket(bucket).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to copy is aborted if the
	// object's generation number does not match your precondition.
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %w", err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// See the StorageClass documentation for other valid storage classes:
	// https://cloud.google.com/storage/docs/storage-classes
	newStorageClass := "COLDLINE"
	// You can't change an object's storage class directly, the only way is
	// to rewrite the object with the desired storage class.
	copier := o.CopierFrom(o)
	copier.StorageClass = newStorageClass
	if _, err := copier.Run(ctx); err != nil {
		return fmt.Errorf("copier.Run: %w", err)
	}
	fmt.Fprintf(w, "Object %v in bucket %v had its storage class set to %v\n", object, bucket, newStorageClass)
	return nil
}

// [END storage_change_file_storage_class]
