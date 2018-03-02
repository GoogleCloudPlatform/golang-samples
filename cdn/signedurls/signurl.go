// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command signedurls creates a signed URL for a Cloud CDN endpoint with the
// given key.
package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

// [START example]

// SignURL creates a signed URL for an endpoint on Cloud CDN. url must start
// with "https://" and should not have the "Expires", "KeyName", or "Signature"
// query parameters. key should be in raw form (not base64url-encoded) which is
// 16-bytes long. keyName should be added to the backend service or bucket.
func SignURL(url, keyName string, key []byte, expiration time.Time) string {
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
	key, err := readKeyFile("/path/to/key")
	if err != nil {
		log.Fatal(err)
	}
	url := SignURL("https://example.com", "MY-KEY", key, time.Now().Add(time.Hour*24))
	fmt.Println(url)
}

// [END example]
