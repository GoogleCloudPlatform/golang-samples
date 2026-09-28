// Copyright 2026 Google LLC
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

package bidi

// [START storage_optimize_write_latency_pool]
import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// optimizeWriteLatencyPool uses a pre-warmed pool of writers for a zonal bucket.
func optimizeWriteLatencyPool(out io.Writer, bucketName, keyPrefix string) error {
	// bucketName := "bucket-name"
	// keyPrefix := "pooled-object"
	ctx := context.Background()
	poolSize := 3

	client, err := storage.NewGRPCClient(ctx, storage.WithAppendableUploads())
	if err != nil {
		return fmt.Errorf("storage.NewGRPCClient: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	newPrewarmedWriter := func(name string) (*storage.Writer, error) {
		w := bucket.Object(name).If(storage.Conditions{DoesNotExist: true}).NewWriter(ctx)
		// TODO(https://github.com/googleapis/google-cloud-go/issues/20580): Flush may
		// no longer be necessary in the future.
		if _, err := w.Flush(); err != nil {
			_ = w.Close()
			return nil, fmt.Errorf("Writer.Flush(%s): %w", name, err)
		}
		return w, nil
	}

	// 1. Init pool: A buffered channel is a thread-safe FIFO queue of
	// pre-warmed writers. Flushing incurs operation charges, so size the pool
	// carefully.
	pool := make(chan *storage.Writer, poolSize)
	// On return, close pooled writers; unused objects remain unfinalized.
	defer func() {
		close(pool)
		for w := range pool {
			_ = w.Close()
		}
	}()

	for i := 0; i < poolSize; i++ {
		w, err := newPrewarmedWriter(fmt.Sprintf("%s_%d", keyPrefix, i))
		if err != nil {
			return err
		}
		pool <- w
	}

	// 2. Write: Take a pre-warmed writer (waiting for a refill if the pool is
	// empty) and commit with Flush() (~1-2 ms) instead of blocking on Close().
	var w *storage.Writer
	select {
	case w = <-pool:
	case <-time.After(10 * time.Second):
		return errors.New("timed out waiting for a pre-warmed writer")
	}
	if _, err := w.Write([]byte("0123456789")); err != nil {
		_ = w.Close()
		return fmt.Errorf("Writer.Write: %w", err)
	}
	if _, err := w.Flush(); err != nil {
		_ = w.Close()
		return fmt.Errorf("Writer.Flush: %w", err)
	}

	// 3. Pool maintenance (run asynchronously off the critical write path):
	// Close the used writer without finalizing and refill the pool.
	maintenance := make(chan struct{})
	go func(used *storage.Writer) {
		defer close(maintenance)
		_ = used.Close()
		next, err := newPrewarmedWriter(fmt.Sprintf("%s_%d", keyPrefix, poolSize))
		if err != nil {
			fmt.Fprintf(out, "failed to pre-warm replacement writer: %v\n", err)
			return
		}
		pool <- next
	}(w)
	// Wait for maintenance on return; runs before the pool cleanup above.
	defer func() { <-maintenance }()

	// 4. Read: Unfinalized objects are readable after Flush().
	r, err := bucket.Object(fmt.Sprintf("%s_0", keyPrefix)).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object.NewReader: %w", err)
	}
	contents, err := io.ReadAll(r)
	_ = r.Close()
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	fmt.Fprintf(out, "Read unfinalized object %s_0: %s\n", keyPrefix, string(contents))
	return nil
}

// [END storage_optimize_write_latency_pool]
