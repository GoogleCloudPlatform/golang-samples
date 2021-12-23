package retry

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestConfigureRetries(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket, err := testutil.CreateTestBucket(ctx, t, client, tc.ProjectID, "storage-buckets-test")
	if err != nil {
		t.Fatalf("creating bucket: %v", err)
	}
	defer testutil.DeleteBucketIfExists(ctx, client, bucket)
	object := "foo.txt"

	// Upload test object to delete in sample.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, strings.NewReader("hello world")); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		t.Fatalf("Writer.Close: %v", err)
	}

	var buf bytes.Buffer
	if err := configureRetries(&buf, bucket, object); err != nil {
		t.Errorf("configureRetries: %v", err)
	}

	if got, want := buf.String(), "deleted with a customized retry"; !strings.Contains(got, want) {
		t.Errorf("configureRetries: got %q; want to contain %q", got, want)
	}

}
