package inspect

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	dataSetID = "dlp_test_dataset"
	tableID   = "dlp_inspect_test_table_table_id"
)

func TestInspectBigQuerySendToScc(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer

	if err := inspectBigQuerySendToScc(&buf, tc.ProjectID, dataSetID, tableID); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	if want := "Job created successfully:"; !strings.Contains(got, want) {
		t.Errorf("InspectBigQuerySendToScc got %q, want %q", got, want)
	}

	jobName := strings.SplitAfter(got, "Job created successfully: ")
	log.Printf("Job Name : %v", jobName)

	deleteJob(tc.ProjectID, jobName[1])
}
