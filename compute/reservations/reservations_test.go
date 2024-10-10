// Copyright 2024 Google LLC
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

package snippets

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/protobuf/proto"
)

func createTemplate(project, templateName string) error {
	ctx := context.Background()
	client, err := compute.NewInstanceTemplatesRESTClient(ctx)
	if err != nil {
		return err
	}

	disk := &computepb.AttachedDisk{
		AutoDelete: proto.Bool(true),
		Boot:       proto.Bool(true),
		InitializeParams: &computepb.AttachedDiskInitializeParams{
			SourceImage: proto.String("projects/debian-cloud/global/images/family/debian-12"),
			DiskSizeGb:  proto.Int64(25),
			DiskType:    proto.String("pd-balanced"),
		},
	}
	req := &computepb.InsertInstanceTemplateRequest{
		Project: project,
		InstanceTemplateResource: &computepb.InstanceTemplate{
			Name: proto.String(templateName),
			Properties: &computepb.InstanceProperties{
				MachineType: proto.String("n1-standard-4"),
				Disks:       []*computepb.AttachedDisk{disk},
				NetworkInterfaces: []*computepb.NetworkInterface{{
					Name: proto.String("global/networks/default"),
				}},
			},
		},
	}
	op, err := client.Insert(ctx, req)
	if err != nil {
		return err
	}
	return op.Wait(ctx)
}

func getTemplate(project, templateName string) (*computepb.InstanceTemplate, error) {
	ctx := context.Background()
	client, err := compute.NewInstanceTemplatesRESTClient(ctx)
	if err != nil {
		return nil, err
	}

	req := &computepb.GetInstanceTemplateRequest{
		Project:          project,
		InstanceTemplate: templateName,
	}
	return client.Get(ctx, req)
}

func deleteTemplate(project, templateName string) error {
	ctx := context.Background()
	client, err := compute.NewInstanceTemplatesRESTClient(ctx)
	if err != nil {
		return err
	}

	req := &computepb.DeleteInstanceTemplateRequest{
		Project:          project,
		InstanceTemplate: templateName,
	}
	op, err := client.Delete(ctx, req)
	if err != nil {
		return err
	}

	return op.Wait(ctx)
}

func deleteInstance(project, zone, instance string) error {
	ctx := context.Background()
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return err
	}

	req := &computepb.DeleteInstanceRequest{
		Instance: instance,
		Project:  project,
		Zone:     zone,
	}
	op, err := client.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete instance: %w", err)
	}

	return op.Wait(ctx)
}

func deleteInstanceTemplate(project, templateName string) error {
	ctx := context.Background()
	instanceTemplatesClient, err := compute.NewInstanceTemplatesRESTClient(ctx)
	if err != nil {
		return err
	}
	defer instanceTemplatesClient.Close()

	req := &computepb.DeleteInstanceTemplateRequest{
		Project:          project,
		InstanceTemplate: templateName,
	}

	op, err := instanceTemplatesClient.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete instance template: %w", err)
	}

	return op.Wait(ctx)
}

