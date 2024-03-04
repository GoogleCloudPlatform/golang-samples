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

package hmac

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
)

var serviceAccountEmail = os.Getenv("GOLANG_SAMPLES_SERVICE_ACCOUNT_EMAIL")
var storageClient *storage.Client

func TestMain(m *testing.M) {
	if serviceAccountEmail == "" {
		fmt.Fprintln(os.Stderr, "GOLANG_SAMPLES_SERVICE_ACCOUNT_EMAIL not set. Skipping.")
		return
	}
	ctx := context.Background()
	storageClient, _ = storage.NewClient(ctx)
	defer storageClient.Close()

	// Delete all existing HMAC keys in the project to avoid running into
	// resource constraints during the test.
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	if err := deleteAllKeys(projectID); err != nil {
		fmt.Printf("deleting existing keys: %v", err)
	}

	os.Exit(m.Run())
}

func TestListKeys(t *testing.T) {
	tc := testutil.SystemTest(t)
	key, err := createTestKey(tc.ProjectID, t)
	defer deleteTestKey(key)
	if err != nil {
		t.Fatalf("Error in key creation: %s", err)
	}

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		keys, err := listHMACKeys(ioutil.Discard, tc.ProjectID)
		if err != nil {
			r.Errorf("listHMACKeys raised error: %s", err)
		}
		if len(keys) < 1 {
			r.Errorf("Should have at least one key listed.")
		}
	})
}

func TestCreateKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	key, err := createHMACKey(ioutil.Discard, tc.ProjectID, serviceAccountEmail)
	defer deleteTestKey(key)

	if err != nil {
		t.Fatalf("createHMACKey raised error: %s", err)
	}
	if key == nil {
		t.Fatalf("Returned nil key.")
	}
	if key.State != "ACTIVE" {
		t.Errorf("State of key is %s, should be ACTIVE", key.State)
	}
}

func TestActivateKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	key, err := createTestKey(tc.ProjectID, t)
	defer deleteTestKey(key)
	if err != nil {
		t.Fatalf("Error in key creation: %s", err)
	}

	// Key must first be deactivated in order to update to active state.
	ctx := context.Background()
	handle := storageClient.HMACKeyHandle(key.ProjectID, key.AccessID)
	handle.Update(ctx, storage.HMACKeyAttrsToUpdate{State: "INACTIVE"})

	key, err = activateHMACKey(ioutil.Discard, key.AccessID, key.ProjectID)

	if err != nil {
		t.Fatalf("Error in activateHMACKey: %s", err)
	}
	if key.State != "ACTIVE" {
		t.Errorf("State of key is %s, should be ACTIVE", key.State)
	}
}

func TestDeactivateKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	key, err := createTestKey(tc.ProjectID, t)
	defer deleteTestKey(key)
	if err != nil {
		t.Fatalf("Error in key creation: %s", err)
	}

	key, err = deactivateHMACKey(ioutil.Discard, key.AccessID, key.ProjectID)
	if err != nil {
		t.Fatalf("Error in deactivateHMACKey: %s", err)
	}
	if key.State != "INACTIVE" {
		t.Errorf("State of key is %s, should be INACTIVE", key.State)
	}
}

func TestGetKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	key, err := createTestKey(tc.ProjectID, t)
	defer deleteTestKey(key)
	if err != nil {
		t.Fatalf("Error in key creation: %s", err)
	}

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		gotKey, err := getHMACKey(ioutil.Discard, key.AccessID, key.ProjectID)
		if err != nil {
			r.Errorf("Error in getHMACKey: %s", err)
			return
		}
		if gotKey == nil {
			r.Errorf("Returned nil key.")
		}
	})
}

func TestDeleteKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	key, err := createTestKey(tc.ProjectID, t)
	defer deleteTestKey(key)
	if err != nil {
		t.Fatalf("Error in key creation: %s", err)
	}

	// Keys must be in INACTIVE state before deletion.
	ctx := context.Background()
	handle := storageClient.HMACKeyHandle(key.ProjectID, key.AccessID)
	_, err = handle.Update(ctx, storage.HMACKeyAttrsToUpdate{State: "INACTIVE"})
	if err != nil {
		t.Errorf("Error updating HMAC key: %s", err)
	}

	// Retry as HMAC key updates can take up to 3 minutes to propagate.
	testutil.Retry(t, 18, 10*time.Second, func(r *testutil.R) {
		err = deleteHMACKey(io.Discard, key.AccessID, key.ProjectID)
		if err != nil {
			// 400 error with reason "invalid" means the key was already deleted;
			// this can happen if there was a previous call that returned an error
			// but succeeded on the server side.
			var ge *googleapi.Error
			if errors.As(err, &ge) {
				if ge.Code == 400 && len(ge.Errors) > 0 && ge.Errors[0].Reason == "invalid" {
					return
				}
			}
			r.Errorf("Error in deleteHMACKey: %s", err)
		}
		key, _ = handle.Get(ctx)
		if key != nil && key.State != "DELETED" {
			r.Errorf("State of key is %s, should be DELETED", key.State)
		}

	})
}

// Delete all HMAC keys in the project.
func deleteAllKeys(projectID string) error {
	iter := storageClient.ListHMACKeys(context.Background(), projectID)
	for {
		key, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("ListHMACKeys: %w", err)
		}
		deleteTestKey(key)
	}
	return nil
}

// Create a key for testing purposes.
func createTestKey(projectID string, t *testing.T) (*storage.HMACKey, error) {
	ctx := context.Background()
	var key *storage.HMACKey
	var err error

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		key, err = storageClient.CreateHMACKey(ctx, projectID, serviceAccountEmail)
		if err != nil {
			r.Errorf("Error in CreateHMACKey: %v", err)
			return
		}

		// Nil key check should not happen but is added to handle flaky
		// "nil pointer dereference" error.
		if key == nil {
			r.Errorf("CreateHMACKey returned nil key.")
			err = errors.New("CreateHMACKey returned nil key")
			return
		}
	})

	return key, err
}

// Deactivate and delete the given key. Should operate as a teardown method.
func deleteTestKey(key *storage.HMACKey) {
	if key == nil {
		return
	}
	ctx := context.Background()
	handle := storageClient.HMACKeyHandle(key.ProjectID, key.AccessID)
	if key.State == "ACTIVE" {
		handle.Update(ctx, storage.HMACKeyAttrsToUpdate{State: "INACTIVE"})
	}
	if key.State != "DELETED" {
		handle.Delete(ctx)
	}
}
