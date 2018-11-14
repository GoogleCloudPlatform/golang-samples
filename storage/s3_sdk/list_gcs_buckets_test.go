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

	"cloud.google.com/go/storage"
)

var (
	storageClient *storage.Client
	bucketName    string
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	s := m.Run()
	log.SetOutput(os.Stderr)
	os.Exit(s)
}

func setup(t *testing.T) {
	tc := testutil.SystemTest(t)

	bucketName = tc.ProjectID + "-storage-buckets-tests"
}

func TestList(t *testing.T) {
	setup(t)

	buckets, err := list_gcs_buckets()
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
