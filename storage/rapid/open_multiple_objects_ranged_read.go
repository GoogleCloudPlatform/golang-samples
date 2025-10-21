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

// [START storage_open_multiple_objects_ranged_read]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"cloud.google.com/go/storage/experimental"
	"cloud.google.com/go/storage/transfermanager"
)

// openMultipleObjectsRangedRead reads ranges from multiple objects using
// transfer manager to download ranges in parallel.
func openMultipleObjectsRangedRead(w io.Writer, bucket string, objects []string) ([][]byte, error) {
	// bucket := "bucket-name"
	// objects := []string{"object-name1", "object-name2"}
	ctx := context.Background()
	client, err := storage.NewGRPCClient(ctx, experimental.WithZonalBucketAPIs())
	if err != nil {
		return nil, fmt.Errorf("storage.NewGRPCClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	// Initialize a transfermanager Downloader and add individual objects, downloading
	// the first KiB of each.
	databufs := make([]*transfermanager.DownloadBuffer, 0, len(objects))
	d, err := transfermanager.NewDownloader(client, transfermanager.WithWorkers(16))
	if err != nil {
		return nil, fmt.Errorf("NewDownloader: %w", err)
	}
	for _, obj := range objects {
		b := make([]byte, 1024)
		buf := transfermanager.NewDownloadBuffer(b)
		databufs = append(databufs, buf)
		d.DownloadObject(ctx, &transfermanager.DownloadObjectInput{
			Bucket:      bucket,
			Object:      obj,
			Destination: buf,
			Range: &transfermanager.DownloadRange{
				Offset: 0,
				Length: 1024,
			},
		})
	}

	// Wait for all Download jobs to complete.
	outs, err := d.WaitAndClose()
	if err != nil {
		return nil, fmt.Errorf("downloading: %w", err)
	}
	for _, out := range outs {
		fmt.Fprintf(w, "Downloaded 1 KiB of %v from bucket %v\n", out.Object, out.Bucket)
	}
	var byteSlices [][]byte
	for _, buf := range databufs {
		byteSlices = append(byteSlices, buf.Bytes())
	}
	return byteSlices, nil
}

// [END storage_open_multiple_objects_ranged_read]
