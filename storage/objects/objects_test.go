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

package objects

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	testPrefix      = "storage-objects-test"
	bucketExpiryAge = time.Hour * 24
)

func TestMain(m *testing.M) {
	// Run tests
	exit := m.Run()

	// Delete old buckets whose name begins with our test prefix
	tc, _ := testutil.ContextMain(m)

	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	defer c.Close()

	if err := testutil.DeleteExpiredBuckets(c, tc.ProjectID, testPrefix, bucketExpiryAge); err != nil {
		// Don't fail the test if cleanup fails
		log.Printf("Post-test cleanup failed: %v", err)
	}
	os.Exit(exit)
}

// TestObjects runs most of the samples tests of the package.
func TestObjects(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	dir, err := os.MkdirTemp("", "objectsTestTempDir")
	if err != nil {
		t.Fatalf("os.MkdirTemp: %v", err)
	}
	defer os.RemoveAll(dir) // clean up

	var (
		bucket           = testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)
		dstBucket        = testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)
		bucketVersioning = testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)
		object1          = "foo.txt"
		object2          = "foo/a.txt"
		object3          = "bar.txt"
		dstObj           = "foobar.txt"
	)

	if err := enableVersioning(io.Discard, bucketVersioning); err != nil {
		t.Fatalf("enableVersioning: %v", err)
	}

	if err := uploadFile(io.Discard, bucket, object1); err != nil {
		t.Fatalf("uploadFile(%q): %v", object1, err)
	}
	if err := uploadFile(io.Discard, bucket, object2); err != nil {
		t.Fatalf("uploadFile(%q): %v", object2, err)
	}

	if err := streamFileUpload(io.Discard, bucketVersioning, object1); err != nil {
		t.Fatalf("streamFileUpload(%q): %v", object1, err)
	}
	// Check enableVersioning correctly work.
	bkt := client.Bucket(bucketVersioning)
	bAttrs, err := bkt.Attrs(ctx)
	if !bAttrs.VersioningEnabled {
		t.Fatalf("object versioning is not enabled")
	}
	obj := bkt.Object(object1)
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Object(%q).Attrs: %v", bucketVersioning, object1, err)
	}
	// Keep the original generation of object1 before re-uploading
	// to use in the versioning samples.
	gen := attrs.Generation
	if err := streamFileUpload(io.Discard, bucketVersioning, object1); err != nil {
		t.Fatalf("streamFileUpload(%q): %v", object1, err)
	}

	{
		// Should only show "foo/a.txt", not "foo.txt"
		var buf bytes.Buffer
		if err := listFiles(&buf, bucket); err != nil {
			t.Fatalf("listFiles: %v", err)
		}
		if got, want := buf.String(), object1; !strings.Contains(got, want) {
			t.Errorf("List() got %q; want to contain %q", got, want)
		}
		if got, want := buf.String(), object2; !strings.Contains(got, want) {
			t.Errorf("List() got %q; want to contain %q", got, want)
		}
	}

	{
		// Should only show "foo/a.txt", not "foo.txt"
		const prefix = "foo/"
		var buf bytes.Buffer
		if err := listFilesWithPrefix(&buf, bucket, prefix, ""); err != nil {
			t.Fatalf("listFilesWithPrefix: %v", err)
		}
		if got, want := buf.String(), object1; strings.Contains(got, want) {
			t.Errorf("List(%q) got %q; want NOT to contain %q", prefix, got, want)
		}
		if got, want := buf.String(), object2; !strings.Contains(got, want) {
			t.Errorf("List(%q) got %q; want to contain %q", prefix, got, want)
		}
	}

	{
		// Should show 2 versions of foo.txt
		var buf bytes.Buffer
		if err := listFilesAllVersion(&buf, bucketVersioning); err != nil {
			t.Fatalf("listFilesAllVersion: %v", err)
		}

		i := 0
		for _, line := range strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n") {
			if got, want := line, object1; !strings.Contains(got, want) {
				t.Errorf("List(Versions: true) got %q; want to contain %q", got, want)
			}
			i++
		}
		if i != 2 {
			t.Errorf("listFilesAllVersion should show 2 versions of foo.txt; got %d", i)
		}
	}

	{
		if err := downloadUsingRequesterPays(io.Discard, bucket, object1, tc.ProjectID); err != nil {
			t.Errorf("downloadUsingRequesterPays: %v", err)
		}
	}
	t.Run("changeObjectStorageClass", func(t *testing.T) {
		bkt := client.Bucket(bucket)
		obj := bkt.Object(object1)
		if err := changeObjectStorageClass(io.Discard, bucket, object1); err != nil {
			t.Errorf("changeObjectStorageClass: %v", err)
		}
		wantStorageClass := "COLDLINE"
		oattrs, err := obj.Attrs(ctx)
		if err != nil {
			t.Errorf("obj.Attrs: %v", err)
		}
		if oattrs.StorageClass != wantStorageClass {
			t.Errorf("object storage class: got %q, want %q", oattrs.StorageClass, wantStorageClass)
		}
	})
	if err := copyOldVersionOfObject(io.Discard, bucketVersioning, object1, object3, gen); err != nil {
		t.Fatalf("copyOldVersionOfObject: %v", err)
	}
	// Delete the first version of an object1 for a bucketVersioning.
	if err := deleteOldVersionOfObject(io.Discard, bucketVersioning, object1, gen); err != nil {
		t.Fatalf("deleteOldVersionOfObject: %v", err)
	}
	data, err := downloadFileIntoMemory(io.Discard, bucket, object1)
	if err != nil {
		t.Fatalf("downloadFileIntoMemory: %v", err)
	}
	if got, want := string(data), "Hello\nworld"; got != want {
		t.Errorf("contents = %q; want %q", got, want)
	}

	t.Run("setMetadata", func(t *testing.T) {
		bkt := client.Bucket(bucket)
		obj := bkt.Object(object1)
		err = setMetadata(io.Discard, bucket, object1)
		if err != nil {
			t.Errorf("setMetadata: %v", err)
		}
		attrs, err := obj.Attrs(ctx)
		if err != nil {
			t.Errorf("object.Attrs: %v", err)
		}
		if got, want := attrs.Metadata["keyToAddOrUpdate"], "value"; got != want {
			t.Errorf("object content = %q; want %q", got, want)
		}
	})
	_, err = getMetadata(io.Discard, bucket, object1)
	if err != nil {
		t.Errorf("getMetadata: %v", err)
	}
	t.Run("publicFile", func(t *testing.T) {
		t.Skip("Skipping due to project permissions changes, see: b/445769988")
		if err := makePublic(io.Discard, bucket, object1); err != nil {
			t.Errorf("makePublic: %v", err)
		}
		data, err = downloadPublicFile(io.Discard, bucket, object1)
		if err != nil {
			t.Fatalf("downloadPublicFile: %v", err)
		}
		if got, want := string(data), "Hello\nworld"; got != want {
			t.Errorf("contents = %q; want %q", got, want)
		}
	})

	t.Run("downloadByteRange", func(t *testing.T) {
		destination := filepath.Join(dir, "fileDownloadByteRangeDestination.txt")
		err = downloadByteRange(io.Discard, bucket, object1, 1, 4, destination)
		if err != nil {
			t.Fatalf("downloadFile: %v", err)
		}
		data, err := os.ReadFile(destination)
		if err != nil {
			t.Fatalf("os.ReadFile: %v", err)
		}
		if got, want := string(data), "ell"; got != want {
			t.Errorf("contents = %q; want %q", got, want)
		}
	})

	t.Run("downloadFile", func(t *testing.T) {
		destination := filepath.Join(dir, "fileDownloadDestination.txt")
		err = downloadFile(io.Discard, bucket, object1, destination)
		if err != nil {
			t.Fatalf("downloadFile: %v", err)
		}
		data, err := os.ReadFile(destination)
		if err != nil {
			t.Fatalf("os.ReadFile: %v", err)
		}
		if got, want := string(data), "Hello\nworld"; got != want {
			t.Errorf("contents = %q; want %q", got, want)
		}
	})

	err = moveFile(io.Discard, bucket, object1)
	if err != nil {
		t.Fatalf("moveFile: %v", err)
	}
	// object1's new name.
	object1 = object1 + "-rename"

	if err := copyFile(io.Discard, dstBucket, bucket, object1); err != nil {
		t.Errorf("copyFile: %v", err)
	}
	t.Run("composeFile", func(t *testing.T) {
		if err := composeFile(io.Discard, bucket, object1, object2, dstObj); err != nil {
			t.Errorf("composeFile: %v", err)
		}
		bkt := client.Bucket(bucket)
		obj := bkt.Object(dstObj)
		_, err = obj.Attrs(ctx)
		if err == storage.ErrObjectNotExist {
			t.Errorf("Destination object was not created")
		} else if err != nil {
			t.Errorf("object.Attrs: %v", err)
		}
	})

	if err := deleteFile(io.Discard, bucket, object1); err != nil {
		t.Errorf("deleteFile: %v", err)
	}

	key := []byte("my-secret-AES-256-encryption-key")
	newKey := []byte("My-secret-AES-256-encryption-key")

	if err := generateEncryptionKey(io.Discard); err != nil {
		t.Errorf("generateEncryptionKey: %v", err)
	}
	if err := uploadEncryptedFile(io.Discard, bucket, object1, key); err != nil {
		t.Errorf("uploadEncryptedFile: %v", err)
	}
	data, err = downloadEncryptedFile(io.Discard, bucket, object1, key)
	if err != nil {
		t.Errorf("downloadEncryptedFile: %v", err)
	}
	if got, want := string(data), "top secret"; got != want {
		t.Errorf("object content = %q; want %q", got, want)
	}
	if err := rotateEncryptionKey(io.Discard, bucket, object1, key, newKey); err != nil {
		t.Errorf("rotateEncryptionKey: %v", err)
	}
	o := client.Bucket(bucket).Object(dstObj)
	if err := o.Delete(ctx); err != nil {
		t.Errorf("Object(%q).Delete: %v", dstObj, err)
	}
	if err := disableVersioning(io.Discard, bucketVersioning); err != nil {
		t.Fatalf("disableVersioning: %v", err)
	}
	bAttrs, err = bkt.Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketVersioning, err)
	}
	if bAttrs.VersioningEnabled {
		t.Fatalf("object versioning is not disabled")
	}
}

