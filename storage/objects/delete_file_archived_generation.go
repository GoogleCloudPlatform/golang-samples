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

// [START storage_delete_file_archived_generation]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// deleteOldVersionOfObject deletes a noncurrent version of an object.
func deleteOldVersionOfObject(w io.Writer, bucketName, objectName string, gen int64) error {
	// bucketName := "bucket-name"
	// objectName := "object-name"

	// gen is the generation of objectName to delete.
	// gen := 1587012235914578
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	obj := client.Bucket(bucketName).Object(objectName)
	if err := obj.Generation(gen).Delete(ctx); err != nil {
		return fmt.Errorf("Bucket(%q).Object(%q).Generation(%v).Delete: %w", bucketName, objectName, gen, err)
	}
	fmt.Fprintf(w, "Generation %v of object %v was deleted from %v\n", gen, objectName, bucketName)
	return nil
}

// [END storage_delete_file_archived_generation]
