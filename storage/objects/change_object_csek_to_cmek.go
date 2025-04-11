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

// [START storage_object_csek_to_cmek]
import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// сhangeObjectCSEKtoKMS changes the key used to encrypt an object from
// a customer-supplied encryption key to a customer-managed encryption key.
func сhangeObjectCSEKToKMS(w io.Writer, bucket, object string, encryptionKey []byte, kmsKeyName string) error {
	// bucket := "bucket-name"
	// object := "object-name"

	// encryptionKey is the Base64 encoded decryption key, which should be the same
	// key originally used to encrypt the object.
	// encryptionKey := []byte("TIbv/fjexq+VmtXzAlc63J4z5kFmWJ6NdAPQulQBT7g=")

	// kmsKeyName is the name of the KMS key to manage this object with.
	// kmsKeyName := "projects/projectId/locations/global/keyRings/keyRingID/cryptoKeys/cryptoKeyID"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to copy is aborted if the
	// object's generation number does not match your precondition.
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %w", err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// You can't change an object's encryption key directly. Instead, you must
	// rewrite the object using the new key.
	src := o.Key(encryptionKey)
	c := o.CopierFrom(src)
	c.DestinationKMSKeyName = kmsKeyName
	if _, err := c.Run(ctx); err != nil {
		return fmt.Errorf("Copier.Run: %w", err)
	}
	fmt.Fprintf(w, "Object %v in bucket %v is now managed by the KMS key %v instead of a customer-supplied encryption key\n", object, bucket, kmsKeyName)
	return nil
}

// [END storage_object_csek_to_cmek]
