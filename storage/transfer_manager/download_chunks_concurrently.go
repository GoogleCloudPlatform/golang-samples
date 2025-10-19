// Copyright 2025 Google LLC
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

// [START storage_download_chunks_concurrently]
package transfermanager

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"cloud.google.com/go/storage/transfermanager"
)

// downloadChunksConcurrently downloads a single file in chunks, concurrently in a process pool.
func downloadChunksConcurrently(w io.Writer, bucketName, blobName, filename string) error {
	// bucketName := "your-bucket-name"
	// blobName := "target-file"
	// filename := ""

	// The size of each chunk. The performance impact of this value depends on
	// the use case. The remote service has a minimum of 5 MiB and a
	// maximum of 5 GiB.
	chunkSize := 32 * 1024 * 1024 // 32 MiB

	// The maximum number of processes to use for the operation. The performance
	// impact of this value depends on the use case, but smaller files usually
	// benefit from a higher number of processes. Each additional process
	// occupies some CPU and memory resources until finished.
	workers := 8

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	d, err := transfermanager.NewDownloader(client, transfermanager.WithPartSize(int64(chunkSize)), transfermanager.WithWorkers(workers))
	if err != nil {
		return fmt.Errorf("transfermanager.NewDownloader: %w", err)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}

	in := &transfermanager.DownloadObjectInput{
		Bucket:      bucketName,
		Object:      blobName,
		Destination: f,
	}

	if err := d.DownloadObject(ctx, in); err != nil {
		return fmt.Errorf("d.DownloadObject: %w", err)
	}

	results, err := d.WaitAndClose()
	if err != nil {
		return fmt.Errorf("d.WaitAndClose: %w", err)
	}

	// Iterate through completed downloads and process results.
	for _, out := range results {
		if out.Err != nil {
			fmt.Fprintf(w, "download of %v failed with error %v\n", out.Object, out.Err)
		} else {
			fmt.Fprintf(w, "Downloaded %v to %v.\n", blobName, filename)
		}
	}
	return nil
}

// [END storage_download_chunks_concurrently]
