package snippets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestGetSink(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	sinkName := "_Default"
	if err := getSink(buf, tc.ProjectID, sinkName); err != nil {
		t.Fatalf("getSink(%q, %q) failed: %v", tc.ProjectID, sinkName, err)
	}
	if !strings.Contains(buf.String(), sinkName) {
		t.Errorf("getSink got %q, want to contain %q", buf.String(), sinkName)
	}
}
