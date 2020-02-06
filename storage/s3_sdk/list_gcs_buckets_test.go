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
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var serviceAccountEmail = os.Getenv("GOLANG_SAMPLES_SERVICE_ACCOUNT_EMAIL")
var storageClient *storage.Client

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	s := m.Run()
	log.SetOutput(os.Stderr)

	ctx := context.Background()
	storageClient, _ = storage.NewClient(ctx)
	defer storageClient.Close()

	tc, _ := testutil.ContextMain(m)
	key, err := createTestKey(tc.ProjectID)
	if err != nil {
		log.Fatalf("error setting up HMAC key: %v", err)
	}
	defer deleteTestKey(key)

	os.Exit(s)
}

// Create a key for testing purposes and set environment variables
func createTestKey(projectID string) (*storage.HMACKey, error) {
	ctx := context.Background()
	key, err := storageClient.CreateHMACKey(ctx, projectID, serviceAccountEmail)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	if err = os.Setenv("STORAGE_HMAC_ACCESS_KEY_ID", key.AccessID); err != nil {
		return key, fmt.Errorf("error setting STORAGE_HMAC_ACCESS_KEY_ID: %v", err)
	}
	if err = os.Setenv(key.Secret, "STORAGE_HMAC_ACCESS_SECRET_KEY"); err != nil {
		return key, fmt.Errorf("error setting STORAGE_HMAC_ACCESS_SECRET_KEY: %v", err)
	}
	return key, nil
}

// Deactivate and delete the given key. Should operate as a teardown method.
func deleteTestKey(key *storage.HMACKey) {
	ctx := context.Background()
	handle := storageClient.HMACKeyHandle(key.ProjectID, key.AccessID)
	if key.State == "ACTIVE" {
		handle.Update(ctx, storage.HMACKeyAttrsToUpdate{State: "INACTIVE"})
	}
	if key.State != "DELETED" {
		handle.Delete(ctx)
	}
}

func TestList(t *testing.T) {
	googleAccessKeyID := os.Getenv("STORAGE_HMAC_ACCESS_KEY_ID")
	googleAccessKeySecret := os.Getenv("STORAGE_HMAC_ACCESS_SECRET_KEY")

	if googleAccessKeyID == ""  {
		t.Errorf("STORAGE_HMAC_ACCESS_KEY_ID must be set.")
	}
	if googleAccessKeySecret == "" {
		t.Errorf("STORAGE_HMAC_ACCESS_SECRET_KEY must be set.")
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
