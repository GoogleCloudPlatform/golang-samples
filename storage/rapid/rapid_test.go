// Copyright 2025 Google LLC
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

package rapid

import (
	"bytes"
	"context"
	"log"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"cloud.google.com/go/storage/experimental"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/google/uuid"
)

const (
	testPrefix        = "storage-objects-test"
	testZonalLocation = "us-west4"
	testZonalZone     = "us-west4-a"
)

var (
	zonalBucketName string
	client          *storage.Client
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Skip tests by default for now, until b/452725162 is resolved
	if os.Getenv("STORAGE_RUN_RAPID_TESTS") == "" {
		os.Exit(0)
	}

	// Create fixture client & bucket to use across tests.
	tc, _ := testutil.ContextMain(m)
	var err error
	client, err = storage.NewGRPCClient(context.Background(), experimental.WithZonalBucketAPIs())
	if err != nil {
		log.Fatalf("storage.NewGRPCClient: %v", err)
	}
	zonalBucketName = strings.Join([]string{testPrefix, uuid.NewString()}, "-")
	if err := client.Bucket(zonalBucketName).Create(ctx, tc.ProjectID, &storage.BucketAttrs{
		Location: testZonalLocation,
		CustomPlacementConfig: &storage.CustomPlacementConfig{
			DataLocations: []string{testZonalZone},
		},
		StorageClass: "RAPID",
		HierarchicalNamespace: &storage.HierarchicalNamespace{
			Enabled: true,
		},
		UniformBucketLevelAccess: storage.UniformBucketLevelAccess{
			Enabled: true,
		},
	}); err != nil {
		log.Fatalf("BucketHandle.Create: %v", err)

	}
	// Run tests.
	exit := m.Run()

	// Cleanup bucket and objects.
	if err := testutil.DeleteBucketIfExists(ctx, client, zonalBucketName); err != nil {
		log.Printf("deleting bucket: %v", err)
	}
	os.Exit(exit)
}

func TestCreateAndWriteAppendableObject(t *testing.T) {
	var b bytes.Buffer
	object := "obj-appendable"
	if err := createAndWriteAppendableObject(&b, zonalBucketName, object); err != nil {
		t.Fatalf("running sample: %v, output: %v", err, b.String())
	}

	// Check that object was created & is unfinalized
	attrs, err := client.Bucket(zonalBucketName).Object(object).Attrs(context.Background())
	if err != nil {
		t.Fatalf("object.Attrs: %v", err)
	}
	if !attrs.Finalized.IsZero() {
		t.Errorf("got finalized object, want unfinalized")
	}
}

func TestFinalizeAppendableObject(t *testing.T) {
	var b bytes.Buffer
	object := "obj-finalize"
	if err := finalizeAppendableObject(&b, zonalBucketName, object); err != nil {
		t.Fatalf("running sample: %v, output: %v", err, b.String())
	}

	// Check that object was created & is finalized
	attrs, err := client.Bucket(zonalBucketName).Object(object).Attrs(context.Background())
	if err != nil {
		t.Fatalf("object.Attrs: %v", err)
	}
	if attrs.Finalized.IsZero() {
		t.Errorf("got unfinalized object, want finalized")
	}
}

func TestPauseAndResumeAppendableUpload(t *testing.T) {
	var b bytes.Buffer
	object := "obj-pause"
	if err := pauseAndResumeAppendableUpload(&b, zonalBucketName, object); err != nil {
		t.Fatalf("running sample: %v, output: %v", err, b.String())
	}

	// Check that object was created & is finalized
	attrs, err := client.Bucket(zonalBucketName).Object(object).Attrs(context.Background())
	if err != nil {
		t.Fatalf("object.Attrs: %v", err)
	}
	if attrs.Finalized.IsZero() {
		t.Errorf("got unfinalized object, want finalized")
	}
}
