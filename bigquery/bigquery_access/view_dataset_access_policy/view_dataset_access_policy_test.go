package viewdatasetaccesspolicy

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestViewDatasetAccessPolicies(t *testing.T) {
	tc := testutil.SystemTest(t)

	datasetName := "my_new_dataset"

	b := bytes.Buffer{}

	if err := viewDatasetAccessPolicies(&b, tc.ProjectID, datasetName); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Role"; !strings.Contains(got, want) {
		t.Errorf("viewDatasetAccessPolicies: expected %q to contain %q", got, want)
	}

}
