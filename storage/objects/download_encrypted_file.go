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

// [START storage_download_encrypted_file]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// downloadEncryptedFile reads an encrypted object.
func downloadEncryptedFile(w io.Writer, bucket, object string, secretKey []byte) ([]byte, error) {
	// bucket := "bucket-name"
	// object := "object-name"
	// key := []byte("secret-encryption-key")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	obj := client.Bucket(bucket).Object(object)

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := obj.Key(secretKey).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).Key(%q).NewReader: %w", object, secretKey, err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	fmt.Fprintf(w, "File %v downloaded with encryption key.\n", object)
	return data, nil
}

// [END storage_download_encrypted_file]
