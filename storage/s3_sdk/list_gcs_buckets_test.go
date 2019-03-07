// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/aws/aws-sdk-go/aws"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	s := m.Run()
	log.SetOutput(os.Stderr)
	os.Exit(s)
}

func TestList(t *testing.T) {
	tc := testutil.SystemTest(t)

	bucketName := os.Getenv("CLOUD_CLOUD_PROJECT_S3_SDK")
	googleAccessKeyID := os.Getenv("STORAGE_HMAC_ACCESS_KEY_ID")
	googleAccessKeySecret := os.Getenv("STORAGE_HMAC_ACCESS_SECRET_KEY")

	buckets, err := list_gcs_buckets(googleAccessKeyID, googleAccessKeySecret)
	if err != nil {
		t.Errorf("Unable to list GCS buckets: %v", err)
	}

	var ok bool
outer:
	for attempt := 0; attempt < 5; attempt++ { // for eventual consistency
		for _, b := range buckets {
			if aws.StringValue(b.Name) == bucketName {
				ok = true
				break outer
			}
		}
		time.Sleep(2 * time.Second)
	}
	if !ok {
		t.Errorf("got bucket list: %v; want %q in the list", buckets, bucketName)
	}
}
