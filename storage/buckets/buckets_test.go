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
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	testPrefix      = "storage-buckets-test"
	bucketExpiryAge = time.Hour * 24
)

var client *storage.Client

func TestMain(m *testing.M) {
	// Initialize global vars
	tc, _ := testutil.ContextMain(m)

	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	client = c
	defer client.Close()

	// Run tests
	exit := m.Run()

	// Delete old buckets whose name begins with our test prefix
	if err := testutil.DeleteExpiredBuckets(client, tc.ProjectID, testPrefix, bucketExpiryAge); err != nil {
		// Don't fail the test if cleanup fails
		log.Printf("Post-test cleanup failed: %v", err)
	}
	os.Exit(exit)
}

func TestCreate(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := testutil.UniqueBucketName(testPrefix)
	ctx := context.Background()

	defer testutil.DeleteBucketIfExists(ctx, client, bucketName)

	if err := createBucket(ioutil.Discard, tc.ProjectID, bucketName); err != nil {
		t.Fatalf("createBucket: %v", err)
	}
}

func TestCreateBucketClassLocation(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := testutil.UniqueBucketName(testPrefix)
	ctx := context.Background()

	defer testutil.DeleteBucketIfExists(ctx, client, bucketName)

	if err := createBucketClassLocation(ioutil.Discard, tc.ProjectID, bucketName); err != nil {
		t.Fatalf("createBucketClassLocation: %v", err)
	}
}

