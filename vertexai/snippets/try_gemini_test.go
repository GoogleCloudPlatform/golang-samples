package snippets

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func Test_tryGemini(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	err := tryGemini(buf, tc.ProjectID, "us-central1", "gemini-pro-vision")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	fmt.Println(buf)
}
