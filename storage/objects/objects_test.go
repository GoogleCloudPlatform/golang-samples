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
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestObjects runs all samples tests of the package.
func TestObjects(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	var (
		bucket                = tc.ProjectID + "-samples-object-bucket-1"
		dstBucket             = tc.ProjectID + "-samples-object-bucket-2"
		object1               = "foo.txt"
		object2               = "foo/a.txt"
		allAuthenticatedUsers = storage.AllAuthenticatedUsers
		roleReader            = storage.RoleReader
	)

	cleanBucket(t, ctx, client, tc.ProjectID, bucket)
	cleanBucket(t, ctx, client, tc.ProjectID, dstBucket)

	if err := uploadFile(ioutil.Discard, bucket, object1); err != nil {
		t.Fatalf("uploadFile(%q): %v", object1, err)
	}
	if err := uploadFile(ioutil.Discard, bucket, object2); err != nil {
		t.Fatalf("uploadFile(%q): %v", object2, err)
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
		if err := downloadUsingRequesterPays(ioutil.Discard, bucket, object1, tc.ProjectID); err != nil {
			t.Errorf("downloadUsingRequesterPays: %v", err)
		}
	}

	data, err := downloadFile(ioutil.Discard, bucket, object1)
	if err != nil {
		t.Fatalf("downloadFile: %v", err)
	}
	if got, want := string(data), "Hello\nworld"; got != want {
		t.Errorf("contents = %q; want %q", got, want)
	}

	_, err = getMetadata(ioutil.Discard, bucket, object1)
	if err != nil {
		t.Errorf("getMetadata: %v", err)
	}
	if err := makePublic(ioutil.Discard, bucket, object1, allAuthenticatedUsers, roleReader); err != nil {
		t.Errorf("makePublic: %v", err)
	}

	err = moveFile(ioutil.Discard, bucket, object1)
	if err != nil {
		t.Fatalf("moveFile: %v", err)
	}
	// object1's new name.
	object1 = object1 + "-rename"

	if err := copyFile(ioutil.Discard, dstBucket, bucket, object1); err != nil {
		t.Errorf("copyFile: %v", err)
	}

	key := []byte("my-secret-AES-256-encryption-key")
	newKey := []byte("My-secret-AES-256-encryption-key")

	if err := uploadEncryptedFile(ioutil.Discard, bucket, object1, key); err != nil {
		t.Errorf("uploadEncryptedFile: %v", err)
	}
	data, err = downloadEncryptedFile(ioutil.Discard, bucket, object1, key)
	if err != nil {
		t.Errorf("downloadEncryptedFile: %v", err)
	}
	if got, want := string(data), "top secret"; got != want {
		t.Errorf("object content = %q; want %q", got, want)
	}
	if err := rotateEncryptionKey(ioutil.Discard, bucket, object1, key, newKey); err != nil {
		t.Errorf("rotateEncryptionKey: %v", err)
	}
	if err := deleteFile(ioutil.Discard, bucket, object1); err != nil {
		t.Errorf("deleteFile: %v", err)
	}
	if err := deleteFile(ioutil.Discard, bucket, object2); err != nil {
		t.Errorf("deleteFile: %v", err)
	}

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		// Cleanup, this part won't be executed if Fatal happens.
		// TODO(jbd): Implement garbage cleaning.
		if err := client.Bucket(bucket).Delete(ctx); err != nil {
			r.Errorf("Bucket(%q).Delete: %v", bucket, err)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		if err := deleteFile(ioutil.Discard, dstBucket, object1+"-copy"); err != nil {
			r.Errorf("deleteFile: %v", err)
		}
	})

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		if err := client.Bucket(dstBucket).Delete(ctx); err != nil {
			r.Errorf("Bucket(%q).Delete: %v", dstBucket, err)
		}
	})
}

