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
	defer client.Close()

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
	defer client.Close()

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
	defer client.Close()

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

func createInstance(projectID, zone, instanceName string) error {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.InsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb:  proto.Int64(375),
						SourceImage: proto.String("projects/debian-cloud/global/images/family/debian-12"),
					},
					Interface:  proto.String("SCSI"),
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
				},
			},
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/%s", zone, "n1-standard-1")),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					Name: proto.String("global/networks/default"),
				},
			},
		},
	}

	op, err := instancesClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create instance: %w", err)
	}

	return op.Wait(ctx)
}

func deleteInstance(projectID, zone, instance string) error {
	ctx := context.Background()
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	req := &computepb.DeleteInstanceRequest{
		Instance: instance,
		Project:  projectID,
		Zone:     zone,
	}
	op, err := client.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to delete instance: %w", err)
	}

	return op.Wait(ctx)
}

func createSpecificSharedReservation(client ClientInterface, projectID, baseProjectId, zone, reservationName string) error {
	ctx := context.Background()

	shareSettings := map[string]*computepb.ShareSettingsProjectConfig{
		projectID: {ProjectId: proto.String(projectID)},
	}

	req := &computepb.InsertReservationRequest{
		Project: baseProjectId,
		ReservationResource: &computepb.Reservation{
			Name: proto.String(reservationName),
			Zone: proto.String(zone),
			SpecificReservation: &computepb.AllocationSpecificSKUReservation{
				Count: proto.Int64(2),
				InstanceProperties: &computepb.AllocationSpecificSKUAllocationReservedInstanceProperties{
					MachineType:    proto.String("n2-standard-32"),
					MinCpuPlatform: proto.String("Intel Cascade Lake"),
				},
			},
			ShareSettings: &computepb.ShareSettings{
				ProjectMap: shareSettings,
				ShareType:  proto.String("SPECIFIC_PROJECTS"),
			},
			SpecificReservationRequired: proto.Bool(true),
		},
		Zone: zone,
	}

	op, err := client.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create reservation: %w", err)
	}

	if op != nil {
		if err = op.Wait(ctx); err != nil {
			return fmt.Errorf("unable to wait for the operation: %w", err)
		}
	}
	return nil
}

