// Package uptime contains samples of uptime checks.
package uptime

import (
	"fmt"
	"io"
	"log"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/golang/protobuf/ptypes/duration"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// [START monitoring_uptime_check_create]

// create creates an example uptime check.
func create(w io.Writer, projectID string) {
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
	fmt.Fprintf(w, "Successfully created uptime check %q\n", uc.GetDisplayName())
}

// [END monitoring_uptime_check_create]

// [START monitoring_uptime_check_list_configs]

// list is an example of listing the uptime checks in projectID.
func list(w io.Writer, projectID string) {
	ctx := context.Background()
	client, err := monitoring.NewUptimeCheckClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	req := &monitoringpb.ListUptimeCheckConfigsRequest{
		Parent: "projects/" + projectID,
	}
	it := client.ListUptimeCheckConfigs(ctx, req)
	for {
		config, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error getting uptime checks: %v", err)
		}
		fmt.Fprintln(w, config)
	}
	fmt.Fprintln(w, "Done listing uptime checks")
}

// [END monitoring_uptime_check_list_configs]

// [START monitoring_uptime_check_list_ips]

// listIPs is an example of listing uptime check IPs.
func listIPs(w io.Writer) {
	ctx := context.Background()
	client, err := monitoring.NewUptimeCheckClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	req := &monitoringpb.ListUptimeCheckIpsRequest{}
	it := client.ListUptimeCheckIps(ctx, req)
	for {
		config, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error getting uptime checks: %v", err)
		}
		fmt.Fprintln(w, config)
	}
	fmt.Fprintln(w, "Done listing uptime check IPs")
}

// [END monitoring_uptime_check_list_ips]

// [START monitoring_uptime_check_get]

// get is an example of getting an uptime check. resourceName should be
// of the form `projects/[PROJECT_ID]/uptimeCheckConfigs/[UPTIME_CHECK_ID]`.
func get(w io.Writer, resourceName string) {
	ctx := context.Background()
	client, err := monitoring.NewUptimeCheckClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	req := &monitoringpb.GetUptimeCheckConfigRequest{
		Name: resourceName,
	}
	config, err := client.GetUptimeCheckConfig(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Config: %v", config)
}

// [END monitoring_uptime_check_get]

// [START monitoring_uptime_check_delete]

// delete is an example of deleting an uptime check. resourceName should be
// of the form `projects/[PROJECT_ID]/uptimeCheckConfigs/[UPTIME_CHECK_ID]`.
func delete(w io.Writer, resourceName string) {
	ctx := context.Background()
	client, err := monitoring.NewUptimeCheckClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	req := &monitoringpb.DeleteUptimeCheckConfigRequest{
		Name: resourceName,
	}
	if err := client.DeleteUptimeCheckConfig(ctx, req); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Successfully deleted %q", resourceName)
}

// [END monitoring_uptime_check_delete]
