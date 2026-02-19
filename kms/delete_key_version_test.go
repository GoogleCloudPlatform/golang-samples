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

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestDeleteCryptoKeyVersion(t *testing.T) {
	testutil.SystemTest(t)

	// Create a fresh key and version to delete.
	keyName, err := fixture.CreateSymmetricKey(fixture.KeyRingName)
	if err != nil {
		t.Fatalf("failed to create key for version deletion: %v", err)
	}
	versionName := fmt.Sprintf("%s/cryptoKeyVersions/1", keyName)

	// Wait for the version to be ready before deleting (good practice).
	if err := fixture.WaitForKeyVersionReady(versionName); err != nil {
		t.Fatalf("failed to wait for key version ready: %v", err)
	}

	var b bytes.Buffer
	if err := deleteCryptoKeyVersion(&b, versionName); err != nil {
		t.Fatalf("deleteCryptoKeyVersion: %v", err)
	}

	if got, want := b.String(), "Deleted crypto key version"; !strings.Contains(got, want) {
		t.Errorf("deleteCryptoKeyVersion: expected %q to contain %q", got, want)
	}
}
