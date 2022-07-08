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

package snippets

// [START mediacdn_sign_cookie]
import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"time"
)

// signCookie prints the Signed cookie value for the specified URL prefix and configuration.
func signCookie(urlPrefix, keyName string, base64Key []byte, expires time.Time) string {
	// urlPrefix := "http://example.com"
	// keyName := "your_key_name"
	// base64Key := "[]byte{34, 31, ...}"
	// expires := time.Unix(1558131350, 0)

	toSign := fmt.Sprintf(
		"URLPrefix=%s:Expires=%d:KeyName=%s",
		base64.RawURLEncoding.EncodeToString([]byte(urlPrefix)),
		expires.Unix(),
		keyName,
	)
	sig := ed25519.Sign(base64Key, []byte(toSign))

	return fmt.Sprintf(
		"Edge-Cache-Cookie=%s:Signature=%s",
		toSign,
		base64.RawURLEncoding.EncodeToString(sig),
	)
}

// [END mediacdn_sign_cookie]
