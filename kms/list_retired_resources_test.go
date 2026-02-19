// Copyright 2026 Google LLC
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

package kms

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestListRetiredResources(t *testing.T) {
	testutil.SystemTest(t)

	// Create a fresh key and version to delete.
	keyName, err := fixture.CreateSymmetricKey(fixture.KeyRingName)
	if err != nil {
		t.Fatalf("failed to create key for retired resource list: %v", err)
	}
	versionName := fmt.Sprintf("%s/cryptoKeyVersions/1", keyName)

	if err := fixture.WaitForKeyVersionReady(versionName); err != nil {
		t.Fatalf("failed to wait for key version ready: %v", err)
	}

	// Delete the version.
	var b bytes.Buffer
	if err := deleteCryptoKeyVersion(&b, versionName); err != nil {
		t.Fatalf("deleteCryptoKeyVersion: %v", err)
	}

	// Test ListRetiredResources
	b.Reset()

	// Retry a few times as listing might have eventual consistency.
	deadline := time.Now().Add(30 * time.Second)
	var listErr error
	found := false

	// List from the Location, as per proto.
	for time.Now().Before(deadline) {
		b.Reset()
		// fixture.LocationName is like "projects/p/locations/l"
		if listErr = listRetiredResources(&b, fixture.LocationName); listErr == nil {
			if strings.Contains(b.String(), versionName) {
				found = true
				break
			}
		}
		time.Sleep(2 * time.Second)
	}

	if listErr != nil {
		t.Fatalf("listRetiredResources: %v", listErr)
	}

	if !found {
		t.Errorf("listRetiredResources: expected to find %q in output:\n%s", versionName, b.String())
	}
}
