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

// [START storage_rotate_encryption_key]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// rotateEncryptionKey encrypts an object with the newKey.
func rotateEncryptionKey(w io.Writer, bucket, object string, key, newKey []byte) error {
	// bucket := "bucket-name"
	// object := "object-name"
	// key := []byte("encryption-key")
	// newKey := []byte("new-encryption-key")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	obj := client.Bucket(bucket).Object(object)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// obj is encrypted with key, we are encrypting it with the newKey.
	_, err = obj.Key(newKey).CopierFrom(obj.Key(key)).Run(ctx)
	if err != nil {
		return fmt.Errorf("Key(%q).CopierFrom(%q).Run: %v", newKey, key, err)
	}
	fmt.Fprintf(w, "Key rotation complete for blob %v.\n", object)
	return nil
}

// [END storage_rotate_encryption_key]
