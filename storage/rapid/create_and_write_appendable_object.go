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

// [START storage_create_and_write_appendable_object]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"cloud.google.com/go/storage/experimental"
)

// createAndWriteAppendableObject creates and uploads a new appendable object in
// a rapid bucket. The object will not be finalized.
func createAndWriteAppendableObject(w io.Writer, bucket, object string) error {
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

	// Create a Writer and write some data.
	writer := client.Bucket(bucket).Object(object).NewWriter(ctx)

	if _, err := writer.Write([]byte("Some data\n")); err != nil {
		return fmt.Errorf("Writer.Write: %w", err)
	}

	// Flush the buffered data to the service. This is not a terminal
	// operation. The Writer can be used after the flush completes.
	// After a flush, the data is visible to readers.
	size, err := writer.Flush()
	if err != nil {
		return fmt.Errorf("Writer.Flush: %w", err)
	}
	fmt.Fprintf(w, "Flush completed. Persisted size is now %d", size)

	// The Writer is still open. We can write more data.
	if _, err := writer.Write([]byte("Some more data\n")); err != nil {
		return fmt.Errorf("Writer.Write: %w", err)
	}

	// Close the Writer to flush any remaining buffered data.
	// The object will be unfinalized, which means another writer can
	// later append to the object.
	if err := writer.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	fmt.Fprintf(w, "Uploaded object %v", object)

	return nil
}

// [END storage_create_and_write_appendable_object]
