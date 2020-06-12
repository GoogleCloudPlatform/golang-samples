// Copyright 2020 Google LLC
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

package buckets

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreate(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	// Clean up bucket before running tests.
	deleteBucket(ioutil.Discard, bucketName)
	if err := createBucket(ioutil.Discard, tc.ProjectID, bucketName); err != nil {
		t.Fatalf("createBucket: %v", err)
	}
}

func TestCreateBucketClassLocation(t *testing.T) {
	tc := testutil.SystemTest(t)
	name := tc.ProjectID + "-storage-buckets-tests-attrs"

	// Clean up bucket before running the test.
	deleteBucket(ioutil.Discard, name)
	if err := createBucketClassLocation(ioutil.Discard, tc.ProjectID, name); err != nil {
		t.Fatalf("createBucketClassLocation: %v", err)
	}
	if err := deleteBucket(ioutil.Discard, name); err != nil {
		t.Fatalf("deleteBucket: %v", err)
	}
}

func TestListBuckets(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	buckets, err := listBuckets(ioutil.Discard, tc.ProjectID)
	if err != nil {
		t.Fatalf("listBuckets: %v", err)
	}

	var ok bool
	testutil.Retry(t, 5, 2*time.Second, func(r *testutil.R) { // for eventual consistency
		for _, b := range buckets {
			if b == bucketName {
				ok = true
				break
			}
		}
		if !ok {
			r.Errorf("got bucket list: %v; want %q in the list", buckets, bucketName)
		}
	})
}

func TestGetBucketMetadata(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	buf := new(bytes.Buffer)
	if _, err := getBucketMetadata(buf, bucketName); err != nil {
		t.Errorf("getBucketMetadata: %#v", err)
	}

	got := buf.String()
	if want := "BucketName:"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestIAM(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	if _, err := getBucketPolicy(ioutil.Discard, bucketName); err != nil {
		t.Errorf("getBucketPolicy: %#v", err)
	}
	if err := addBucketIAMMember(ioutil.Discard, bucketName); err != nil {
		t.Errorf("addBucketIAMMember: %v", err)
	}
	if err := removeBucketIAMMember(ioutil.Discard, bucketName); err != nil {
		t.Errorf("removeBucketIAMMember: %v", err)
	}

	// Uniform bucket-level access is required to use IAM with conditions.
	if err := enableUniformBucketLevelAccess(ioutil.Discard, bucketName); err != nil {
		t.Fatalf("enableUniformBucketLevelAccess:  %v", err)
	}

	role := "roles/storage.objectViewer"
	member := "group:cloud-logs@google.com"
	title := "title"
	description := "description"
	expression := "resource.name.startsWith(\"projects/_/buckets/bucket-name/objects/prefix-a-\")"

	if err := addBucketConditionalIAMBinding(ioutil.Discard, bucketName, role, member, title, description, expression); err != nil {
		t.Errorf("addBucketConditionalIAMBinding: %v", err)
	}
	if err := removeBucketConditionalIAMBinding(ioutil.Discard, bucketName, role, title, description, expression); err != nil {
		t.Errorf("removeBucketConditionalIAMBinding: %v", err)
	}
}

func TestRequesterPays(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	if err := enableRequesterPays(ioutil.Discard, bucketName); err != nil {
		t.Errorf("enableRequesterPays: %#v", err)
	}
	if err := disableRequesterPays(ioutil.Discard, bucketName); err != nil {
		t.Errorf("disableRequesterPays: %#v", err)
	}
	if err := getRequesterPaysStatus(ioutil.Discard, bucketName); err != nil {
		t.Errorf("getRequesterPaysStatus: %#v", err)
	}
}

func TestKMS(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	keyRingID := os.Getenv("GOLANG_SAMPLES_KMS_KEYRING")
	cryptoKeyID := os.Getenv("GOLANG_SAMPLES_KMS_CRYPTOKEY")

	if keyRingID == "" || cryptoKeyID == "" {
		t.Skip("GOLANG_SAMPLES_KMS_KEYRING and GOLANG_SAMPLES_KMS_CRYPTOKEY must be set")
	}

	kmsKeyName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", tc.ProjectID, "global", keyRingID, cryptoKeyID)
	if err := setBucketDefaultKMSKey(ioutil.Discard, bucketName, kmsKeyName); err != nil {
		t.Fatalf("setBucketDefaultKmsKey: failed to enable default kms key (%q): %v", kmsKeyName, err)
	}
}

func TestBucketLock(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	retentionPeriod := 5 * time.Second
	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := setRetentionPolicy(ioutil.Discard, bucketName, retentionPeriod); err != nil {
			r.Errorf("setRetentionPolicy: %v", err)
		}
	})

	attrs, err := getRetentionPolicy(ioutil.Discard, bucketName)
	if err != nil {
		t.Fatalf("getRetentionPolicy: %v", err)
	}
	if attrs.RetentionPolicy.RetentionPeriod != retentionPeriod {
		t.Fatalf("retention period is not the expected value (%q): %v", retentionPeriod, attrs.RetentionPolicy.RetentionPeriod)
	}
	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := enableDefaultEventBasedHold(ioutil.Discard, bucketName); err != nil {
			r.Errorf("enableDefaultEventBasedHold: %v", err)
		}
	})

	attrs, err = getDefaultEventBasedHold(ioutil.Discard, bucketName)
	if err != nil {
		t.Fatalf("getDefaultEventBasedHold: %v", err)
	}
	if !attrs.DefaultEventBasedHold {
		t.Fatalf("default event-based hold was not enabled")
	}
	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := disableDefaultEventBasedHold(ioutil.Discard, bucketName); err != nil {
			r.Errorf("disableDefaultEventBasedHold: %v", err)
		}
	})

	attrs, err = getDefaultEventBasedHold(ioutil.Discard, bucketName)
	if err != nil {
		t.Fatalf("getDefaultEventBasedHold: %v", err)
	}
	if attrs.DefaultEventBasedHold {
		t.Fatalf("default event-based hold was not disabled")
	}
	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := removeRetentionPolicy(ioutil.Discard, bucketName); err != nil {
			r.Errorf("removeRetentionPolicy: %v", err)
		}
	})

	attrs, err = getRetentionPolicy(ioutil.Discard, bucketName)
	if err != nil {
		t.Fatalf("getRetentionPolicy: %v", err)
	}
	if attrs.RetentionPolicy != nil {
		t.Fatalf("retention period to not be set")
	}
	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := setRetentionPolicy(ioutil.Discard, bucketName, retentionPeriod); err != nil {
			r.Errorf("setRetentionPolicy: %v", err)
		}
	})

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := lockRetentionPolicy(ioutil.Discard, bucketName); err != nil {
			r.Errorf("lockRetentionPolicy: %v", err)
		}
		attrs, err := getRetentionPolicy(ioutil.Discard, bucketName)
		if err != nil {
			r.Errorf("getRetentionPolicy: %v", err)
		}
		if !attrs.RetentionPolicy.IsLocked {
			r.Errorf("retention policy is not locked")
		}
	})

	time.Sleep(5 * time.Second)
	deleteBucket(ioutil.Discard, bucketName)
	time.Sleep(5 * time.Second)

	if err := createBucket(ioutil.Discard, tc.ProjectID, bucketName); err != nil {
		t.Fatalf("createBucket: %v", err)
	}
}