func TestKMSObjects(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	keyRingID := os.Getenv("GOLANG_SAMPLES_KMS_KEYRING")
	cryptoKeyID := os.Getenv("GOLANG_SAMPLES_KMS_CRYPTOKEY")
	if keyRingID == "" || cryptoKeyID == "" {
		t.Skip("GOLANG_SAMPLES_KMS_KEYRING and GOLANG_SAMPLES_KMS_CRYPTOKEY must be set")
	}

	bucket := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)
	object := "foo.txt"

	kmsKeyName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", tc.ProjectID, "global", keyRingID, cryptoKeyID)
	t.Run("сhangeObjectCSEKtoKMS", func(t *testing.T) {
		object1 := "foo1.txt"
		key := []byte("my-secret-AES-256-encryption-key")
		obj := client.Bucket(bucket).Object(object1)

		testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
			wc := obj.Key(key).NewWriter(ctx)
			if _, err := wc.Write([]byte("top secret")); err != nil {
				r.Errorf("Writer.Write: %v", err)
			}
			if err := wc.Close(); err != nil {
				r.Errorf("Writer.Close: %v", err)
			}
		})
		if err := сhangeObjectCSEKToKMS(io.Discard, bucket, object1, key, kmsKeyName); err != nil {
			t.Errorf("сhangeObjectCSEKtoKMS: %v", err)
		}
		attrs, err := obj.Attrs(ctx)
		if err != nil {
			t.Errorf("obj.Attrs: %v", err)
		}
		if got, want := attrs.KMSKeyName, kmsKeyName; !strings.Contains(got, want) {
			t.Errorf("attrs.KMSKeyName expected %q to contain %q", got, want)
		}
	})

	if err := uploadWithKMSKey(io.Discard, bucket, object, kmsKeyName); err != nil {
		t.Errorf("uploadWithKMSKey: %v", err)
	}
}

