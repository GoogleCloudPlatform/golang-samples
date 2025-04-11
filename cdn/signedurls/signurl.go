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

// Package signedurl creates a signed URL for a Cloud CDN endpoint with the
// given key.
package signedurl

// [START cloudcdn_sign_url]
// [START cloudcdn_sign_url_prefix]
import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// [END cloudcdn_sign_url_prefix]
// SignURL creates a signed URL for an endpoint on Cloud CDN.
//
// - url must start with "https://" and should not have the "Expires", "KeyName", or "Signature"
// query parameters.
// - key should be in raw form (not base64url-encoded) which is 16-bytes long.
// - keyName must match a key added to the backend service or bucket.
func signURL(url, keyName string, key []byte, expiration time.Time) string {
	sep := "?"
	if strings.Contains(url, "?") {
		sep = "&"
	}
	url += sep
	url += fmt.Sprintf("Expires=%d", expiration.Unix())
	url += fmt.Sprintf("&KeyName=%s", keyName)

	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(url))
	sig := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	url += fmt.Sprintf("&Signature=%s", sig)
	return url
}

// [END cloudcdn_sign_url]
// [START cloudcdn_sign_url_prefix]
// SignURLWithPrefix creates a signed URL prefix for an endpoint on Cloud CDN.
// Prefixes allow access to any URL with the same prefix, and can be useful for
// granting access broader content without signing multiple URLs.
//
// - urlPrefix must start with "https://" and should not include query parameters.
// - key should be in raw form (not base64url-encoded) which is 16-bytes long.
// - keyName must match a key added to the backend service or bucket.
func signURLWithPrefix(urlPrefix, keyName string, key []byte, expiration time.Time) (string, error) {
	if strings.Contains(urlPrefix, "?") {
		return "", fmt.Errorf("urlPrefix must not include query params: %s", urlPrefix)
	}

	encodedURLPrefix := base64.URLEncoding.EncodeToString([]byte(urlPrefix))
	input := fmt.Sprintf("URLPrefix=%s&Expires=%d&KeyName=%s",
		encodedURLPrefix, expiration.Unix(), keyName)

	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(input))
	sig := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	signedValue := fmt.Sprintf("%s&Signature=%s", input, sig)

	return signedValue, nil
}

// [END cloudcdn_sign_url_prefix]
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

func generateSignedURLs(w io.Writer) error {
	// The path to a file containing the base64-encoded signing key
	keyPath := os.Getenv("KEY_PATH")

	// Note: consider using the GCP Secret Manager for managing access to your
	// signing key(s).
	key, err := readKeyFile(keyPath)
	if err != nil {
		return err
	}

	// Generate a signed URL for a specific URL
	url := signURL("https://example.com/media/1234.m3e8", "my-key", key, time.Now().Add(time.Hour*24))
	fmt.Fprintln(w, url)

	// Generate a signed URL for
	urlPrefix, err := signURLWithPrefix("https://www.google.com/", "my-key", key, time.Unix(1549751401, 0))
	if err != nil {
		return err
	}

	fmt.Fprintln(w, urlPrefix)

	return nil
}