func TestKMSObjects(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	keyRingID := os.Getenv("GOLANG_SAMPLES_KMS_KEYRING")
	cryptoKeyID := os.Getenv("GOLANG_SAMPLES_KMS_CRYPTOKEY")
	if keyRingID == "" || cryptoKeyID == "" {
		t.Skip("GOLANG_SAMPLES_KMS_KEYRING and GOLANG_SAMPLES_KMS_CRYPTOKEY must be set")
	}

	var (
		bucket    = tc.ProjectID + "-samples-object-bucket-1"
		dstBucket = tc.ProjectID + "-samples-object-bucket-2"
		object1   = "foo.txt"
	)

	cleanBucket(t, ctx, client, tc.ProjectID, bucket)
	cleanBucket(t, ctx, client, tc.ProjectID, dstBucket)

	kmsKeyName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", tc.ProjectID, "global", keyRingID, cryptoKeyID)

	if err := uploadWithKMSKey(ioutil.Discard, bucket, object1, kmsKeyName); err != nil {
		t.Errorf("uploadWithKMSKey: %v", err)
	}
}

func TestV4SignedURL(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucketName := tc.ProjectID + "-signed-url-bucket-name"
	objectName := "foo.txt"
	serviceAccount := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	cleanBucket(t, ctx, client, tc.ProjectID, bucketName)
	putBuf := new(bytes.Buffer)
	putURL, err := generateV4PutObjectSignedURL(putBuf, bucketName, objectName, serviceAccount)
	if err != nil {
		t.Errorf("generateV4PutObjectSignedURL: %v", err)
	}
	got := putBuf.String()
	if want := "Generated PUT signed URL:"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}

	httpClient := &http.Client{}
	request, err := http.NewRequest("PUT", putURL, strings.NewReader("hello world"))
	request.ContentLength = 11
	request.Header.Set("Content-Type", "application/octet-stream")
	response, err := httpClient.Do(request)
	if err != nil {
		t.Errorf("httpClient.Do: %v", err)
	}
	getBuf := new(bytes.Buffer)
	getURL, err := generateV4GetObjectSignedURL(getBuf, bucketName, objectName, serviceAccount)
	if err != nil {
		t.Errorf("generateV4GetObjectSignedURL: %v", err)
	}
	got = getBuf.String()
	if want := "Generated GET signed URL:"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}

	response, err = http.Get(getURL)
	if err != nil {
		t.Errorf("http.Get: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll: %v", err)
	}

	if got, want := string(body), "hello world"; got != want {
		t.Errorf("object content = %q; want %q", got, want)
	}
}

