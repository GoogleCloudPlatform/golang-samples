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
package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"
	"bytes"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	s := m.Run()
	log.SetOutput(os.Stderr)
	os.Exit(s)
}

func TestList(t *testing.T) {
	bucketName := os.Getenv("GOOGLE_CLOUD_PROJECT_S3_SDK")
	googleAccessKeyID := os.Getenv("STORAGE_HMAC_ACCESS_KEY_ID")
	googleAccessKeySecret := os.Getenv("STORAGE_HMAC_ACCESS_SECRET_KEY")

	buf := new(bytes.Buffer)
	buckets, err := listGCSBuckets(buf, googleAccessKeyID, googleAccessKeySecret)
	if err != nil {
		t.Errorf("Unable to list GCS buckets: %v", err)
	}

	got := buf.String()
	if want := bucketName; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) { // for eventual consistency
		for _, b := range buckets {
			if aws.StringValue(b.Name) == bucketName {
				return
			}
		}
		time.Sleep(2 * time.Second)
		r.Errorf("got bucket list: %v; want %q in the list", buckets, bucketName)
	})
}
