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

// [START storage_copy_file_archived_generation]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// copyFileArchivedGeneration copies a noncurrent version of an object.
func copyFileArchivedGeneration(w io.Writer, dstBucket, srcBucket, srcObject string, gen int64) error {
	// dstBucket := "destination-bucket-name"
	// srcBucket := "source-bucket-name"
	// srcObject := "source-object-name"
	// gen := 1587012235914578
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	dstObject := fmt.Sprintf("%v-copy-%v", srcObject, gen)
	src := client.Bucket(srcBucket).Object(srcObject)
	dst := client.Bucket(dstBucket).Object(dstObject)

	if _, err := dst.CopierFrom(src.Generation(gen)).Run(ctx); err != nil {
		return fmt.Errorf("Object(%q).CopierFrom(%q).Generation(%v).Run: %v", dstObject, srcObject, gen, err)
	}
	fmt.Fprintf(w, "Generation %v of object %v in bucket %v was copied to %v in bucket %v\n", gen, srcObject, srcBucket, dstObject, dstBucket)
	return nil
}

// [END storage_copy_file_archived_generation]
