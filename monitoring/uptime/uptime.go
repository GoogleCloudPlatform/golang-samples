// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package uptime contains samples of uptime checks.
package uptime

import (
	"context"
	"fmt"
	"io"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/golang/protobuf/ptypes/duration"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
	"google.golang.org/genproto/protobuf/field_mask"
)

// [START monitoring_uptime_check_create]

// create creates an example uptime check.
func create(w io.Writer, projectID string) (*monitoringpb.UptimeCheckConfig, error) {
	ctx := context.Background()
	client, err := monitoring.NewUptimeCheckClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewUptimeCheckClient: %v", err)
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
	config, err := client.CreateUptimeCheckConfig(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("CreateUptimeCheckConfig: %v", err)
	}
	fmt.Fprintf(w, "Successfully created uptime check %q\n", config.GetDisplayName())
	return config, nil
}

// [END monitoring_uptime_check_create]

// [START monitoring_uptime_check_list_configs]

// list is an example of listing the uptime checks in projectID.
func list(w io.Writer, projectID string) error {
	ctx := context.Background()
	client, err := monitoring.NewUptimeCheckClient(ctx)
	if err != nil {
		return fmt.Errorf("NewUptimeCheckClient: %v", err)
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
			return fmt.Errorf("ListUptimeCheckConfigs: %v", err)
		}
		fmt.Fprintln(w, config)
	}
	fmt.Fprintln(w, "Done listing uptime checks")
	return nil
}

// [END monitoring_uptime_check_list_configs]

// [START monitoring_uptime_check_list_ips]

// listIPs is an example of listing uptime check IPs.
func listIPs(w io.Writer) error {
	ctx := context.Background()
	client, err := monitoring.NewUptimeCheckClient(ctx)
	if err != nil {
		return fmt.Errorf("NewUptimeCheckClient: %v", err)
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
			return fmt.Errorf("ListUptimeCheckIps: %v", err)
		}
		fmt.Fprintln(w, config)
	}
	fmt.Fprintln(w, "Done listing uptime check IPs")
	return nil
}

// [END monitoring_uptime_check_list_ips]

// [START monitoring_uptime_check_get]

// get is an example of getting an uptime check. resourceName should be
// of the form `projects/[PROJECT_ID]/uptimeCheckConfigs/[UPTIME_CHECK_ID]`.
func get(w io.Writer, resourceName string) (*monitoringpb.UptimeCheckConfig, error) {
	ctx := context.Background()
	client, err := monitoring.NewUptimeCheckClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewUptimeCheckClient: %v", err)
	}
	defer client.Close()
	req := &monitoringpb.GetUptimeCheckConfigRequest{
		Name: resourceName,
	}
	config, err := client.GetUptimeCheckConfig(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetUptimeCheckConfig: %v", err)
	}
	fmt.Fprintf(w, "Config: %v", config)
	return config, nil
}

// [END monitoring_uptime_check_get]

// [START monitoring_uptime_check_update]

// update is an example of updating an uptime check. resourceName should be
// of the form `projects/[PROJECT_ID]/uptimeCheckConfigs/[UPTIME_CHECK_ID]`.
func update(w io.Writer, resourceName, displayName, httpCheckPath string) (*monitoringpb.UptimeCheckConfig, error) {
	ctx := context.Background()
	client, err := monitoring.NewUptimeCheckClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewUptimeCheckClient: %v", err)
	}
	defer client.Close()
	getReq := &monitoringpb.GetUptimeCheckConfigRequest{
		Name: resourceName,
	}
	config, err := client.GetUptimeCheckConfig(ctx, getReq)
	if err != nil {
		return nil, fmt.Errorf("GetUptimeCheckConfig: %v", err)
	}
	config.DisplayName = displayName
	config.GetHttpCheck().Path = httpCheckPath
	req := &monitoringpb.UpdateUptimeCheckConfigRequest{
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"display_name", "http_check.path"},
		},
		UptimeCheckConfig: config,
	}
	config, err = client.UpdateUptimeCheckConfig(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("UpdateUptimeCheckConfig: %v", err)
	}
	fmt.Fprintf(w, "Successfully updated %v", resourceName)
	return config, nil
}

// [END monitoring_uptime_check_update]

// [START monitoring_uptime_check_delete]

// delete is an example of deleting an uptime check. resourceName should be
// of the form `projects/[PROJECT_ID]/uptimeCheckConfigs/[UPTIME_CHECK_ID]`.
func delete(w io.Writer, resourceName string) error {
	ctx := context.Background()
	client, err := monitoring.NewUptimeCheckClient(ctx)
	if err != nil {
		return fmt.Errorf("NewUptimeCheckClient: %v", err)
	}
	defer client.Close()
	req := &monitoringpb.DeleteUptimeCheckConfigRequest{
		Name: resourceName,
	}
	if err := client.DeleteUptimeCheckConfig(ctx, req); err != nil {
		return fmt.Errorf("DeleteUptimeCheckConfig: %v", err)
	}
	fmt.Fprintf(w, "Successfully deleted %q", resourceName)
	return nil
}

// [END monitoring_uptime_check_delete]
