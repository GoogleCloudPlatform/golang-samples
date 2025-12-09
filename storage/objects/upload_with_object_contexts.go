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

package objects

// [START storage_upload_with_object_contexts]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// uploadWithObjectContexts sets an object's contexts.
func uploadWithObjectContexts(w io.Writer, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// To set contexts on a new object during write:
	writer := o.NewWriter(ctx)
	writer.Contexts = &storage.ObjectContexts{
		Custom: map[string]storage.ObjectCustomContextPayload{
			"key1": {Value: "value1"},
			"key2": {Value: "value2"},
		},
	}
	if _, err := writer.Write([]byte("test")); err != nil {
		return fmt.Errorf("Writer.Write: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	fmt.Fprintf(w, "Added new contexts to object %v in bucket %v.\n", object, bucket)

	return nil
}

// [END storage_upload_with_object_contexts]