func TestCreateBucketDualRegion(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	bucketName := testutil.UniqueBucketName(testPrefix)
	ctx := context.Background()

	defer testutil.DeleteBucketIfExists(ctx, client, bucketName)

	location := "US"
	region1 := "US-EAST1"
	region2 := "US-WEST1"
	if err := createBucketDualRegion(buf, tc.ProjectID, bucketName); err != nil {
		t.Fatalf("createBucketDualRegion: %v", err)
	}
	got := buf.String()
	if want := bucketName; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
	if want := location; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
	if want := "dual-region"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
	if want := fmt.Sprintf("%s %s", region1, region2); !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStorageClass(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

	if err := changeDefaultStorageClass(ioutil.Discard, bucketName); err != nil {
		t.Errorf("changeDefaultStorageClass: %v", err)
	}
	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	got := attrs.StorageClass
	if want := "COLDLINE"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestListBuckets(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

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
	ctx := context.Background()

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

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
	ctx := context.Background()

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

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
func TestCORSConfiguration(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

	want := []storage.CORS{
		{
			MaxAge:          time.Hour,
			Methods:         []string{"GET"},
			Origins:         []string{"some-origin.com"},
			ResponseHeaders: []string{"Content-Type"},
		},
	}
	if err := setBucketCORSConfiguration(ioutil.Discard, bucketName, want[0].MaxAge, want[0].Methods, want[0].Origins, want[0].ResponseHeaders); err != nil {
		t.Fatalf("setBucketCORSConfiguration: %v", err)
	}
	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	if !reflect.DeepEqual(attrs.CORS, want) {
		t.Fatalf("Unexpected CORS Configuration: got: %v, want: %v", attrs.CORS, want)
	}
	if err := removeBucketCORSConfiguration(ioutil.Discard, bucketName); err != nil {
		t.Fatalf("removeBucketCORSConfiguration: %v", err)
	}
	attrs, err = client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	if attrs.CORS != nil {
		t.Fatalf("Unexpected CORS Configuration: got: %v, want: %v", attrs.CORS, []storage.CORS{})
	}
}

func TestRequesterPays(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

	// Tests which update the bucket metadata must be retried in order to avoid
	// flakes from rate limits.
	testutil.Retry(t, 5, 2*time.Second, func(r *testutil.R) {
		if err := enableRequesterPays(ioutil.Discard, bucketName); err != nil {
			r.Errorf("enableRequesterPays: %#v", err)
		}
	})
	testutil.Retry(t, 5, 2*time.Second, func(r *testutil.R) {
		if err := disableRequesterPays(ioutil.Discard, bucketName); err != nil {
			r.Errorf("disableRequesterPays: %#v", err)
		}
	})
	if err := getRequesterPaysStatus(ioutil.Discard, bucketName); err != nil {
		t.Errorf("getRequesterPaysStatus: %#v", err)
	}
}

func TestKMS(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

	keyRingID := os.Getenv("GOLANG_SAMPLES_KMS_KEYRING")
	cryptoKeyID := os.Getenv("GOLANG_SAMPLES_KMS_CRYPTOKEY")

	if keyRingID == "" || cryptoKeyID == "" {
		t.Skip("GOLANG_SAMPLES_KMS_KEYRING and GOLANG_SAMPLES_KMS_CRYPTOKEY must be set")
	}

	kmsKeyName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", tc.ProjectID, "global", keyRingID, cryptoKeyID)
	testutil.Retry(t, 5, 2*time.Second, func(r *testutil.R) {
		if err := setBucketDefaultKMSKey(ioutil.Discard, bucketName, kmsKeyName); err != nil {
			r.Errorf("setBucketDefaultKMSKey: failed to enable default KMS key (%q): %v", kmsKeyName, err)
		}
	})
	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	if attrs.Encryption.DefaultKMSKeyName != kmsKeyName {
		t.Fatalf("Default KMS key was not set correctly: got %v, want %v", attrs.Encryption.DefaultKMSKeyName, kmsKeyName)
	}
	testutil.Retry(t, 5, 2*time.Second, func(r *testutil.R) {
		if err := removeBucketDefaultKMSKey(ioutil.Discard, bucketName); err != nil {
			r.Errorf("removeBucketDefaultKMSKey: failed to remove default KMS key: %v", err)
		}
	})
	attrs, err = client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	if attrs.Encryption != nil {
		t.Fatalf("Default KMS key was not removed from a bucket(%v)", bucketName)
	}
}

func TestBucketLock(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

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
	ctx := context.Background()

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

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

func TestPublicAccessPrevention(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

	if err := setPublicAccessPreventionEnforced(ioutil.Discard, bucketName); err != nil {
		t.Errorf("setPublicAccessPreventionEnforced: %v", err)
	}
	// Verify that PublicAccessPrevention was set correctly.
	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	if attrs.PublicAccessPrevention != storage.PublicAccessPreventionEnforced {
		t.Errorf("PublicAccessPrevention: got %s, want %s", attrs.PublicAccessPrevention, storage.PublicAccessPreventionEnforced)
	}

	buf := new(bytes.Buffer)
	if err := getPublicAccessPrevention(buf, bucketName); err != nil {
		t.Errorf("getPublicAccessPrevention: %v", err)
	}
	// Verify that the correct value was printed.
	got := buf.String()
	want := "Public access prevention is enforced"
	if !strings.Contains(got, want) {
		t.Errorf("getPublicAccessPrevention: got %v, want %v", got, want)
	}

	if err := setPublicAccessPreventionInherited(ioutil.Discard, bucketName); err != nil {
		t.Errorf("setPublicAccessPreventionInherited: %v", err)
	}
	// Verify that PublicAccessPrevention was set correctly.
	attrs, err = client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	if attrs.PublicAccessPrevention != storage.PublicAccessPreventionInherited {
		t.Errorf("PublicAccessPrevention: got %s, want %s", attrs.PublicAccessPrevention, storage.PublicAccessPreventionInherited)
	}

}

func TestLifecycleManagement(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

	if err := enableBucketLifecycleManagement(ioutil.Discard, bucketName); err != nil {
		t.Fatalf("enableBucketLifecycleManagement: %v", err)
	}

	// Verify lifecycle is set
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

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := disableBucketLifecycleManagement(ioutil.Discard, bucketName); err != nil {
			r.Errorf("disableBucketLifecycleManagement: %v", err)
		}
	})

	attrs, err = client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}

	if n := len(attrs.Lifecycle.Rules); n != 0 {
		t.Fatalf("Length of lifecycle rules should be 0, got %d", n)
	}
}

func TestBucketLabel(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

	labelName := "label-name"
	labelValue := "label-value"
	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := addBucketLabel(ioutil.Discard, bucketName, labelName, labelValue); err != nil {
			r.Errorf("addBucketLabel: %v", err)
		}
	})
	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	if got, ok := attrs.Labels[labelName]; ok {
		if got != labelValue {
			t.Fatalf("The label(%q) was set incorrectly on a bucket(%v): got value %v, want value %v", labelName, bucketName, got, labelValue)
		}
	} else {
		t.Fatalf("The label(%q) was not set on a bucket(%v)", labelName, bucketName)
	}
	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := removeBucketLabel(ioutil.Discard, bucketName, labelName); err != nil {
			r.Errorf("removeBucketLabel: %v", err)
		}
	})
	attrs, err = client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	if _, ok := attrs.Labels[labelName]; ok {
		t.Fatalf("The label(%q) was not removed from a bucket(%v)", labelName, bucketName)
	}
}

