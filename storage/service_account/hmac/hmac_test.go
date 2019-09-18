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
	"fmt"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"io/ioutil"
	"os"
	"testing"

	"cloud.google.com/go/storage"
)

var serviceAccountEmail = os.Getenv("GOLANG_SAMPLES_SERVICE_ACCOUNT_EMAIL")
var storageClient *storage.Client

func TestMain(m *testing.M) {
	ctx := context.Background()
	storageClient, _ = storage.NewClient(ctx)
	defer storageClient.Close()

	os.Exit(m.Run())
}

func TestListKeys(t *testing.T) {
	tc := testutil.SystemTest(t)
	key, err := createTestKey(tc.ProjectID)
	defer deleteTestKey(key)
	if err != nil {
		t.Errorf("Error in key creation: %s", err)
	}

	keys, err := listHMACKeys(ioutil.Discard, tc.ProjectID)
	if err != nil {
		t.Errorf("listHMACKeys raised error: %s", err)
	}
	if len(keys) < 1 {
		t.Errorf("Should have at least one key listed.")
	}
}

func TestCreateKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	key, err := createHMACKey(ioutil.Discard, tc.ProjectID, serviceAccountEmail)
	defer deleteTestKey(key)

	if err != nil {
		t.Errorf("createHMACKey raised error: %s", err)
	}
	if key == nil {
		t.Errorf("Returned nil key.")
	}
	if key.State != "ACTIVE" {
		t.Errorf("State of key is %s, should be ACTIVE", key.State)
	}
}

func TestActivateKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	key, err := createTestKey(tc.ProjectID)
	defer deleteTestKey(key)
	if err != nil {
		t.Errorf("Error in key creation: %s", err)
	}

	// Key must first be deactivated in order to update to active state.
	ctx := context.Background()
	handle := storageClient.HMACKeyHandle(key.ProjectID, key.AccessID)
	handle.Update(ctx, storage.HMACKeyAttrsToUpdate{State: "INACTIVE"})

	key, err = activateHMACKey(ioutil.Discard, key.AccessID, key.ProjectID)

	if err != nil {
		t.Errorf("Error in activateHMACKey: %s", err)
	}
	if key.State != "ACTIVE" {
		t.Errorf("State of key is %s, should be ACTIVE", key.State)
	}
}

func TestDeactivateKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	key, err := createTestKey(tc.ProjectID)
	defer deleteTestKey(key)
	if err != nil {
		t.Errorf("Error in key creation: %s", err)
	}

	key, err = deactivateHMACKey(ioutil.Discard, key.AccessID, key.ProjectID)
	if err != nil {
		t.Errorf("Error in deactivateHMACKey: %s", err)
	}
	if key.State != "INACTIVE" {
		t.Errorf("State of key is %s, should be INACTIVE", key.State)
	}
}

func TestGetKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	key, err := createTestKey(tc.ProjectID)
	defer deleteTestKey(key)
	if err != nil {
		t.Errorf("Error in key creation: %s", err)
	}

	key, err = getHMACKey(ioutil.Discard, key.AccessID, key.ProjectID)
	if err != nil {
		t.Errorf("Error in getHMACKey: %s", err)
	}
	if key == nil {
		t.Errorf("Returned nil key.")
	}
}

func TestDeleteKey(t *testing.T) {
	tc := testutil.SystemTest(t)
	key, err := createTestKey(tc.ProjectID)
	defer deleteTestKey(key)
	if err != nil {
		t.Errorf("Error in key creation: %s", err)
	}

	// Keys must be in INACTIVE state before deletion.
	ctx := context.Background()
	handle := storageClient.HMACKeyHandle(key.ProjectID, key.AccessID)
	handle.Update(ctx, storage.HMACKeyAttrsToUpdate{State: "INACTIVE"})

	err = deleteHMACKey(ioutil.Discard, key.AccessID, key.ProjectID)
	if err != nil {
		t.Errorf("Error in deleteHMACKey: %s", err)
	}
	key, _ = handle.Get(ctx)
	if key.State != "DELETED" {
		t.Errorf("State of key is %s, should be DELETED", key.State)
	}

}

// Create a key for testing purposes.
func createTestKey(projectID string) (*storage.HMACKey, error) {
	ctx := context.Background()
	key, err := storageClient.CreateHMACKey(ctx, projectID, serviceAccountEmail)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
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