func TestReservations(t *testing.T) {
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-west2-b"
	templateName := fmt.Sprintf("test-template-%v-%v", time.Now().Format("01-02-2006"), r.Int())

	var buf bytes.Buffer

	err := createTemplate(tc.ProjectID, templateName)
	if err != nil {
		t.Errorf("createTemplate got err: %v", err)
	}
	defer deleteTemplate(tc.ProjectID, templateName)

	sourceTemplate, err := getTemplate(tc.ProjectID, templateName)
	if err != nil {
		t.Errorf("getTemplate got err: %v", err)
	}

	t.Run("Test basics", func(t *testing.T) {
		reservationName := fmt.Sprintf("test-reservation-%v-%v", time.Now().Format("01-02-2006"), r.Int())

		want := "Reservation created"
		if err := createReservation(&buf, tc.ProjectID, zone, reservationName, *sourceTemplate.SelfLink); err != nil {
			t.Errorf("createReservation got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createReservation got %s, want %s", got, want)
		}
		buf.Reset()

		want = fmt.Sprintf("Reservation: %s", reservationName)
		if _, err := getReservation(&buf, tc.ProjectID, zone, reservationName); err != nil {
			t.Errorf("getReservation got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("getReservation got %s, want %s", got, want)
		}
		buf.Reset()

		want = fmt.Sprintf("- %s %d", reservationName, 2)
		if err := listReservations(&buf, tc.ProjectID, zone); err != nil {
			t.Errorf("listReservations got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("listReservations got %s, want %s", got, want)
		}
		buf.Reset()

		want = "Reservation deleted"
		if err := deleteReservation(&buf, tc.ProjectID, zone, reservationName); err != nil {
			t.Errorf("deleteReservation got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("deleteReservation got %s, want %s", got, want)
		}
	})

	t.Run("Test update VMs", func(t *testing.T) {
		reservationName := fmt.Sprintf("test-reservation-%v-%v", time.Now().Format("01-02-2006"), r.Int())

		if err := createReservation(&buf, tc.ProjectID, zone, reservationName, *sourceTemplate.SelfLink); err != nil {
			t.Fatalf("createReservation got err: %v", err)
		}
		buf.Reset()

		numberOfVMs := int64(5)
		if err := updateReservationVMS(&buf, tc.ProjectID, zone, reservationName, numberOfVMs); err != nil {
			t.Errorf("updateReservationVMS got err: %v", err)
		}

		reservation, err := getReservation(&buf, tc.ProjectID, zone, reservationName)
		if err != nil {
			t.Errorf("getReservation got err: %v", err)
		}
		count := reservation.GetSpecificReservation().GetCount()
		if count != numberOfVMs {
			t.Errorf("reservation wasn't updated got: %d want: %d", count, numberOfVMs)
		}

		if err := deleteReservation(&buf, tc.ProjectID, zone, reservationName); err != nil {
			t.Errorf("deleteReservation got err: %v", err)
		}
	})

	t.Run("Test without template", func(t *testing.T) {
		reservationName := fmt.Sprintf("test-reservation-%v-%v", time.Now().Format("01-02-2006"), r.Int())

		want := "Reservation created"
		if err := createBaseReservation(&buf, tc.ProjectID, zone, reservationName); err != nil {
			t.Fatalf("createBaseReservation got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createBaseReservation got %s, want %s", got, want)
		}
		buf.Reset()

		if err := deleteReservation(&buf, tc.ProjectID, zone, reservationName); err != nil {
			t.Errorf("deleteReservation got err: %v", err)
		}
	})

	t.Run("Create instance without consuming reservation", func(t *testing.T) {
		reservationName := fmt.Sprintf("test-reservation-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		instanceName := fmt.Sprintf("test-instance-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		if err = createBaseReservation(&buf, tc.ProjectID, zone, reservationName); err != nil {
			t.Errorf("createBaseReservation got err: %v", err)
		}

		ctx := context.Background()
		reservationsClient, err := compute.NewReservationsRESTClient(ctx)
		if err != nil {
			t.Errorf("reservationsClient got err: %v", err)
		}
		defer reservationsClient.Close()

		err = createInstanceNotConsumeReservation(&buf, tc.ProjectID, zone, instanceName)
		if err != nil {
			t.Errorf("createInstanceNotConsumeReservation failed: %v", err)
		}

		want := "Instance created"
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createInstanceNotConsumeReservation got %s, want %s", got, want)
		}

		req := &computepb.GetReservationRequest{
			Project:     tc.ProjectID,
			Zone:        zone,
			Reservation: reservationName,
		}

		resp, err := reservationsClient.Get(ctx, req)
		if err != nil {
			t.Errorf("get reservation failed: %v", err)
		}
		inUseAfter := resp.GetSpecificReservation().GetInUseCount()
		if inUseAfter != 0 {
			t.Errorf("reservation was consumed. Expected 0, got %d", inUseAfter)
		}

		if err = deleteInstance(tc.ProjectID, zone, instanceName); err != nil {
			t.Errorf("deleteInstance got err: %v", err)
		}
		if err := deleteReservation(&buf, tc.ProjectID, zone, reservationName); err != nil {
			t.Errorf("deleteReservation got err: %v", err)
		}
	})

	t.Run("Create template without consuming reservation", func(t *testing.T) {
		buf.Reset()
		reservationName := fmt.Sprintf("test-reservation-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		templateName := fmt.Sprintf("test-instance-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		if err = createBaseReservation(&buf, tc.ProjectID, zone, reservationName); err != nil {
			t.Errorf("createBaseReservation got err: %v", err)
		}

		ctx := context.Background()
		reservationsClient, err := compute.NewReservationsRESTClient(ctx)
		if err != nil {
			t.Errorf("reservationsClient got err: %v", err)
		}
		defer reservationsClient.Close()

		err = createTemplateNotConsumeReservation(&buf, tc.ProjectID, templateName)
		if err != nil {
			t.Errorf("createTemplateNotConsumeReservation failed: %v", err)
		}

		want := "Instance template created"
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createTemplateNotConsumeReservation got %s, want %s", got, want)
		}

		req := &computepb.GetReservationRequest{
			Project:     tc.ProjectID,
			Zone:        zone,
			Reservation: reservationName,
		}

		resp, err := reservationsClient.Get(ctx, req)
		if err != nil {
			t.Errorf("get reservation failed: %v", err)
		}
		inUseAfter := resp.GetSpecificReservation().GetInUseCount()
		if inUseAfter != 0 {
			t.Errorf("reservation was consumed. Expected 0, got %d", inUseAfter)
		}

		if err = deleteInstanceTemplate(tc.ProjectID, templateName); err != nil {
			t.Errorf("deleteInstanceTemplate got err: %v", err)
		}
		if err := deleteReservation(&buf, tc.ProjectID, zone, reservationName); err != nil {
			t.Errorf("deleteReservation got err: %v", err)
		}
	})
}