func TestBucketWebsiteInfo(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

	index := "index.html"
	notFoundPage := "404.html"
	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := setBucketWebsiteInfo(ioutil.Discard, bucketName, index, notFoundPage); err != nil {
			r.Errorf("setBucketWebsiteInfo: %v", err)
		}
	})
	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	if attrs.Website.MainPageSuffix != index {
		t.Fatalf("got index page: %v, want %v", attrs.Website.MainPageSuffix, index)
	}
	if attrs.Website.NotFoundPage != notFoundPage {
		t.Fatalf("got not found page: %v, want %v", attrs.Website.NotFoundPage, notFoundPage)
	}
}

func TestSetBucketPublicIAM(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

	if err := setBucketPublicIAM(ioutil.Discard, bucketName); err != nil {
		t.Fatalf("setBucketPublicIAM: %v", err)
	}
	policy, err := client.Bucket(bucketName).IAM().V3().Policy(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).IAM().V3().Policy: %v", bucketName, err)
	}
	want := new(iam.Policy3)
	want.Bindings = append(want.Bindings, &iampb.Binding{
		Role:    "roles/storage.objectViewer",
		Members: []string{iam.AllUsers},
	})
	if !reflect.DeepEqual(policy.Bindings[len((policy.Bindings))-1], want.Bindings[0]) {
		t.Fatalf("Public policy was not set: \ngot: %v, \nwant: %v\n", policy.Bindings[len((policy.Bindings))-1], want.Bindings[0])
	}
}

func TestDelete(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)

	if err := deleteBucket(ioutil.Discard, bucketName); err != nil {
		t.Fatalf("deleteBucket: %v", err)
	}
}

