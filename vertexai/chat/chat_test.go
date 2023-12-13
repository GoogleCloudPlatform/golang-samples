package chat

import (
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func Test_makeChatRequests(t *testing.T) {
	tc := testutil.SystemTest(t)
	err := makeChatRequests(tc.ProjectID, "us-central1", "gemini-pro-vision")
	if err != nil {
		t.Errorf("unexpected error: %v", err.Error())
	}
}