func TestV4SignedURL(t *testing.T) {
	t.Skip("Skipping due to project permissions changes, see: b/445769988")
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)
	objectName := "foo.txt"

	// Generate PUT URL.
	putBuf := new(bytes.Buffer)
	putURL, err := generateV4PutObjectSignedURL(putBuf, bucketName, objectName)
	if err != nil {
		t.Errorf("generateV4PutObjectSignedURL: %v", err)
	}
	got := putBuf.String()
	if want := "Generated PUT signed URL:"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}

	// Generate GET URL.
	getBuf := new(bytes.Buffer)
	getURL, err := generateV4GetObjectSignedURL(getBuf, bucketName, objectName)
	if err != nil {
		t.Errorf("generateV4GetObjectSignedURL: %v", err)
	}
	got = getBuf.String()
	if want := "Generated GET signed URL:"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}

	// Create PUT request.
	httpClient := &http.Client{}
	request, err := http.NewRequest("PUT", putURL, strings.NewReader("hello world"))
	if err != nil {
		t.Fatalf("failed to compose HTTP request: %v", err)
	}
	request.ContentLength = 11
	request.Header.Set("Content-Type", "application/octet-stream")

	// Test PUT and GET requests.
	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		_, err = httpClient.Do(request)
		if err != nil {
			r.Errorf("httpClient.Do: %v", err)
		}

		response, err := http.Get(getURL)
		if err != nil {
			r.Errorf("http.Get: %v", err)
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			r.Errorf("io.ReadAll: %v", err)
		}

		if got, want := string(body), "hello world"; got != want {
			r.Errorf("object content = %q; want %q", got, want)
		}
	})
}

