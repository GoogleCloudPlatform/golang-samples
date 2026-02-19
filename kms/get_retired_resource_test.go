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
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
)

func TestGetRetiredResource(t *testing.T) {
	testutil.SystemTest(t)

	// Create a fresh key and version to delete.
	keyName, err := fixture.CreateSymmetricKey(fixture.KeyRingName)
	if err != nil {
		t.Fatalf("failed to create key for retired resource: %v", err)
	}
	versionName := fmt.Sprintf("%s/cryptoKeyVersions/1", keyName)

	if err := fixture.WaitForKeyVersionReady(versionName); err != nil {
		t.Fatalf("failed to wait for key version ready: %v", err)
	}

	// Delete the version to create a retired resource.
	var b bytes.Buffer
	if err := deleteCryptoKeyVersion(&b, versionName); err != nil {
		t.Fatalf("deleteCryptoKeyVersion: %v", err)
	}

	// We need to find the name of the retired resource to call GetRetiredResource.
	// Since we created a FRESH key, there should be exactly one retired resource under it.

	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	deadline := time.Now().Add(60 * time.Second)
	var retiredResourceName string

	for time.Now().Before(deadline) {
		req := &kmspb.ListRetiredResourcesRequest{
			Parent: keyName,
		}
		it := client.ListRetiredResources(ctx, req)

		// Just take the first one.
		resp, err := it.Next()
		if err == iterator.Done {
			// Not found yet.
			time.Sleep(2 * time.Second)
			continue
		}
		if err != nil {
			// Retrying list on error
			time.Sleep(2 * time.Second)
			continue
		}

		retiredResourceName = resp.Name
		break
	}

	if retiredResourceName == "" {
		t.Fatalf("failed to find ANY retired resource for %s", keyName)
	}

	// Now call the sample with the correct name.
	b.Reset()
	if err := getRetiredResource(&b, retiredResourceName); err != nil {
		t.Fatalf("getRetiredResource: %v", err)
	}

	if got, want := b.String(), "Got retired resource"; !strings.Contains(got, want) {
		t.Errorf("getRetiredResource: expected %q to contain %q", got, want)
	}
}
