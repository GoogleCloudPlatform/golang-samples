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

// Package signedcookie creates a signed URL for a Cloud CDN endpoint with the
// given key.
package signedcookie

// [START cloudcdn_sign_cookie]
import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// signCookie creates a signed cookie for an endpoint served by Cloud CDN.
//
// - urlPrefix must start with "https://" and should include the path prefix
// for which the cookie will authorize access to.
// - key should be in raw form (not base64url-encoded) which is
// 16-bytes long.
// - keyName must match a key added to the backend service or bucket.
func signCookie(urlPrefix, keyName string, key []byte, expiration time.Time) (string, error) {
	encodedURLPrefix := base64.URLEncoding.EncodeToString([]byte(urlPrefix))
	input := fmt.Sprintf("URLPrefix=%s:Expires=%d:KeyName=%s",
		encodedURLPrefix, expiration.Unix(), keyName)

	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(input))
	sig := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	signedValue := fmt.Sprintf("%s:Signature=%s",
		input,
		sig,
	)

	return signedValue, nil
}

// [END cloudcdn_sign_cookie]

// readKeyFile reads the base64url-encoded key file and decodes it.
func readKeyFile(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %+v", err)
	}
	d := make([]byte, base64.URLEncoding.DecodedLen(len(b)))
	n, err := base64.URLEncoding.Decode(d, b)
	if err != nil {
		return nil, fmt.Errorf("failed to base64url decode: %+v", err)
	}
	return d[:n], nil
}

func generateSignedCookie(w io.Writer) error {
	// The path to a file containing the base64-encoded signing key
	keyPath := os.Getenv("KEY_PATH")

	// Note: consider using the GCP Secret Manager for managing access to your
	// signing key(s).
	key, err := readKeyFile(keyPath)
	if err != nil {
		return err
	}

	var (
		// domain and path should match the user-facing URL for accessing
		// content.
		domain     = "media.example.com"
		path       = "/segments/"
		keyName    = "my-key"
		expiration = time.Hour * 2
	)

	signedValue, err := signCookie(fmt.Sprintf("https://%s%s", domain,
		path), keyName, key, time.Now().Add(expiration))
	if err != nil {
		return err
	}

	// Use Go's http.Cookie type to construct a cookie.
	cookie := &http.Cookie{
		Name:   "Cloud-CDN-Cookie",
		Value:  signedValue,
		Path:   path, // Best practice: only send the cookie for paths it is valid for
		Domain: domain,
		MaxAge: int(expiration.Seconds()),
	}

	// We print this to stdout in this example. In a real application, use the
	// SetCookie method on a http.ResponseWriter to write the cookie to the
	// user.
	fmt.Fprintln(w, cookie)

	return nil
}