func TestUniformBucketLevelAccess(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := enableUniformBucketLevelAccess(ioutil.Discard, bucketName); err != nil {
			r.Errorf("enableUniformBucketLevelAccess: %v", err)
		}
	})

	attrs, err := getUniformBucketLevelAccess(ioutil.Discard, bucketName)
	if err != nil {
		t.Fatalf("getUniformBucketLevelAccess: %v", err)
	}
	if !attrs.UniformBucketLevelAccess.Enabled {
		t.Fatalf("Uniform bucket-level access was not enabled for (%q).", bucketName)
	}

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := disableUniformBucketLevelAccess(ioutil.Discard, bucketName); err != nil {
			r.Errorf("disableUniformBucketLevelAccess: %v", err)
		}
	})

	attrs, err = getUniformBucketLevelAccess(ioutil.Discard, bucketName)
	if err != nil {
		t.Fatalf("getUniformBucketLevelAccess: %v", err)
	}
	if attrs.UniformBucketLevelAccess.Enabled {
		t.Fatalf("Uniform bucket-level access was not disabled for (%q).", bucketName)
	}
}

func TestLifecycleManagement(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	ctx := context.Background()
	testutil.CleanBucket(ctx, t, tc.ProjectID, bucketName)

	if err := enableBucketLifecycleManagement(ioutil.Discard, bucketName); err != nil {
		t.Fatalf("enableBucketLifecycleManagement: %v", err)
	}

	// verify lifecycle is set
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}

	want := storage.LifecycleRule{
		Action:    storage.LifecycleAction{Type: "Delete"},
		Condition: storage.LifecycleCondition{AgeInDays: 100},
	}

	r := attrs.Lifecycle.Rules
	if len(r) != 1 {
		t.Fatalf("Length of lifecycle rules should be 1, got %d", len(r))
	}

	if !reflect.DeepEqual(r[0], want) {
		t.Fatalf("Unexpected lifecycle rule: got: %v, want: %v", r, want)
	}
}

func TestDelete(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	if err := deleteBucket(ioutil.Discard, bucketName); err != nil {
		t.Fatalf("deleteBucket: %v", err)
	}
}