func TestPostPolicyV4(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)
	objectName := "foo.txt"

	putBuf := new(bytes.Buffer)
	policy, err := generateSignedPostPolicyV4(putBuf, bucketName, objectName)
	if err != nil {
		t.Fatalf("generateSignedPostPolicyV4: %v", err)
	}
	got := putBuf.String()
	if want := "<form action="; !strings.HasPrefix(got, want) {
		t.Errorf("got output %q, should start with %q", got, want)
	}
	missing := false
	for k, v := range policy.Fields {
		if !strings.Contains(got, k) || !strings.Contains(got, v) {
			t.Errorf("output missing form field %v: %v", k, v)
			missing = true
		}
	}
	if missing {
		t.Fatalf("got output %q", got)
	}

	// The signed post policy allows an unauthenticated client to make a POST to
	// the bucket using policy.URL and a form containing the values from
	// policy.Fields. We test that this actually works against the live service
	// using the policy generated by the sample.

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		// Create a form using the returned post policy.
		formBuf := new(bytes.Buffer)
		mw := multipart.NewWriter(formBuf)
		for fieldName, value := range policy.Fields {
			if err := mw.WriteField(fieldName, value); err != nil {
				t.Errorf("writing form: %v", err)
			}
		}

		// Create a file for upload.
		fileBody := bytes.Repeat([]byte("z"), 25)
		mf, err := mw.CreateFormFile("file", "bar.txt")
		if err != nil {
			t.Fatalf("CreateFormFile: %v", err)
		}
		if _, err := mf.Write(fileBody); err != nil {
			t.Fatalf("Write: %v", err)
		}
		if err := mw.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}

		// Compose the HTTP request.
		req, err := http.NewRequest("POST", policy.URL, formBuf)
		if err != nil {
			t.Errorf("failed to compose HTTP request: %v", err)
		}
		req.Header.Set("Content-Type", mw.FormDataContentType())

		// Dump the request for logging.
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			r.Logf("requestDump: %v", err)
		}

		// Make request.
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			r.Errorf("client.Do: %v", err)
		}
		if g, w := res.StatusCode, 204; g != w {
			responseDump, _ := httputil.DumpResponse(res, true)
			r.Errorf("status code in response mismatch: got %d want %d\nRequest: %v\n\nResponse: %s\n",
				g, w, string(requestDump), responseDump)
		}

		io.Copy(io.Discard, res.Body)
		if err := res.Body.Close(); err != nil {
			r.Errorf("Body.Close: %v", err)
		}
	})

	// Verify that the file was uploaded by reading back its attributes.
	bkt := client.Bucket(bucketName)
	obj := bkt.Object(objectName)

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		t.Fatalf("Failed to retrieve attributes: %v", err)
	}
	if attrs.Name != objectName {
		t.Errorf("object name: got %q, want %q", attrs.Name, objectName)
	}
}