// TestRPO tests the following samples:
// createBucketTurboReplication, setRPODefault, setRPOAsyncTurbo, getRPO
func TestRPO(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := testutil.UniqueBucketName(testPrefix)
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	// Clean up bucket before running the test
	if err := testutil.DeleteBucketIfExists(ctx, client, bucketName); err != nil {
		t.Fatalf("Error deleting bucket: %v", err)
	}

	location := "NAM4" // must be dual-region
	if err := createBucketTurboReplication(ioutil.Discard, tc.ProjectID, bucketName, location); err != nil {
		t.Fatalf("createBucketTurboReplication: %v", err)
	}

	// Verify that RPO was set correctly on creation
	attrs, err := bucket.Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	if attrs.RPO != storage.RPOAsyncTurbo {
		t.Errorf("createBucketTurboReplication: got %s, want %s", attrs.RPO, storage.RPOAsyncTurbo)
	}

	// Test disable turbo replication:
	if err := setRPODefault(ioutil.Discard, bucketName); err != nil {
		t.Errorf("setRPODefault: %v", err)
	}
	// Verify that RPO was set correctly
	attrs, err = bucket.Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	if attrs.RPO != storage.RPODefault {
		t.Errorf("setRPODefault: got %s, want %s", attrs.RPO, storage.RPODefault)
	}

	// Test enable turbo replication:
	if err := setRPOAsyncTurbo(ioutil.Discard, bucketName); err != nil {
		t.Fatalf("setRPOAsyncTurbo: %v", err)
	}

	// Verify that RPO was set correctly
	attrs, err = bucket.Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	if attrs.RPO != storage.RPOAsyncTurbo {
		t.Errorf("setRPOAsyncTurbo: got %s, want %s", attrs.RPO, storage.RPOAsyncTurbo)
	}

	// Test get turbo replication:
	buf := new(bytes.Buffer)
	if err := getRPO(buf, bucketName); err != nil {
		t.Errorf("getRPO: %v", err)
	}
	// Verify that the correct value was printed
	got := buf.String()
	want := "RPO is ASYNC_TURBO"
	if !strings.Contains(got, want) {
		t.Errorf("getRPO: got %v, want %v", got, want)
	}

	if err := deleteBucket(ioutil.Discard, bucketName); err != nil {
		t.Fatalf("deleteBucket: %v", err)
	}
}

// TestAutoclass tests the following samples:
// getAutoclass, setAutoclass
func TestAutoclass(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	bucketName := testutil.UniqueBucketName(testPrefix)
	defer testutil.DeleteBucketIfExists(ctx, client, bucketName)

	// Test create new bucket with Autoclass enabled.
	autoclassConfig := &storage.BucketAttrs{
		Autoclass: &storage.Autoclass{
			Enabled: true,
		},
	}
	bucket := client.Bucket(bucketName)
	if err := bucket.Create(ctx, tc.ProjectID, autoclassConfig); err != nil {
		t.Fatalf("Bucket creation failed: %v", err)
	}

	// Test get Autoclass config.
	buf := new(bytes.Buffer)
	if err := getAutoclass(buf, bucketName); err != nil {
		t.Errorf("getAutoclass: %#v", err)
	}
	if got, want := buf.String(), "Autoclass enabled was set to true"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}

	// Test set Autoclass terminal storage class to "ARCHIVE".
	if err := setAutoclass(buf, bucketName); err != nil {
		t.Errorf("setAutoclass: %#v", err)
	}
	if got, want := buf.String(), "Autoclass terminal storage class was last updated to ARCHIVE"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCreateBucketHierarchicalNamespace(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucketName := testutil.UniqueBucketName(testPrefix)
	defer testutil.DeleteBucketIfExists(ctx, client, bucketName)

	// Test creating new bucket with HNS enabled.
	buf := new(bytes.Buffer)
	if err := createBucketHierarchicalNamespace(buf, tc.ProjectID, bucketName); err != nil {
		t.Fatalf("createBucketHierarchicalNamespace: %v", err)
	}

	if got, want := buf.String(), "Created bucket"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}

	// Verify that HNS was set as expected.
	attrs, err := client.Bucket(bucketName).Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketName, err)
	}
	if got, want := (attrs.HierarchicalNamespace), (&storage.HierarchicalNamespace{Enabled: true}); got == nil || !got.Enabled {
		t.Errorf("Attrs.HierarchicalNamespace: got %v, want %v", got, want)
	}
}

func TestCreateBucketObjectRetention(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := testutil.UniqueBucketName(testPrefix)
	ctx := context.Background()

	defer testutil.DeleteBucketIfExists(ctx, client, bucketName)

	buf := new(bytes.Buffer)

	if err := createBucketObjectRetention(buf, tc.ProjectID, bucketName); err != nil {
		t.Fatalf("createBucketObjectRetention: %v", err)
	}

	if got, want := buf.String(), "Enabled"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
