// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// The alert sample demonstrates interacting with the monitoring
// and alerting API.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/api/iterator"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
	fieldmask "google.golang.org/genproto/protobuf/field_mask"
)

func main() {
	projectID := flag.String("projectID", "", "Project ID")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage:
alert -projectID <project ID> <action> <args>

Actions:
* list: print alert policies to stdout.
* backup: serialize policies and channels to stdout.
* restore: restore policies from stdin (based on backout output).
* replaceChannels <alertPolicyID> <comma,separated,channelIDs>: replace project channels.
* enable: enable all policies.
* disable: disable all policies.
`)
	}
	flag.Parse()
	if *projectID == "" {
		fmt.Fprintf(os.Stderr, "projectID must be set\n")
		os.Exit(1)
	}
	switch flag.Arg(0) {
	case "list":
		if err := listAlertPolicies(os.Stdout, *projectID); err != nil {
			log.Fatal(err)
		}
	case "backup":
		if err := backupPolicies(os.Stdout, *projectID); err != nil {
			log.Fatal(err)
		}
	case "restore":
		if err := restorePolicies(os.Stdout, *projectID, os.Stdin); err != nil {
			log.Fatal(err)
		}
	case "replaceChannels":
		if err := replaceChannels(os.Stdout, *projectID, flag.Arg(1), strings.Split(flag.Arg(2), ",")); err != nil {
			log.Fatal(err)
		}
	case "enable":
		if err := enablePolicies(os.Stdout, *projectID, true); err != nil {
			log.Fatal(err)
		}
	case "disable":
		if err := enablePolicies(os.Stdout, *projectID, false); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid action: %q\n", flag.Arg(0))
		flag.Usage()
		os.Exit(1)
	}
}

// [START monitoring_alert_list_policies]

// listAlertPolicies lists the alert policies in the project.
func listAlertPolicies(w io.Writer, projectID string) error {
	ctx := context.Background()
	client, err := monitoring.NewAlertPolicyClient(ctx)
	if err != nil {
		return err
	}

	req := &monitoringpb.ListAlertPoliciesRequest{
		Name: "projects/" + projectID,
		// Filter:  "", // See https://cloud.google.com/monitoring/api/v3/sorting-and-filtering.
		// OrderBy: "", // See https://cloud.google.com/monitoring/api/v3/sorting-and-filtering.
	}
	it := client.ListAlertPolicies(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			fmt.Fprintln(w, "Done")
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "  Name: %q\n", resp.GetName())
		fmt.Fprintf(w, "  Display Name: %q\n", resp.GetDisplayName())
		fmt.Fprintf(w, "  Documentation Content: %q\n\n", resp.GetDocumentation().GetContent())
	}
	return nil
}

// [END monitoring_alert_list_policies]

// [START monitoring_alert_backup_policies]

// backupPolicies writes a JSON representation of the project's alert
// policies and notification channels.
func backupPolicies(w io.Writer, projectID string) error {
	b := backup{ProjectID: projectID}
	ctx := context.Background()

	alertClient, err := monitoring.NewAlertPolicyClient(ctx)
	if err != nil {
		return err
	}
	alertReq := &monitoringpb.ListAlertPoliciesRequest{
		Name: "projects/" + projectID,
		// Filter:  "", // See https://cloud.google.com/monitoring/api/v3/sorting-and-filtering.
		// OrderBy: "", // See https://cloud.google.com/monitoring/api/v3/sorting-and-filtering.
	}
	alertIt := alertClient.ListAlertPolicies(ctx, alertReq)
	for {
		resp, err := alertIt.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		b.AlertPolicies = append(b.AlertPolicies, &alertPolicy{resp})
	}

	// [START monitoring_alert_list_channels]
	channelClient, err := monitoring.NewNotificationChannelClient(ctx)
	if err != nil {
		return err
	}
	channelReq := &monitoringpb.ListNotificationChannelsRequest{
		Name: "projects/" + projectID,
		// Filter:  "", // See https://cloud.google.com/monitoring/api/v3/sorting-and-filtering.
		// OrderBy: "", // See https://cloud.google.com/monitoring/api/v3/sorting-and-filtering.
	}
	channelIt := channelClient.ListNotificationChannels(ctx, channelReq)
	for {
		resp, err := channelIt.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		b.Channels = append(b.Channels, &channel{resp})
	}
	// [END monitoring_alert_list_channels]
	bs, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	if _, err := w.Write(bs); err != nil {
		return err
	}
	return nil
}

// alertPolicy is a wrapper around the AlertPolicy proto to
// ensure JSON marshaling/unmarshaling works correctly.
type alertPolicy struct {
	*monitoringpb.AlertPolicy
}

// channel is a wrapper around the NotificationChannel proto to
// ensure JSON marshaling/unmarshaling works correctly.
type channel struct {
	*monitoringpb.NotificationChannel
}

// backup is used to backup and restore a project's policies.
type backup struct {
	ProjectID     string
	AlertPolicies []*alertPolicy
	Channels      []*channel
}

func (a *alertPolicy) MarshalJSON() ([]byte, error) {
	m := &jsonpb.Marshaler{EmitDefaults: true}
	b := new(bytes.Buffer)
	m.Marshal(b, a.AlertPolicy)
	return b.Bytes(), nil
}

func (a *alertPolicy) UnmarshalJSON(b []byte) error {
	u := &jsonpb.Unmarshaler{}
	a.AlertPolicy = new(monitoringpb.AlertPolicy)
	return u.Unmarshal(bytes.NewReader(b), a.AlertPolicy)
}

func (c *channel) MarshalJSON() ([]byte, error) {
	m := &jsonpb.Marshaler{}
	b := new(bytes.Buffer)
	m.Marshal(b, c.NotificationChannel)
	return b.Bytes(), nil
}

func (c *channel) UnmarshalJSON(b []byte) error {
	u := &jsonpb.Unmarshaler{}
	c.NotificationChannel = new(monitoringpb.NotificationChannel)
	return u.Unmarshal(bytes.NewReader(b), c.NotificationChannel)
}

// [END monitoring_alert_backup_policies]

// [START monitoring_alert_restore_policies]
// [START monitoring_alert_create_policy]
// [START monitoring_alert_create_channel]
// [START monitoring_alert_update_channel]

// restorePolicies updates the project with the alert policies and
// notification channels in r.
func restorePolicies(w io.Writer, projectID string, r io.Reader) error {
	b := backup{}
	if err := json.NewDecoder(r).Decode(&b); err != nil {
		return err
	}
	sameProject := projectID == b.ProjectID

	ctx := context.Background()

	alertClient, err := monitoring.NewAlertPolicyClient(ctx)
	if err != nil {
		return err
	}
	channelClient, err := monitoring.NewNotificationChannelClient(ctx)
	if err != nil {
		return err
	}

	// When a channel is recreated, rather than updated, it will get
	// a new name.  We have to update the AlertPolicy with the new
	// name.  channelNames keeps track of the new names.
	channelNames := make(map[string]string)
	for _, c := range b.Channels {
		fmt.Fprintf(w, "Updating channel %q\n", c.GetDisplayName())
		c.VerificationStatus = monitoringpb.NotificationChannel_VERIFICATION_STATUS_UNSPECIFIED
		updated := false
		if sameProject {
			req := &monitoringpb.UpdateNotificationChannelRequest{
				NotificationChannel: c.NotificationChannel,
			}
			_, err := channelClient.UpdateNotificationChannel(ctx, req)
			if err == nil {
				updated = true
			}
		}
		if !updated {
			req := &monitoringpb.CreateNotificationChannelRequest{
				Name:                "projects/" + projectID,
				NotificationChannel: c.NotificationChannel,
			}
			oldName := c.GetName()
			c.Name = ""
			newC, err := channelClient.CreateNotificationChannel(ctx, req)
			if err != nil {
				return err
			}
			channelNames[oldName] = newC.GetName()
		}
	}

	for _, policy := range b.AlertPolicies {
		fmt.Fprintf(w, "Updating alert %q\n", policy.GetDisplayName())
		policy.CreationRecord = nil
		policy.MutationRecord = nil
		for i, aChannel := range policy.GetNotificationChannels() {
			if c, ok := channelNames[aChannel]; ok {
				policy.NotificationChannels[i] = c
			}
		}
		updated := false
		if sameProject {
			req := &monitoringpb.UpdateAlertPolicyRequest{
				AlertPolicy: policy.AlertPolicy,
			}
			_, err := alertClient.UpdateAlertPolicy(ctx, req)
			if err == nil {
				updated = true
			}
		}
		if !updated {
			req := &monitoringpb.CreateAlertPolicyRequest{
				Name:        "projects/" + projectID,
				AlertPolicy: policy.AlertPolicy,
			}
			if _, err = alertClient.CreateAlertPolicy(ctx, req); err != nil {
				log.Fatal(err)
			}
		}
	}
	fmt.Fprintf(w, "Successfully restored alerts.")
	return nil
}

// [END monitoring_alert_restore_policies]
// [END monitoring_alert_create_policy]
// [END monitoring_alert_create_channel]
// [END monitoring_alert_update_channel]

// [START monitoring_alert_replace_channels]

// replaceChannels replaces the notification channels in the alert policy
// with channelIDs.
func replaceChannels(w io.Writer, projectID, alertPolicyID string, channelIDs []string) error {
	ctx := context.Background()

	client, err := monitoring.NewAlertPolicyClient(ctx)
	if err != nil {
		return err
	}

	policy := &monitoringpb.AlertPolicy{
		Name: "projects/" + projectID + "/alertPolicies/" + alertPolicyID,
	}
	for _, c := range channelIDs {
		c = "projects/" + projectID + "/notificationChannels/" + c
		policy.NotificationChannels = append(policy.NotificationChannels, c)
	}
	req := &monitoringpb.UpdateAlertPolicyRequest{
		AlertPolicy: policy,
		UpdateMask: &fieldmask.FieldMask{
			Paths: []string{"notification_channels"},
		},
	}
	if _, err := client.UpdateAlertPolicy(ctx, req); err != nil {
		return fmt.Errorf("UpdateAlertPolicy: %v", err)
	}
	fmt.Fprintf(w, "Successfully replaced channels.")
	return nil
}

// [END monitoring_alert_replace_channels]

// [START monitoring_alert_disable_policies]
// [START monitoring_alert_enable_policies]

// enablePolicies enables or disables all alert policies in the project.
func enablePolicies(w io.Writer, projectID string, enable bool) error {
	ctx := context.Background()

	client, err := monitoring.NewAlertPolicyClient(ctx)
	if err != nil {
		return err
	}

	req := &monitoringpb.ListAlertPoliciesRequest{
		Name: "projects/" + projectID,
		// Filter:  "", // See https://cloud.google.com/monitoring/api/v3/sorting-and-filtering.
		// OrderBy: "", // See https://cloud.google.com/monitoring/api/v3/sorting-and-filtering.
	}
	it := client.ListAlertPolicies(ctx, req)
	for {
		a, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		if a.GetEnabled().GetValue() == enable {
			fmt.Fprintf(w, "Policy %q already has enabled=%v", a.GetDisplayName(), enable)
			continue
		}
		a.Enabled = &wrappers.BoolValue{Value: enable}
		req := &monitoringpb.UpdateAlertPolicyRequest{
			AlertPolicy: a,
			UpdateMask: &fieldmask.FieldMask{
				Paths: []string{"enabled"},
			},
		}
		if _, err := client.UpdateAlertPolicy(ctx, req); err != nil {
			return err
		}
	}
	fmt.Fprintln(w, "Successfully updated alerts.")
	return nil
}

// [END monitoring_alert_enable_policies]
// [END monitoring_alert_disable_policies]
