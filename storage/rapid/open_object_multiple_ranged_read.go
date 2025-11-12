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

// [START storage_open_object_multiple_ranged_read]
import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"cloud.google.com/go/storage/experimental"
)

// openObjectMultipleRangedRead opens a single object using
// MultiRangeDownloader to download multiple ranges.
func openObjectMultipleRangedRead(w io.Writer, bucket, object string) ([][]byte, error) {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewGRPCClient(ctx, experimental.WithZonalBucketAPIs())
	if err != nil {
		return nil, fmt.Errorf("storage.NewGRPCClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// Create the MultiRangeDownloader, which opens a stream to the object.
	mrd, err := client.Bucket(bucket).Object(object).NewMultiRangeDownloader(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewMultiRangeDownloader: %w", err)
	}

	// Add some 1 KiB ranges to download. This call is non-blocking. The
	// provided callback is invoked when the range download is complete.
	startOffsets := []int64{0, 1024, 2048}
	var dataBufs [3]bytes.Buffer
	var errs []error
	for i, off := range startOffsets {
		mrd.Add(&dataBufs[i], off, 1024, func(off, length int64, err error) {
			if err != nil {
				errs = append(errs, err)
			} else {
				fmt.Fprintf(w, "downloaded range at offset %v", off)
			}
		})
	}

	// Wait for all downloads to complete.
	mrd.Wait()
	if len(errs) > 0 {
		return nil, fmt.Errorf("one or more downloads failed; errors: %v", errs)
	}
	if err := mrd.Close(); err != nil {
		return nil, fmt.Errorf("MultiRangeDownloader.Close: %w", err)
	}

	fmt.Fprintf(w, "Read the ranges of %v into memory\n", object)

	// Collect the byte slices
	var byteSlices [][]byte
	for _, buf := range dataBufs {
		byteSlices = append(byteSlices, buf.Bytes())
	}

	return byteSlices, nil
}

// [END storage_open_object_multiple_ranged_read]
