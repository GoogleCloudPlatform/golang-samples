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

package rapid

// [START storage_finalize_appendable_object_upload]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"cloud.google.com/go/storage/experimental"
)

// finalizeAppendableObject creates, uploads and finalizes a new object in
// a rapid bucket.
func finalizeAppendableObject(w io.Writer, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewGRPCClient(ctx, experimental.WithZonalBucketAPIs())
	if err != nil {
		return fmt.Errorf("storage.NewGRPCClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// Create a Writer and set FinalizeOnClose so that the object will be
	// finalized after the write is complete.
	writer := client.Bucket(bucket).Object(object).NewWriter(ctx)
	writer.FinalizeOnClose = true

	if _, err := writer.Write([]byte("some data to finalize\n")); err != nil {
		return fmt.Errorf("Writer.Write: %w", err)
	}

	// Close the Writer to flush any remaining buffered data and finalize
	// the upload. This makes the object non-appendable.
	// No more data can be written to this object.
	if err := writer.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	fmt.Fprintf(w, "Uploaded and finalized object %v", object)

	return nil
}

// [END storage_finalize_appendable_object_upload]
