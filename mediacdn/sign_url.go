// Copyright 2022 Google LLC
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

package snippets

// [START mediacdn_sign_url]
import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"
)

// signURL prints the signed URL string for the specified URL and configuration.
func signURL(w io.Writer, url, keyName string, privateKey []byte, expires time.Time) error {
	// url := "http://example.com"
	// keyName := "your_key_name"
	// privateKey := "[]byte{34, 31, ...}"
	// expires := time.Unix(1558131350, 0)

	sep := '?'
	if strings.ContainsRune(url, '?') {
		sep = '&'
	}
	toSign := fmt.Sprintf("%s%cExpires=%d&KeyName=%s", url, sep, expires.Unix(), keyName)
	sig := ed25519.Sign(privateKey, []byte(toSign))

	fmt.Fprintf(w, "%s&Signature=%s", toSign, base64.RawURLEncoding.EncodeToString(sig))

	return nil
}

// [END mediacdn_sign_url]
