// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

var (
	c          *storage.Client
	bucketName string
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	c, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	tc, ok := testutil.ContextMain(m)
	if !ok {
		log.Println("Skipping tests: GOLANG_SAMPLES_PROJECT_ID must be set")
		return
	}

	bucketName = tc.ProjectID + "-storage-tests"
	// Clean up buckets before running tests.
	delete(c, bucketName)
	delete(c, bucketName+"-attrs")
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	tc := testutil.SystemTest(t)
	if err := create(c, tc.ProjectID, bucketName); err != nil {
		t.Fatalf("failed to create bucket (%q): %v", bucketName, err)
	}
}

func TestCreateWithAttrs(t *testing.T) {
	tc := testutil.SystemTest(t)
	name := bucketName + "-attrs"
	if err := createWithAttrs(c, tc.ProjectID, name); err != nil {
		t.Fatalf("failed to create bucket (%q): %v", name, err)
	}
	if err := delete(c, name); err != nil {
		t.Fatalf("failed to delete bucket (%q): %v", name, err)
	}
}

func TestList(t *testing.T) {
	tc := testutil.SystemTest(t)
	buckets, err := list(c, tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
	var ok bool
outer:
	for attempt := 0; attempt < 5; attempt++ { // for eventual consistency
		for _, b := range buckets {
			if b == bucketName {
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

func TestIAM(t *testing.T) {
	if _, err := getPolicy(c, bucketName); err != nil {
		t.Errorf("getPolicy: %#v", err)
	}
	if err := addUser(c, bucketName); err != nil {
		t.Errorf("addUser: %v", err)
	}
	if err := removeUser(c, bucketName); err != nil {
		t.Errorf("removeUser: %v", err)
	}
}

func TestRequesterPays(t *testing.T) {
	if err := enableRequesterPays(c, bucketName); err != nil {
		t.Errorf("enableRequesterPay: %#v", err)
	}
	if err := disableRequesterPays(c, bucketName); err != nil {
		t.Errorf("enableRequesterPay: %#v", err)
	}
	if err := checkRequesterPays(c, bucketName); err != nil {
		t.Errorf("enableRequesterPay: %#v", err)
	}
}

func TestDelete(t *testing.T) {
	if err := delete(c, bucketName); err != nil {
		t.Fatalf("failed to delete bucket (%q): %v", bucketName, err)
	}
}
