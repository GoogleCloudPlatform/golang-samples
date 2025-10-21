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

// [START storage_read_appendable_object_tail]
import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"cloud.google.com/go/storage/experimental"
)

// readAppendableObjectTail simulates a "tail -f" command on a GCS object. It
// repeatedly polls an appendable object for new content. In a real
// application, the object would be written to by a separate process.
func readAppendableObjectTail(w io.Writer, bucket, object string) ([]byte, error) {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewGRPCClient(ctx, experimental.WithZonalBucketAPIs())
	if err != nil {
		return nil, fmt.Errorf("storage.NewGRPCClient: %w", err)
	}
	defer client.Close()

	// Set a context timeout. When this timeout is reached, the read stream
	// will be closed, so omit this to tail indefinitely.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Create a new appendable object and write some data.
	writer := client.Bucket(bucket).Object(object).If(storage.Conditions{DoesNotExist: true}).NewWriter(ctx)
	if _, err := writer.Write([]byte("Some data\n")); err != nil {
		return nil, fmt.Errorf("Writer.Write: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("Writer.Close: %w", err)
	}
	gen := writer.Attrs().Generation

	// Create the MultiRangeDownloader, which opens a read stream to the object.
	mrd, err := client.Bucket(bucket).Object(object).NewMultiRangeDownloader(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewMultiRangeDownloader: %w", err)
	}

	// In a goroutine, poll the object. In this example we continue until all the
	// bytes we expect to see were received, but in a real application this could
	// continue to poll indefinitely until some signal is received.
	var buf bytes.Buffer
	var mrdErr error
	done := make(chan bool)
	go func() {
		var currOff int64
		rangeDownloaded := make(chan bool)
		for buf.Len() < 100 {
			// Add the current range and wait for it to be downloaded.
			// Using a length of 0 will read to the current end of the object.
			// The callback will give the actual number of bytes that were
			// read in each iteration.
			mrd.Add(&buf, currOff, 0, func(offset, length int64, err error) {
				// After each range is received, update
				// the starting offset based on how many bytes were received.
				if err != nil {
					mrdErr = err
				}
				currOff += length
				rangeDownloaded <- true
			})
			<-rangeDownloaded
			if mrdErr != nil {
				break
			}
			time.Sleep(1 * time.Second)
		}
		// After exiting the loop, close MultiRangeDownloader and signal that
		// all ranges have been read.
		if err := mrd.Close(); err != nil {
			mrdErr = err
		}
		done <- true
	}()

	// Meanwhile, continue to write 10 bytes at a time to the object.
	// This could be done by calling NewWriterFromAppendable object repeatedly
	// (as in the example) or calling Writer.Flush without closing the Writer.
	for range 9 {
		appendWriter, offset, err := client.Bucket(bucket).Object(object).Generation(gen).NewWriterFromAppendableObject(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("NewWriterFromAppendableObject: %w", err)
		}
		if _, err := appendWriter.Write([]byte("more data\n")); err != nil {
			return nil, fmt.Errorf("appendWriter.Write: %w", err)
		}
		if err := appendWriter.Close(); err != nil {
			return nil, fmt.Errorf("appendWriter.Close: %w", err)
		}
		fmt.Fprintf(w, "Wrote 10 bytes at offset %v", offset)
	}

	// Wait for tailing goroutine to exit.
	<-done
	if mrdErr != nil {
		return nil, fmt.Errorf("MultiRangeDownloader: %w", err)
	}
	fmt.Fprintf(w, "Read %v bytes from object %v", buf.Len(), object)
	return buf.Bytes(), nil
}

// [END storage_read_appendable_object_tail]
