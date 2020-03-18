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

// Command signedurls creates a signed URL for a Cloud CDN endpoint with the
// given key.
package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// [START example]

// SignCookie creates a signed cookie for an endpoint served by Cloud CDN.
//
// - urlPrefix must start with "https://" and should include the path prefix
// for which the cookie will authorize access to.
// - key should be in raw form (not base64url-encoded) which is
// 16-bytes long.
// - keyName must match a key added to the backend service or bucket.
func SignCookie(urlPrefix, keyName string, key []byte, expiration time.Time) (string, error) {
	if !strings.HasPrefix(urlPrefix, "http") {
		return "", fmt.Errorf("the provided urlPrefix is missing the protocol: %s", urlPrefix)
	}

	encodedURLPrefix := base64.URLEncoding.EncodeToString([]byte(urlPrefix))
	input := fmt.Sprintf("URLPrefix=%s:Expires=%d:Keyname=%s",
		encodedURLPrefix,
		expiration.Unix(),
		keyName,
	)

	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(input))
	sig := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	signedValue := fmt.Sprintf("%s:Signature=%s",
		input,
		sig,
	)

	return signedValue, nil
}

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

func main() {
	var keyPath string
	flag.StringVar(&keyPath, "key-file", "", "The path to a file containing the base64-encoded signing key")
	flag.Parse()

	key, err := readKeyFile(keyPath)
	if err != nil {
		log.Fatal(err)
	}

	var (
		// domain and path should match the user-facing URL for accessing
		// content.
		domain     = "media.example.com"
		path       = "/segments/"
		keyName    = "my-key"
		expiration = time.Hour * 2
	)
	signedValue, err := SignCookie(
		fmt.Sprintf("https://%s%s", domain, path),
		keyName,
		key,
		//time.Now().Add(expiration),
		time.Unix(1558131350, 0),
	)

	// Use Go's http.Cookie type to construct a cookie.
	cookie := &http.Cookie{
		Name:   "Cloud-CDN-Cookie",
		Value:  signedValue,
		Path:   "/segments/", // Best practice: only send the cookie for paths it is valid for
		Domain: "media.example.com",
		MaxAge: int(expiration.Seconds()),
	}

	// We print this to stdout in this example.
	// Use http.ResponseWriter.SetCookie to write a cookie to an authenticated
	// client in a real application.
	fmt.Println(cookie)
}

// [END example]
