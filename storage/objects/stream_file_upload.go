// Copyright 2021 Google LLC
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

// [START storage_stream_file_upload]
// [START storage_file_upload_from_memory]
import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// streamFileUpload uploads an object via a stream.
func streamFileUpload(w io.Writer, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	b := []byte("Hello world.")
	buf := bytes.NewBuffer(b)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// Set a generation-match precondition. The request to upload is aborted
	// if the object's generation number does not match your precondition
	// criteria. This avoids race conditions and data corruption.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	o = o.If(storage.Conditions{DoesNotExist: true})

	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	// 	return fmt.Errorf("object.Attrs: %v", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	wc.ChunkSize = 0 // note retries are not supported for chunk size 0.

	if _, err = io.Copy(wc, buf); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	// Data can continue to be added to the file until the writer is closed.
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Fprintf(w, "%v uploaded to %v.\n", object, bucket)

	return nil
}

// [END storage_stream_file_upload]
// [END storage_file_upload_from_memory]