func TestObjectBucketLock(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	var (
		bucketName      = testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)
		objectName      = "foo.txt"
		retentionPeriod = 5 * time.Second
	)

	bucket := client.Bucket(bucketName)
	bucketAttrs, err := bucket.Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Update: %v", bucketName, err)
	}

	if err := uploadFile(io.Discard, bucketName, objectName); err != nil {
		t.Fatalf("uploadFile(%q): %v", objectName, err)
	}
	// Updating a bucket is conditionally idempotent, so we set metageneration match and let the library handle the retry
	if _, err := bucket.If(storage.BucketConditions{MetagenerationMatch: bucketAttrs.MetaGeneration}).
		Update(ctx, storage.BucketAttrsToUpdate{
			RetentionPolicy: &storage.RetentionPolicy{
				RetentionPeriod: retentionPeriod,
			},
		}); err != nil {
		t.Errorf("Bucket(%q).Update: %v", bucketName, err)
	}
	if err := setEventBasedHold(io.Discard, bucketName, objectName); err != nil {
		t.Errorf("setEventBasedHold(%q, %q): %v", bucketName, objectName, err)
	}
	oAttrs, err := getMetadata(io.Discard, bucketName, objectName)
	if err != nil {
		t.Errorf("getMetadata: %v", err)
	}
	if !oAttrs.EventBasedHold {
		t.Errorf("event-based hold is not enabled")
	}
	if err := releaseEventBasedHold(io.Discard, bucketName, objectName); err != nil {
		t.Errorf("releaseEventBasedHold(%q, %q): %v", bucketName, objectName, err)
	}
	oAttrs, err = getMetadata(io.Discard, bucketName, objectName)
	if err != nil {
		t.Errorf("getMetadata: %v", err)
	}
	if oAttrs.EventBasedHold {
		t.Errorf("event-based hold is not disabled")
	}

	bucketAttrs, err = bucket.Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Update: %v", bucketName, err)
	}
	// Updating a bucket is conditionally idempotent, so we set metageneration match and let the library handle the retry
	if _, err := bucket.If(storage.BucketConditions{MetagenerationMatch: bucketAttrs.MetaGeneration}).
		Update(ctx, storage.BucketAttrsToUpdate{
			RetentionPolicy: &storage.RetentionPolicy{},
		}); err != nil {
		t.Errorf("Bucket(%q).Update: %v", bucketName, err)
	}
	if err := setTemporaryHold(io.Discard, bucketName, objectName); err != nil {
		t.Errorf("setTemporaryHold(%q, %q): %v", bucketName, objectName, err)
	}
	oAttrs, err = getMetadata(io.Discard, bucketName, objectName)
	if err != nil {
		t.Errorf("getMetadata: %v", err)
	}
	if !oAttrs.TemporaryHold {
		t.Errorf("temporary hold is not disabled")
	}
	if err := releaseTemporaryHold(io.Discard, bucketName, objectName); err != nil {
		t.Errorf("releaseTemporaryHold(%q, %q): %v", bucketName, objectName, err)
	}
	oAttrs, err = getMetadata(io.Discard, bucketName, objectName)
	if err != nil {
		t.Errorf("getMetadata: %v", err)
	}
	if oAttrs.TemporaryHold {
		t.Errorf("temporary hold is not disabled")
	}
}

func TestObjectRetention(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	start := time.Now()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	var (
		bucketName = testutil.UniqueBucketName(testPrefix)
		objectName = "foo.txt"
	)

	bucket := client.Bucket(bucketName).SetObjectRetention(true)
	if err := bucket.Create(ctx, tc.ProjectID, nil); err != nil {
		t.Fatalf("Bucket(%q).Create: %v", bucketName, err)
	}
	defer testutil.DeleteBucketIfExists(ctx, client, bucketName)

	if err := uploadFile(io.Discard, bucketName, objectName); err != nil {
		t.Fatalf("uploadFile(%q): %v", objectName, err)
	}

	err = setObjectRetentionPolicy(io.Discard, bucketName, objectName)
	if err != nil {
		t.Errorf("setObjectRetention: %v", err)
	}
	attrs, err := bucket.Object(objectName).Attrs(ctx)
	if err != nil {
		t.Errorf("object.Attrs: %v", err)
	}

	if attrs.Retention == nil {
		t.Errorf("mismatching retention config, got nil, wanted %+v", attrs.Retention)
	}

	if got, want := attrs.Retention.RetainUntil, start.Add(time.Hour*24*9); got.Before(want) {
		t.Errorf("retention time should be more than 9 days from the start of the test; got %v, want after %v", got, want)
	}
	if got, want := attrs.Retention.RetainUntil, start.Add(time.Hour*24*10); got.After(want) {
		t.Errorf("retention time should be less than 10 days from the start of the test; got %v, want sooner than %v", got, want)
	}
}

