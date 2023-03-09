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
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestListGCSBuckets(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Set up service account HMAC key to use for this test.
	key, err := createTestKey(ctx, t, client, tc.ProjectID)
	if err != nil {
		t.Fatalf("error setting up HMAC key: %v", err)
	}
	defer deleteTestKey(ctx, client, key)

	buf := new(bytes.Buffer)
	// New HMAC key may take up to 15s to propagate, so we need to retry for up
	// to that amount of time.
	testutil.Retry(t, 75, time.Millisecond*200, func(r *testutil.R) {
		buf.Reset()
		if err := listGCSBuckets(buf, key.AccessID, key.Secret); err != nil {
			r.Errorf("listGCSBuckets: %v", err)
		}
	})

	got := buf.String()
	if want := "Buckets:"; !strings.Contains(got, want) {
		t.Fatalf("listGCSBuckets got\n----\n%s\n----\nWant to contain\n----\n%s\n----", got, want)
	}
}

// Create a key for testing purposes and set environment variables
func createTestKey(ctx context.Context, t *testing.T, client *storage.Client, projectID string) (*storage.HMACKey, error) {
	email := os.Getenv("GOLANG_SAMPLES_SERVICE_ACCOUNT_EMAIL")
	if email == "" {
		t.Skip("GOLANG_SAMPLES_SERVICE_ACCOUNT_EMAIL must be defined in the environment")
		return nil, nil
	}
	var key *storage.HMACKey
	var err error
	// TODO: replace testutil.Retry with retry config on client when available.
	testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		key, err = client.CreateHMACKey(ctx, projectID, email)
		if err != nil {
			r.Errorf("CreateHMACKey: %v", err)
		}
	})
	return key, err
}

// Deactivate and delete the given key. Should operate as a teardown method.
func deleteTestKey(ctx context.Context, client *storage.Client, key *storage.HMACKey) {
	handle := client.HMACKeyHandle(key.ProjectID, key.AccessID)
	if key.State == "ACTIVE" {
		handle.Update(ctx, storage.HMACKeyAttrsToUpdate{State: "INACTIVE"})
	}
	if key.State != "DELETED" {
		handle.Delete(ctx)
	}
}

func TestListGCSObjects(t *testing.T) {
	ctx := context.Background()
	tc := testutil.SystemTest(t)

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Set up service account HMAC key to use for this test.
	key, err := createTestKey(ctx, t, client, tc.ProjectID)
	if err != nil {
		t.Fatalf("error setting up HMAC key: %v", err)
	}
	defer deleteTestKey(ctx, client, key)

	buf := new(bytes.Buffer)
	testutil.Retry(t, 5, time.Second, func(r *testutil.R) {
		if err := listGCSObjects(buf, "cloud-samples-data", key.AccessID, key.Secret); err != nil {
			r.Errorf("listGCSObjects: %v", err)
		}

		got := buf.String()
		if want := "Objects:"; !strings.Contains(got, want) {
			r.Errorf("listGCSObjects got\n----\n%s\n----\nWant to contain\n----\n%s\n----", got, want)
		}
	})

}
