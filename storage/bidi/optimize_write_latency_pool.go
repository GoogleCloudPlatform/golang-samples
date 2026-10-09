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
	"sync/atomic"
	"time"

	"cloud.google.com/go/storage"
)

// optimizeWriteLatencyPool uses a pre-warmed pool of writers for appendable
// objects in a bucket with the Rapid storage class.
func optimizeWriteLatencyPool(out io.Writer, bucketName, keyPrefix string) error {
	// bucketName := "bucket-name"
	// keyPrefix := "pooled-object"
	const poolSize = 3
	// Pool writers live until ctx is canceled (here, on return).
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := storage.NewGRPCClient(ctx, storage.WithAppendableUploads())
	if err != nil {
		return fmt.Errorf("storage.NewGRPCClient: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	// Each pre-warmed writer gets the next unique object number.
	var nextID atomic.Int64
	newPrewarmedWriter := func() (*storage.Writer, error) {
		name := fmt.Sprintf("%s_%d", keyPrefix, nextID.Add(1)-1)
		w := bucket.Object(name).If(storage.Conditions{DoesNotExist: true}).NewWriter(ctx)
		if _, err := w.Flush(); err != nil {
			return nil, errors.Join(fmt.Errorf("Writer.Flush(%q): %w", name, err), w.Close())
		}
		return w, nil
	}

	// Init pool: A buffered channel is a thread-safe FIFO queue of pre-warmed
	// writers. Flushing incurs operation charges, so size the pool carefully.
	pool := make(chan *storage.Writer, poolSize)
	// On return, close pooled writers; unused objects remain unfinalized.
	defer func() {
		close(pool)
		for w := range pool {
			if err := w.Close(); err != nil {
				fmt.Fprintf(out, "Writer.Close: %v\n", err)
			}
		}
	}()

	for range poolSize {
		w, err := newPrewarmedWriter()
		if err != nil {
			return err
		}
		pool <- w
	}

	// Write: Take a pre-warmed writer (waiting for a refill if the pool is
	// empty) and commit with the faster Flush() instead of blocking on Close().
	var w *storage.Writer
	select {
	case w = <-pool:
	case <-time.After(10 * time.Second):
		return errors.New("timed out waiting for a pre-warmed writer")
	}
	if _, err := w.Write([]byte("0123456789")); err != nil {
		return errors.Join(fmt.Errorf("Writer.Write: %w", err), w.Close())
	}
	if _, err := w.Flush(); err != nil {
		return errors.Join(fmt.Errorf("Writer.Flush: %w", err), w.Close())
	}
	fmt.Fprintf(out, "Wrote and flushed %q without finalizing\n", w.Name)

	// Pool maintenance (run asynchronously off the critical write path):
	// Close the used writer without finalizing and refill the pool.
	maintenance := make(chan struct{})
	go func(used *storage.Writer) {
		defer close(maintenance)
		if err := used.Close(); err != nil {
			fmt.Fprintf(out, "Writer.Close: %v\n", err)
		}
		next, err := newPrewarmedWriter()
		if err != nil {
			fmt.Fprintf(out, "failed to pre-warm replacement writer: %v\n", err)
			return
		}
		pool <- next
	}(w)
	// Wait for maintenance on return; runs before the pool cleanup above.
	defer func() { <-maintenance }()

	// Read: Unfinalized objects are readable after Flush().
	r, err := bucket.Object(fmt.Sprintf("%s_0", keyPrefix)).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object.NewReader: %w", err)
	}
	contents, err := io.ReadAll(r)
	if err := errors.Join(err, r.Close()); err != nil {
		return fmt.Errorf("reading %s_0: %w", keyPrefix, err)
	}

	fmt.Fprintf(out, "Read unfinalized object %s_0: %s\n", keyPrefix, contents)
	return nil
}

// [END storage_optimize_write_latency_pool]