func TestListSoftDeletedObjects(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	var (
		bucketName = testutil.UniqueBucketName(testPrefix)
		objectName = "soft-deleted-object.txt"
	)

	bucket := client.Bucket(bucketName)
	if err := bucket.Create(ctx, tc.ProjectID, &storage.BucketAttrs{SoftDeletePolicy: &storage.SoftDeletePolicy{
		RetentionDuration: 10 * 24 * time.Hour, // 10 days in hours
	}}); err != nil {
		t.Fatalf("Bucket.Create(%q): %v", bucketName, err)
	}
	defer testutil.DeleteBucketIfExists(ctx, client, bucketName)

	// Upload the object to the bucket.
	if err := uploadFile(io.Discard, bucketName, objectName); err != nil {
		t.Fatalf("uploadFile(%q): %v", objectName, err)
	}

	obj := client.Bucket(bucketName).Object(objectName)
	// Simulate soft deletion by deleting the object.
	if err := obj.Delete(ctx); err != nil {
		t.Fatalf("Object(%q).Delete: %v", objectName, err)
	}

	var buf bytes.Buffer
	if err := listSoftDeletedObjects(&buf, bucketName); err != nil {
		t.Fatalf("listSoftDeletedObjects: %v", err)
	}
	// Verify the output was printed as expected.
	got := buf.String()
	want := fmt.Sprintf("Soft-deleted object: %s\n", objectName)
	if !strings.HasPrefix(got, want) {
		t.Errorf("Output mismatch: got %q, want %q", got, want)
	}
}

func TestRestoreSoftDeletedObject(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	var (
		bucketName = testutil.UniqueBucketName(testPrefix)
		objectName = "soft-deleted-object.txt"
	)

	bucket := client.Bucket(bucketName)
	if err := bucket.Create(ctx, tc.ProjectID, &storage.BucketAttrs{SoftDeletePolicy: &storage.SoftDeletePolicy{
		RetentionDuration: 10 * 24 * time.Hour, // 10 days in hours
	}}); err != nil {
		t.Fatalf("Bucket.Create(%q): %v", bucketName, err)
	}
	defer testutil.DeleteBucketIfExists(ctx, client, bucketName)

	// Upload the object to the bucket.
	if err := uploadFile(io.Discard, bucketName, objectName); err != nil {
		t.Fatalf("uploadFile(%q): %v", objectName, err)
	}

	// Get object attributes to retrieve the generation before deleting the object.
	obj := client.Bucket(bucketName).Object(objectName)
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		t.Fatalf("Object(%q).Attrs: %v", objectName, err)
	}
	generation := attrs.Generation
	// Simulate soft deletion by deleting the object.
	if err := obj.Delete(ctx); err != nil {
		t.Fatalf("Object(%q).Delete: %v", objectName, err)
	}

	var buf bytes.Buffer
	if err := restoreSoftDeletedObject(&buf, bucketName, objectName, generation); err != nil {
		t.Fatalf("restoreSoftDeletedObject: %v", err)
	}
	if !strings.Contains(buf.String(), "has been restored") {
		t.Errorf("restoreSoftDeletedObject output mismatch: got %q", buf.String())
	}

	// Verify the object is restored by checking its attributes.
	restoredAttrs, err := obj.Attrs(ctx)
	if err != nil {
		t.Fatalf("Object(%q).Attrs after restore: %v", objectName, err)
	}
	if !restoredAttrs.Deleted.IsZero() {
		t.Errorf("Object(%q) is still marked as deleted after restore", objectName)
	}
}

