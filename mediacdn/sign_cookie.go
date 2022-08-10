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

// [START mediacdn_sign_cookie]
import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"io"
	"time"
)

// signCookie prints the signed cookie value for the specified URL prefix and configuration.
func signCookie(w io.Writer, urlPrefix, keyName string, privateKey []byte, expires time.Time) error {
	// urlPrefix := "http://example.com"
	// keyName := "your_key_name"
	// privateKey := "[]byte{34, 31, ...}"
	// expires := time.Unix(1558131350, 0)

	toSign := fmt.Sprintf(
		"URLPrefix=%s:Expires=%d:KeyName=%s",
		base64.RawURLEncoding.EncodeToString([]byte(urlPrefix)),
		expires.Unix(),
		keyName,
	)
	sig := ed25519.Sign(privateKey, []byte(toSign))

	fmt.Fprintf(
		w,
		"Edge-Cache-Cookie=%s:Signature=%s",
		toSign,
		base64.RawURLEncoding.EncodeToString(sig),
	)

	return nil
}

// [END mediacdn_sign_cookie]
