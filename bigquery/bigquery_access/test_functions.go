package bigqueryaccess

import (
	"context"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestClient(t *testing.T, ctx context.Context) (*bigquery.Client, error) {
	tc := testutil.SystemTest(t)
	t.Helper()

	// Creates a client.
	client, err := bigquery.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	return client, nil
}

func TestCleanup(t *testing.T, ctx context.Context, client *bigquery.Client, datasetName string) {
	t.Helper()

	if err := client.Dataset(datasetName).DeleteWithContents(ctx); err != nil {
		t.Errorf("Failed to delete table: %v", err)
	}

	if err := client.Close(); err != nil {
		t.Fatalf("Failed to close Big Query client: %v", err)
	}
}
