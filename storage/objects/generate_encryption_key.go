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

// [START storage_generate_encryption_key]
import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// generateEncryptionKey generates a 256 bit (32 byte) AES encryption key and
// prints the base64 representation.
func generateEncryptionKey(w io.Writer) error {
	// This is included for demonstration purposes. You should generate your own
	// key. Please remember that encryption keys should be handled with a
	// comprehensive security policy.
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return fmt.Errorf("rand.Read: %w", err)
	}
	encryptionKey := base64.StdEncoding.EncodeToString(key)
	fmt.Fprintf(w, "Generated base64-encoded encryption key: %v\n", encryptionKey)
	return nil
}

// [END storage_generate_encryption_key]