func createSpecificConsumableReservation(projectID, zone, reservationName string) error {
	ctx := context.Background()
	reservationsClient, err := compute.NewReservationsRESTClient(ctx)
	if err != nil {
		return err
	}
	defer reservationsClient.Close()

	req := &computepb.InsertReservationRequest{
		Project: projectID,
		ReservationResource: &computepb.Reservation{
			Name: proto.String(reservationName),
			Zone: proto.String(zone),
			SpecificReservation: &computepb.AllocationSpecificSKUReservation{
				Count: proto.Int64(2),
				InstanceProperties: &computepb.AllocationSpecificSKUAllocationReservedInstanceProperties{
					MachineType:    proto.String("n2-standard-32"),
					MinCpuPlatform: proto.String("Intel Cascade Lake"),
				},
			},
			SpecificReservationRequired: proto.Bool(true),
		},
		Zone: zone,
	}

	op, err := reservationsClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create reservation: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	return nil
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

	t.Run("Shared reservation CRUD", func(t *testing.T) {
		reservationName := fmt.Sprintf("test-reservation-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		baseProjectID := tc.ProjectID
		// This test require 2 projects, therefore one of them is mocked.
		// If you want to make a real test, please adjust projectID accordingly and uncomment reservationsClient creation.
		// Make sure that base project has proper permissions to share reservations.
		// See: https://cloud.google.com/compute/docs/instances/reservations-shared#shared_reservation_constraint
		destinationProjectID := "some-project"
		ctx := context.Background()

		want := "Reservation created"

		// Uncomment line below if you want to run the test without mocks
		// reservationsClient, err := compute.NewReservationsRESTClient(ctx)
		reservationsClient := ReservationsClient{}
		if err != nil {
			t.Errorf("Couldn't create reservationsClient, err: %v", err)
		}
		defer reservationsClient.Close()

		if err := createSharedReservation(&buf, reservationsClient, destinationProjectID, baseProjectID, zone, reservationName, *sourceTemplate.SelfLink); err != nil {
			t.Errorf("createSharedReservation got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createSharedReservation got %s, want %s", got, want)
		}
		buf.Reset()

		req := &computepb.DeleteReservationRequest{
			Project:     baseProjectID,
			Reservation: reservationName,
			Zone:        zone,
		}

		_, err = reservationsClient.Delete(ctx, req)
		if err != nil {
			t.Errorf("unable to delete reservation: %v", err)
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

		if err = deleteTemplate(tc.ProjectID, templateName); err != nil {
			t.Errorf("deleteTemplate got err: %v", err)
		}
		if err := deleteReservation(&buf, tc.ProjectID, zone, reservationName); err != nil {
			t.Errorf("deleteReservation got err: %v", err)
		}
	})

	t.Run("Test create from exisiting VM", func(t *testing.T) {
		reservationName := fmt.Sprintf("test-reservation-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		existingVM := fmt.Sprintf("test-instance-%v-%v", time.Now().Format("01-02-2006"), r.Int())

		err := createInstance(tc.ProjectID, zone, existingVM)
		if err != nil {
			t.Fatalf("createInstance got err: %v", err)
		}
		defer deleteInstance(tc.ProjectID, zone, existingVM)

		want := "Reservation created"
		if err := createReservationFromVM(&buf, tc.ProjectID, zone, reservationName, existingVM); err != nil {
			t.Errorf("createReservationFromVM got err: %v", err)
		}
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("createReservationFromVM got %s, want %s", got, want)
		}
		buf.Reset()

		if err := deleteReservation(&buf, tc.ProjectID, zone, reservationName); err != nil {
			t.Errorf("deleteReservation got err: %v", err)
		}
	})
}

func TestConsumeReservations(t *testing.T) {
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	zone := "europe-west2-b"
	instanceName := fmt.Sprintf("test-instance-%v-%v", time.Now().Format("01-02-2006"), r.Int())
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

	t.Run("Consume sprecific shared reservation", func(t *testing.T) {
		reservationName := fmt.Sprintf("test-reservation-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		baseProjectID := tc.ProjectID
		// This test require 2 projects, therefore one of them is mocked.
		// If you want to make a real test, please adjust destinationProjectID accordingly and uncomment reservationsClient creation.
		// Make sure that project has proper permissions to share reservations.
		// See: https://cloud.google.com/compute/docs/instances/reservations-shared#shared_reservation_constraint
		destinationProjectID := "softserve-shared"
		ctx := context.Background()

		reservationsClient := ReservationsClient{}
		// Uncomment lines below if you want to run the test without mocks
		//
		// reservationsClient, err := compute.NewReservationsRESTClient(ctx)
		// if err != nil {
		// 	t.Errorf("Couldn't create reservationsClient, err: %v", err)
		// }
		// defer reservationsClient.Close()

		if err := createSpecificSharedReservation(reservationsClient, baseProjectID, destinationProjectID, zone, reservationName); err != nil {
			t.Errorf("createSpecificSharedReservation got err: %v", err)
		}

		instanceClient := InstanceClient{}
		// Uncomment lines below if you want to run the test without mocks
		//
		// instanceClient, err := compute.NewInstancesRESTClient(ctx)
		// if err != nil {
		// 	t.Errorf("Couldn't create instanceClient, err: %v", err)
		// }
		// defer instanceClient.Close()

		if err := consumeSpecificSharedReservation(&buf, instanceClient, baseProjectID, destinationProjectID, zone, instanceName, reservationName); err != nil {
			t.Errorf("consumeSpecificSharedReservation got err: %v", err)
		}

		want := "Instance created from shared reservation"
		if got := buf.String(); !strings.Contains(got, want) {
			t.Errorf("consumeSpecificSharedReservation got %s, want %s", got, want)
		}

		req := &computepb.DeleteInstanceRequest{
			Instance: instanceName,
			Project:  baseProjectID,
			Zone:     zone,
		}
		_, err := instanceClient.Delete(ctx, req)
		if err != nil {
			t.Errorf("unable to delete reservation: %v", err)
		}

		req2 := &computepb.DeleteReservationRequest{
			Project:     destinationProjectID,
			Reservation: reservationName,
			Zone:        zone,
		}

		_, err = reservationsClient.Delete(ctx, req2)
		if err != nil {
			t.Errorf("unable to delete reservation: %v", err)
		}
	})

	t.Run("Consume any reservation", func(t *testing.T) {
		reservationName := fmt.Sprintf("test-reservation-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		if err = createReservation(&buf, tc.ProjectID, zone, reservationName, *sourceTemplate.SelfLink); err != nil {
			t.Errorf("createConsumableReservation got err: %v", err)
		}

		ctx := context.Background()
		reservationsClient, err := compute.NewReservationsRESTClient(ctx)
		if err != nil {
			t.Errorf("reservationsClient got err: %v", err)
		}
		defer reservationsClient.Close()

		req := &computepb.GetReservationRequest{
			Project:     tc.ProjectID,
			Zone:        zone,
			Reservation: reservationName,
		}
		res, err := reservationsClient.Get(ctx, req)
		if err != nil {
			t.Errorf("get reservation got err: %v", err)
		}

		inUseBefore := res.GetSpecificReservation().GetInUseCount()
		if inUseBefore != 0 {
			t.Error("reservation was consumed beforehand")
		}

		if err = consumeAnyReservation(&buf, tc.ProjectID, zone, instanceName, *sourceTemplate.SelfLink); err != nil {
			t.Errorf("consumeAnyReservation got err: %v", err)
		}

		res2, err := reservationsClient.Get(ctx, req)
		if err != nil {
			t.Errorf("get reservation got err: %v", err)
		}

		inUseAfter := res2.GetSpecificReservation().GetInUseCount()
		if inUseAfter != 1 {
			t.Errorf("Reservation wasn't consumed. Expected 1, got %d", inUseAfter)
		}

		if err = deleteInstance(tc.ProjectID, zone, instanceName); err != nil {
			t.Errorf("deleteInstance got err: %v", err)
		}
		if err := deleteReservation(&buf, tc.ProjectID, zone, reservationName); err != nil {
			t.Errorf("deleteReservation got err: %v", err)
		}
	})

	t.Run("Consume specific reservation", func(t *testing.T) {
		reservationName := fmt.Sprintf("test-reservation-%v-%v", time.Now().Format("01-02-2006"), r.Int())
		if err = createSpecificConsumableReservation(tc.ProjectID, zone, reservationName); err != nil {
			t.Errorf("createConsumableReservation got err: %v", err)
		}

		ctx := context.Background()
		reservationsClient, err := compute.NewReservationsRESTClient(ctx)
		if err != nil {
			t.Errorf("reservationsClient got err: %v", err)
		}
		defer reservationsClient.Close()

		req := &computepb.GetReservationRequest{
			Project:     tc.ProjectID,
			Zone:        zone,
			Reservation: reservationName,
		}
		res, err := reservationsClient.Get(ctx, req)
		if err != nil {
			t.Errorf("get reservation got err: %v", err)
		}

		inUseBefore := res.GetSpecificReservation().GetInUseCount()
		if inUseBefore != 0 {
			t.Error("reservation was consumed beforehand")
		}

		if err = consumeSpecificReservation(&buf, tc.ProjectID, zone, instanceName, reservationName); err != nil {
			t.Errorf("consumeAnyReservation got err: %v", err)
		}

		res2, err := reservationsClient.Get(ctx, req)
		if err != nil {
			t.Errorf("get reservation got err: %v", err)
		}

		inUseAfter := res2.GetSpecificReservation().GetInUseCount()
		if inUseAfter != 1 {
			t.Errorf("Reservation wasn't consumed. Expected 1, got %d", inUseAfter)
		}

		if err = deleteInstance(tc.ProjectID, zone, instanceName); err != nil {
			t.Errorf("deleteInstance got err: %v", err)
		}
		if err := deleteReservation(&buf, tc.ProjectID, zone, reservationName); err != nil {
			t.Errorf("deleteReservation got err: %v", err)
		}
	})
}
