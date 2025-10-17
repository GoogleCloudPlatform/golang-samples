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

// [START storage_pause_and_resume_appendable_upload]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"cloud.google.com/go/storage/experimental"
)

// pauseAndResumeAppendableUpload creates a new unfinalized appendable object,
// closes the Writer, then re-opens the object for writing using
// NewWriterFromAppendableObject.
func pauseAndResumeAppendableUpload(w io.Writer, bucket, object string) error {
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

	// Start an appendable upload and write some data.
	writer := client.Bucket(bucket).Object(object).NewWriter(ctx)

	if _, err := writer.Write([]byte("Some data\n")); err != nil {
		return fmt.Errorf("Writer.Write: %w", err)
	}

	// The writer is closed, but the upload is not finalized. This "pauses" the
	// upload, as the object remains appendable.
	if err := writer.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	fmt.Fprintf(w, "Uploaded unfinalize object %v\n", object)

	// To resume the upload we need the object's generation. We can get this
	// from the previous Writer after close.
	gen := writer.Attrs().Generation

	// Now resume the upload. Writer options including finalization can be
	// passed on calling this constructor.
	appendWriter, offset, err := client.Bucket(bucket).Object(object).Generation(gen).NewWriterFromAppendableObject(
		ctx, &storage.AppendableWriterOpts{
			FinalizeOnClose: true,
		},
	)
	if err != nil {
		return fmt.Errorf("NewWriterFromAppendableObject: %v", err)
	}
	fmt.Fprintf(w, "Resuming upload from offset %v\n", offset)

	// Append the rest of the data and close the Writer to finalize.
	if _, err := appendWriter.Write([]byte("resumed data\n")); err != nil {
		return fmt.Errorf("appendWriter.Write: %v", err)
	}
	if err := appendWriter.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	fmt.Fprintf(w, "Uploaded and finalized object %v\n", object)
	return nil
}

// [END storage_pause_and_resume_appendable_upload]
