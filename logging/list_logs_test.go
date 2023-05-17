package snippets

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const logID = "test-list-logs"

var projectID string

func TestMain(m *testing.M) {
	ctx := context.Background()

	tc, ok := testutil.ContextMain(m)
	if !ok {
		log.Fatal("test project not set up properly")
		return
	}
	projectID = tc.ProjectID

	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("logging.NewClient(%q) failed: %v", projectID, err)
	}
	defer client.Close()

	adminClient, err := logadmin.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("logadmin.NewClient(%q) failed: %v", projectID, err)
	}
	defer adminClient.Close()

	// Create a log
	logger := client.Logger(logID)
	logger.Log(logging.Entry{Payload: "create a log"})
	if err := logger.Flush(); err != nil {
		log.Fatalf("logger.Flush() failed: %v", err)
	}
	// Delete the log
	defer adminClient.DeleteLog(ctx, logID)

	m.Run()
}

func TestListLogs(t *testing.T) {
	testutil.Retry(t, 6, 10*time.Second, func(r *testutil.R) {
		buf := &bytes.Buffer{}
		if err := listLogs(buf, projectID); err != nil {
			r.Errorf("listLogs(%q) failed: %v", projectID, err)
			return
		}
		if got, want := buf.String(), logID; !strings.Contains(got, want) {
			r.Errorf("listLogs got %q, want to contain %q", got, want)
		}
	})
}