func TestPostPolicyV4(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucketName := tc.ProjectID + "-post-policy-bucket-name"
	objectName := "foo.txt"
	serviceAccount := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	if err := testutil.CleanBucket(ctx, t, tc.ProjectID, bucketName); err != nil {
		t.Fatalf("CleanBucket: %v", err)
	}
	putBuf := new(bytes.Buffer)
	policy, err := generateSignedPostPolicyV4(putBuf, bucketName, objectName, serviceAccount)
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
		t.Fatalf("failed to compose HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("client.Do: %v", err)
	}
	if g, w := res.StatusCode, 204; g != w {
		blob, _ := httputil.DumpResponse(res, true)
		t.Errorf("status code in response mismatch: got %d want %d\nBody: %s", g, w, blob)
	}
	io.Copy(ioutil.Discard, res.Body)
	if err := res.Body.Close(); err != nil {
		t.Errorf("Body.Close: %v", err)
	}

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
	defer client.Close()

	var (
		bucketName      = tc.ProjectID + "-retent-samples-object-bucket"
		objectName      = "foo.txt"
		retentionPeriod = 5 * time.Second
	)

	cleanBucket(t, ctx, client, tc.ProjectID, bucketName)
	bucket := client.Bucket(bucketName)

	if err := uploadFile(ioutil.Discard, bucketName, objectName); err != nil {
		t.Fatalf("uploadFile(%q): %v", objectName, err)
	}
	if _, err := bucket.Update(ctx, storage.BucketAttrsToUpdate{
		RetentionPolicy: &storage.RetentionPolicy{
			RetentionPeriod: retentionPeriod,
		},
	}); err != nil {
		t.Errorf("Bucket(%q).Update: %v", bucketName, err)
	}
	if err := setEventBasedHold(ioutil.Discard, bucketName, objectName); err != nil {
		t.Errorf("setEventBasedHold(%q, %q): %v", bucketName, objectName, err)
	}
	oAttrs, err := getMetadata(ioutil.Discard, bucketName, objectName)
	if err != nil {
		t.Errorf("getMetadata: %v", err)
	}
	if !oAttrs.EventBasedHold {
		t.Errorf("event-based hold is not enabled")
	}
	if err := releaseEventBasedHold(ioutil.Discard, bucketName, objectName); err != nil {
		t.Errorf("releaseEventBasedHold(%q, %q): %v", bucketName, objectName, err)
	}
	oAttrs, err = getMetadata(ioutil.Discard, bucketName, objectName)
	if err != nil {
		t.Errorf("getMetadata: %v", err)
	}
	if oAttrs.EventBasedHold {
		t.Errorf("event-based hold is not disabled")
	}
	if _, err := bucket.Update(ctx, storage.BucketAttrsToUpdate{
		RetentionPolicy: &storage.RetentionPolicy{},
	}); err != nil {
		t.Errorf("Bucket(%q).Update: %v", bucketName, err)
	}
	if err := setTemporaryHold(ioutil.Discard, bucketName, objectName); err != nil {
		t.Errorf("setTemporaryHold(%q, %q): %v", bucketName, objectName, err)
	}
	oAttrs, err = getMetadata(ioutil.Discard, bucketName, objectName)
	if err != nil {
		t.Errorf("getMetadata: %v", err)
	}
	if !oAttrs.TemporaryHold {
		t.Errorf("temporary hold is not disabled")
	}
	if err := releaseTemporaryHold(ioutil.Discard, bucketName, objectName); err != nil {
		t.Errorf("releaseTemporaryHold(%q, %q): %v", bucketName, objectName, err)
	}
	oAttrs, err = getMetadata(ioutil.Discard, bucketName, objectName)
	if err != nil {
		t.Errorf("getMetadata: %v", err)
	}
	if oAttrs.TemporaryHold {
		t.Errorf("temporary hold is not disabled")
	}
}

// cleanBucket ensures there's a fresh bucket with a given name, deleting the existing bucket if it already exists.
func cleanBucket(t *testing.T, ctx context.Context, client *storage.Client, projectID, bucket string) {
	b := client.Bucket(bucket)
	_, err := b.Attrs(ctx)
	if err == nil {
		it := b.Objects(ctx, nil)
		for {
			attrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				t.Fatalf("Bucket(%q).Objects: %v", bucket, err)
			}
			if attrs.EventBasedHold || attrs.TemporaryHold {
				if _, err := b.Object(attrs.Name).Update(ctx, storage.ObjectAttrsToUpdate{
					TemporaryHold:  false,
					EventBasedHold: false,
				}); err != nil {
					t.Fatalf("Bucket(%q).Object(%q).Update: %v", bucket, attrs.Name, err)
				}
			}
			if err := b.Object(attrs.Name).Delete(ctx); err != nil {
				t.Fatalf("Bucket(%q).Object(%q).Delete: %v", bucket, attrs.Name, err)
			}
		}
		if err := b.Delete(ctx); err != nil {
			t.Fatalf("Bucket(%q).Delete: %v", bucket, err)
		}
	}
	if err := b.Create(ctx, projectID, nil); err != nil && status.Code(err) != codes.AlreadyExists {
		t.Fatalf("Bucket(%q).Create: %v", bucket, err)
	}
}
