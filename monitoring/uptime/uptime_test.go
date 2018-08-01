package uptime

import (
	"bytes"
	"log"
	"strings"
	"testing"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/golang/protobuf/ptypes/duration"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

func TestCreate(t *testing.T) {
	c := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	create(buf, c.ProjectID)
	want := "Successfully"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("%q not found in output: %q", want, got)
	}
}

func TestList(t *testing.T) {
	c := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	list(buf, c.ProjectID)
	want := "Done"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("%q not found in output: %q", want, got)
	}
}

func TestListIPs(t *testing.T) {
	testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	listIPs(buf)
	want := "Done"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("%q not found in output: %q", want, got)
	}
}

func createTestUptimeCheck(projectID string) *monitoringpb.UptimeCheckConfig {
	ctx := context.Background()
	client, err := monitoring.NewUptimeCheckClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	req := &monitoringpb.CreateUptimeCheckConfigRequest{
		Parent: "projects/" + projectID,
		UptimeCheckConfig: &monitoringpb.UptimeCheckConfig{
			DisplayName: "new uptime check",
			Resource: &monitoringpb.UptimeCheckConfig_MonitoredResource{
				MonitoredResource: &monitoredres.MonitoredResource{
					Type: "uptime_url",
					Labels: map[string]string{
						"host": "example.com",
					},
				},
			},
			CheckRequestType: &monitoringpb.UptimeCheckConfig_HttpCheck_{
				HttpCheck: &monitoringpb.UptimeCheckConfig_HttpCheck{
					Path: "/",
					Port: 80,
				},
			},
			Timeout: &duration.Duration{Seconds: 10},
			Period:  &duration.Duration{Seconds: 300},
		},
	}
	uc, err := client.CreateUptimeCheckConfig(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	return uc
}

func TestGet(t *testing.T) {
	c := testutil.SystemTest(t)
	uc := createTestUptimeCheck(c.ProjectID)
	buf := new(bytes.Buffer)
	get(buf, uc.GetName())
	want := "Config:"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("%q not found in output: %q", want, got)
	}
}

func TestDelete(t *testing.T) {
	c := testutil.SystemTest(t)
	uc := createTestUptimeCheck(c.ProjectID)
	buf := new(bytes.Buffer)
	delete(buf, uc.GetName())
	want := "Successfully"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("%q not found in output: %q", want, got)
	}
}
