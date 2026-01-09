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

// [START storage_open_object_read_full_object]
import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"cloud.google.com/go/storage/experimental"
)

// OpenObjectReadFullObject reads a full object's data from a
// rapid bucket.
func openObjectReadFullObject(w io.Writer, bucket, object string) ([]byte, error) {
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

	// Read the first KiB of the file and copy into a buffer.
	r, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewReader: %w", err)
	}
	defer r.Close()
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r); err != nil {
		return nil, fmt.Errorf("copying data: %v", err)
	}

	fmt.Fprintf(w, "Read the data of %v into a buffer\n", object)

	return buf.Bytes(), nil
}

// [END storage_open_object_read_full_object]
