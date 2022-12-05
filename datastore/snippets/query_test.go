package datastore_snippets

import (
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestNotEqualQuery(t *testing.T) {
	tc := testutil.SystemTest(t)
	err := queryNotEquals(tc.ProjectID)
	if err != nil {
		t.Fatal(err)
	}
}