func TestListSoftDeletedVersionsOfObject(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	var (
		bucketName  = testutil.UniqueBucketName(testPrefix)
		objectName1 = "soft-deleted-object.txt"
		objectName2 = "soft-deleted-object-2.txt"
	)

	bucket := client.Bucket(bucketName)
	if err := bucket.Create(ctx, tc.ProjectID, &storage.BucketAttrs{SoftDeletePolicy: &storage.SoftDeletePolicy{
		RetentionDuration: 10 * 24 * time.Hour, // 10 days in hours
	}}); err != nil {
		t.Fatalf("Bucket.Create(%q): %v", bucketName, err)
	}
	defer testutil.DeleteBucketIfExists(ctx, client, bucketName)

	// Upload both objects to the bucket.
	if err := uploadFile(io.Discard, bucketName, objectName1); err != nil {
		t.Fatalf("uploadFile(%q): %v", objectName1, err)
	}
	if err := uploadFile(io.Discard, bucketName, objectName2); err != nil {
		t.Fatalf("uploadFile(%q): %v", objectName2, err)
	}

	// Get object attributes of object 1 to retrieve the generation.
	obj1 := client.Bucket(bucketName).Object(objectName1)
	attrs, err := obj1.Attrs(ctx)
	if err != nil {
		t.Fatalf("Object(%q).Attrs: %v", objectName1, err)
	}
	generation := attrs.Generation
	// Simulate soft deletion by deleting both objects.
	if err := obj1.Delete(ctx); err != nil {
		t.Fatalf("Object(%q).Delete: %v", objectName1, err)
	}
	if err := client.Bucket(bucketName).Object(objectName2).Delete(ctx); err != nil {
		t.Fatalf("Object(%q).Delete: %v", objectName2, err)
	}

	var buf bytes.Buffer
	if err := listSoftDeletedVersionsOfObject(&buf, bucketName, objectName1); err != nil {
		t.Fatalf("listSoftDeletedVersionsOfObject: %v", err)
	}
	// Verify the output was printed as expected-- only objectName1 should be listed.
	got := buf.String()
	want := fmt.Sprintf("Soft-deleted object version: %s (generation: %d)\n", objectName1, generation)
	if !strings.HasPrefix(got, want) {
		t.Errorf("Output mismatch: got %q, want %q", got, want)
	}
}

func TestObjectContexts(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	bucketName := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, testPrefix)
	objectName := "context-object.txt"

	// Set new contexts on object.
	if err := uploadWithObjectContexts(io.Discard, bucketName, objectName); err != nil {
		t.Fatalf("setObjectContexts: %v", err)
	}

	var getBuf bytes.Buffer
	if err := getObjectContexts(&getBuf, bucketName, objectName); err != nil {
		t.Fatalf("getObjectContexts: %v", err)
	}
	// Check for new object contexts.
	got := getBuf.String()
	wantGet1 := "key1 = value1"
	if !strings.Contains(got, wantGet1) {
		t.Errorf("getObjectContexts() got %q; want to contain %q", got, wantGet1)
	}
	wantGet2 := "key2 = value2"
	if !strings.Contains(got, wantGet2) {
		t.Errorf("getObjectContexts() got %q; want to contain %q", got, wantGet2)
	}

	// Patch contexts on existing object.
	var patchBuf bytes.Buffer
	if err := setObjectContexts(&patchBuf, bucketName, objectName); err != nil {
		t.Fatalf("setObjectContexts: %v", err)
	}
	gotPatch := patchBuf.String()
	wantGet1 = "key1 = newValue1"
	if !strings.Contains(gotPatch, wantGet1) {
		t.Errorf("setObjectContexts() got %q; want to contain %q", gotPatch, wantGet1)
	}
	wantGet2 = "key3 = value3"
	if !strings.Contains(gotPatch, wantGet2) {
		t.Errorf("setObjectContexts() got %q; want to contain %q", gotPatch, wantGet2)
	}
	// Object should not contain deleted key.
	absentKey := "key2"
	if strings.Contains(gotPatch, absentKey) {
		t.Errorf("setObjectContexts() got %q; should not contain %q", gotPatch, absentKey)
	}

	var listBuf bytes.Buffer
	filter := "contexts.\"key1\"=\"newValue1\""
	if err := listObjectContexts(&listBuf, bucketName, filter); err != nil {
		t.Fatalf("listObjectContexts: %v", err)
	}
	gotList := listBuf.String()
	if !strings.Contains(gotList, objectName) {
		t.Errorf("listObjectContexts() got %q; want to contain %q", gotList, objectName)
	}

	// Delete all contexts of an object.
	if err := deleteObjectContexts(io.Discard, bucketName, objectName); err != nil {
		t.Fatalf("setObjectContexts: %v", err)
	}

	var getBufAfterDelete bytes.Buffer
	if err := getObjectContexts(&getBufAfterDelete, bucketName, objectName); err != nil {
		t.Fatalf("getObjectContexts: %v", err)
	}
	gotAfterDelete := getBufAfterDelete.String()
	wantAfterDelete := fmt.Sprintf("No contexts found for %v", objectName)
	if !strings.Contains(gotAfterDelete, wantAfterDelete) {
		t.Errorf("getObjectContexts() got %q; want %q", gotAfterDelete, wantAfterDelete)
	}
}
