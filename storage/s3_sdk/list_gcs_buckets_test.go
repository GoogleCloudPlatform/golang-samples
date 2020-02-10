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

// Sample lists GCS buckets using the S3 SDK using interoperability mode.
package s3sdk

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	s := m.Run()
	log.SetOutput(os.Stderr)
	os.Exit(s)
}

func TestList(t *testing.T) {
	t.Skip("HMAC secrets are not set up correctly: https://github.com/GoogleCloudPlatform/golang-samples/issues/1213")
	googleAccessKeyID := os.Getenv("STORAGE_HMAC_ACCESS_KEY_ID")
	googleAccessKeySecret := os.Getenv("STORAGE_HMAC_ACCESS_SECRET_KEY")

	if googleAccessKeyID == "" || googleAccessKeySecret == "" {
		t.Skip("STORAGE_HMAC_ACCESS_KEY_ID and STORAGE_HMAC_ACCESS_SECRET_KEY must be set. Skipping.")
	}

	buf := new(bytes.Buffer)
	_, err := listGCSBuckets(buf, googleAccessKeyID, googleAccessKeySecret)
	if err != nil {
		t.Errorf("listGCSBuckets: %v", err)
	}

	got := buf.String()
	if want := "Buckets:"; !strings.Contains(got, want) {
		t.Errorf("listGCSBuckets got\n----\n%s\n----\nWant to contain\n----\n%s\n----", got, want)
	}
}
