// Copyright 2020 Google LLC
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

// [START storage_upload_encrypted_file]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// uploadEncyptedFile writes an object using AES-256 encryption key.
func uploadEncyptedFile(w io.Writer, bucket, object string, secretKey []byte) error {
	// bucket := "bucket-name"
	// object := "object-name"
	// secretKey := []byte("secret-key")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	obj := client.Bucket(bucket).Object(object)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Encrypt the object's contents.
	wc := obj.Key(secretKey).NewWriter(ctx)
	if _, err := wc.Write([]byte("top secret")); err != nil {
		return fmt.Errorf("Writer.Write: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Fprintf(w, "Uploaded encrypted object %v.\n", object)
	return nil
}

// [END storage_upload_encrypted_file]
