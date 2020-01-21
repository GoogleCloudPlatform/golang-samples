// Copyright 2019 Google LLC
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

// Sample readEncryptedObject demonstrates reading an encrypted object using Cloud KMS key.
package objects

// [START storage_download_encrypted_file]
import (
	"context"
	"io/ioutil"

	"cloud.google.com/go/storage"
)

// readEncryptedObject reads an encrypted object.
func readEncryptedObject(bucket, object string, secretKey []byte) ([]byte, error) {
	// bucket := "bucket-name"
	// object := "object-name"
	// key := []byte("secret-encryption-key")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	obj := client.Bucket(bucket).Object(object)
	rc, err := obj.Key(secretKey).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// [END storage_download_encrypted_file]
