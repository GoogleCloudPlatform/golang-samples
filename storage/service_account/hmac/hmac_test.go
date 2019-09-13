// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//:
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hmac

import (
	"bytes"
	"context"
	"os"
	"testing"

	"cloud.google.com/go/storage"
)

var ProjectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
var ServiceAccountEmail = os.Getenv("HMAC_KEY_TEST_SERVICE_ACCOUNT")
var StorageClient *storage.Client

func TestMain(m *testing.M) {
	ctx := context.Background()
	StorageClient, _ = storage.NewClient(ctx)
	defer StorageClient.Close()

	s := m.Run()
	os.Exit(s)
}

func TestListKeys(t *testing.T) {
	key, err := CreateTestKey()
	defer DeleteTestKey(key)
	if err != nil {
		t.Errorf("Error in key creation: %s", err)
	}

	buf := new(bytes.Buffer)
	keys, err := listHMACKeys(buf, ProjectID)
	if err != nil {
		t.Errorf("listHMACKeys raised error: %s", err)
	}
	if len(keys) < 1 {
		t.Errorf("Should have at least one key listed.")
	}
}

func TestCreateKey(t *testing.T) {
	buf := new(bytes.Buffer)
	key, err := createHMACKey(buf, ProjectID, ServiceAccountEmail)
	defer DeleteTestKey(key)

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
	key, err := CreateTestKey()
	defer DeleteTestKey(key)
	if err != nil {
		t.Errorf("Error in key creation: %s", err)
	}

	// Key must first be deactivated in order to update to active state.
	ctx := context.Background()
	keyHandle := StorageClient.HMACKeyHandle(key.ProjectID, key.AccessID)
	keyHandle.Update(ctx, storage.HMACKeyAttrsToUpdate{State: "INACTIVE"})

	buf := new(bytes.Buffer)
	key, err = activateHMACKey(buf, key.AccessID, key.ProjectID)

	if err != nil {
		t.Errorf("Error in activateHMACKey: %s", err)
	}
	if key.State != "ACTIVE" {
		t.Errorf("State of key is %s, should be ACTIVE", key.State)
	}
}

func TestDeactivateKey(t *testing.T) {
	key, err := CreateTestKey()
	defer DeleteTestKey(key)
	if err != nil {
		t.Errorf("Error in key creation: %s", err)
	}

	buf := new(bytes.Buffer)
	key, err = deactivateHMACKey(buf, key.AccessID, key.ProjectID)
	if err != nil {
		t.Errorf("Error in deactivateHMACKey: %s", err)
	}
	if key.State != "INACTIVE" {
		t.Errorf("State of key is %s, should be INACTIVE", key.State)
	}
}

func TestGetKey(t *testing.T) {
	key, err := CreateTestKey()
	defer DeleteTestKey(key)
	if err != nil {
		t.Errorf("Error in key creation: %s", err)
	}

	buf := new(bytes.Buffer)
	key, err = getHMACKey(buf, key.AccessID, key.ProjectID)
	if err != nil {
		t.Errorf("Error in getHMACKey: %s", err)
	}
	if key == nil {
		t.Errorf("Returned nil key.")
	}
}

func TestDeleteKey(t *testing.T) {
	key, err := CreateTestKey()
	defer DeleteTestKey(key)
	if err != nil {
		t.Errorf("Error in key creation: %s", err)
	}

	// Keys must be in INACTIVE state before deletion.
	ctx := context.Background()
	keyHandle := StorageClient.HMACKeyHandle(key.ProjectID, key.AccessID)
	keyHandle.Update(ctx, storage.HMACKeyAttrsToUpdate{State: "INACTIVE"})

	buf := new(bytes.Buffer)
	err = deleteHMACKey(buf, key.AccessID, key.ProjectID)
	if err != nil {
		t.Errorf("Error in deleteHMACKey: %s", err)
	}
	key, _ = keyHandle.Get(ctx)
	if key.State != "DELETED" {
		t.Errorf("State of key is %s, should be DELETED", key.State)
	}

}

// Create a key for testing purposes.
func CreateTestKey() (*storage.HMACKey, error) {
	ctx := context.Background()
	key, err := StorageClient.CreateHMACKey(ctx, ProjectID, ServiceAccountEmail)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Deactivate and delete the given key. Should operate as a teardown method.
func DeleteTestKey(key *storage.HMACKey) {
	ctx := context.Background()
	keyHandle := StorageClient.HMACKeyHandle(key.ProjectID, key.AccessID)
	if key.State == "ACTIVE" {
		keyHandle.Update(ctx, storage.HMACKeyAttrsToUpdate{State: "INACTIVE"})
	}
	if key.State != "DELETED" {
		keyHandle.Delete(ctx)
	}
}
