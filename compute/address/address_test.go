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

	"google.golang.org/protobuf/proto"

	compute "cloud.google.com/go/compute/apiv1"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// helper functions

func createTestInstance(projectID, zone, instanceName string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	imagesClient, err := compute.NewImagesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewImagesRESTClient: %w", err)
	}
	defer imagesClient.Close()

	// List of public operating system (OS) images: https://cloud.google.com/compute/docs/images/os-details.
	newestDebianReq := &computepb.GetFromFamilyImageRequest{
		Project: "debian-cloud",
		Family:  "debian-12",
	}
	newestDebian, err := imagesClient.GetFromFamily(ctx, newestDebianReq)
	if err != nil {
		return fmt.Errorf("unable to get image from family: %w", err)
	}

	req := &computepb.InsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		InstanceResource: &computepb.Instance{
			Name: proto.String(instanceName),
			Disks: []*computepb.AttachedDisk{
				{
					InitializeParams: &computepb.AttachedDiskInitializeParams{
						DiskSizeGb:  proto.Int64(10),
						SourceImage: newestDebian.SelfLink,
						DiskType:    proto.String(fmt.Sprintf("zones/%s/diskTypes/pd-standard", zone)),
					},
					AutoDelete: proto.Bool(true),
					Boot:       proto.Bool(true),
					Type:       proto.String(computepb.AttachedDisk_PERSISTENT.String()),
				},
			},
			MachineType: proto.String(fmt.Sprintf("zones/%s/machineTypes/n1-standard-1", zone)),
			NetworkInterfaces: []*computepb.NetworkInterface{
				{
					//Name: proto.String("global/networks/default"),
					AccessConfigs: []*computepb.AccessConfig{
						{
							Type:        proto.String(computepb.AccessConfig_ONE_TO_ONE_NAT.String()),
							Name:        proto.String("External NAT"),
							NetworkTier: proto.String(computepb.AccessConfig_PREMIUM.String()),
						},
					},
				},
			},
		},
	}

	op, err := instancesClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create instance: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	return nil
}

func deleteInstance(projectID, zone, instanceName string) error {
	ctx := context.Background()
	client, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.DeleteInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	op, err := client.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteInstance: %w", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait for DeleteInstance operation: %w", err)
	}

	return nil
}

func listIPAddresses(ctx context.Context, projectID, region string) ([]string, error) {
	var addresses []string
	var err error

	if region == "" {
		addresses, err = _listGlobalIPAddresses(ctx, projectID)
	} else {
		addresses, err = _listRegionalIPAddresses(ctx, projectID, region)
	}

	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func _listGlobalIPAddresses(ctx context.Context, projectID string) ([]string, error) {
	client, err := compute.NewGlobalAddressesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewGlobalAddressesRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.ListGlobalAddressesRequest{
		Project: projectID,
	}

	it := client.List(ctx, req)
	var addresses []string
	for {
		address, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("ListGlobalAddresses: %w", err)
		}
		addresses = append(addresses, address.GetName())
	}

	return addresses, nil
}

func _listRegionalIPAddresses(ctx context.Context, projectID, region string) ([]string, error) {
	client, err := compute.NewAddressesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewAddressesRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.ListAddressesRequest{
		Project: projectID,
		Region:  region,
	}

	it := client.List(ctx, req)
	var addresses []string
	for {
		address, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("ListAddresses: %w", err)
		}
		addresses = append(addresses, address.GetName())
	}

	return addresses, nil
}

func deleteIPAddress(ctx context.Context, projectID, region, addressName string) error {

	if region == "" {
		return _deleteGlobalIPAddress(ctx, projectID, addressName)
	}
	return _deleteRegionalIPAddress(ctx, projectID, region, addressName)
}

func _deleteGlobalIPAddress(ctx context.Context, projectID, addressName string) error {
	client, err := compute.NewGlobalAddressesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewGlobalAddressesRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.DeleteGlobalAddressRequest{
		Project: projectID,
		Address: addressName,
	}

	op, err := client.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteGlobalAddress: %w", err)
	}

	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("Wait for DeleteGlobalAddress operation: %w", err)
	}

	return nil
}

func _deleteRegionalIPAddress(ctx context.Context, projectID, region, addressName string) error {
	client, err := compute.NewAddressesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewAddressesRESTClient: %w", err)
	}
	defer client.Close()

	req := &computepb.DeleteAddressRequest{
		Project: projectID,
		Region:  region,
		Address: addressName,
	}

	op, err := client.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteAddress: %w", err)
	}

	if err := op.Wait(ctx); err != nil {
		return fmt.Errorf("Wait for DeleteAddress operation: %w", err)
	}

	return nil
}

// end helper functions

func TestReserveNewRegionalExternal(t *testing.T) {
	ctx := context.Background()
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	addressName := "test-address-" + fmt.Sprint(seededRand.Int())
	region := "us-central1"

	buf := &bytes.Buffer{}

	defer func() {
		if err := deleteIPAddress(ctx, tc.ProjectID, region, addressName); err != nil {
			t.Errorf("deleteIPAddress got err: %v", err)
		}
	}()

	got, err := reserveNewRegionalExternal(buf, tc.ProjectID, region, addressName, true)
	if err != nil {
		t.Errorf("reserveNewExternal got err: %v", err)
	}
	if !strings.Contains(*got.Region, region) {
		t.Errorf("returned global IP address region is not equal to requested: %s != %s", *got.Region, region)
	}
	if *got.NetworkTier != computepb.AccessConfig_PREMIUM.String() {
		t.Errorf("returned global IP network tier is different from requested")
	}

	expectedResult := fmt.Sprintf("Regional address %v reserved: %s", addressName, *got.Address)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("reserveNewExternal got %q, want %q", got, expectedResult)
	}

	// List IP addresses and verify the new address is in the list
	addresses, err := listIPAddresses(ctx, tc.ProjectID, region)
	if err != nil {
		t.Errorf("listIPAddresses got err: %v", err)
	}

	found := false
	for _, address := range addresses {
		if address == addressName {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Address %v not found in the list of reserved addresses", addressName)
	}

}

func TestReserveNewGlobalExternal(t *testing.T) {
	ctx := context.Background()
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	addressName := "test-global-address-" + fmt.Sprint(seededRand.Int())

	buf := &bytes.Buffer{}

	defer func() {
		if err := deleteIPAddress(ctx, tc.ProjectID, "", addressName); err != nil {
			t.Errorf("deleteGlobalIPAddress got err: %v", err)
		}
	}()

	got, err := reserveNewGlobalExternal(buf, tc.ProjectID, addressName, true)
	if err != nil {
		t.Errorf("reserveNewGlobalExternal got err: %v", err)
	}
	if got.Region != nil {
		t.Error("returned global IP address region is not global")
	}
	if *got.IpVersion != computepb.Address_IPV6.String() {
		t.Errorf("returned global IP version is different from requested")
	}

	expectedResult := fmt.Sprintf("Global address %v reserved: %s", addressName, *got.Address)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("reserveNewGlobalExternal got %q, want %q", got, expectedResult)
	}

	// List global IP addresses and verify the new address is in the list
	addresses, err := listIPAddresses(ctx, tc.ProjectID, "")
	if err != nil {
		t.Errorf("listGlobalIPAddresses got err: %v", err)
	}

	found := false
	for _, address := range addresses {
		if address == addressName {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Global address %v not found in the list of reserved addresses", addressName)
	}
}

func TestReleaseGlobalExternal(t *testing.T) {
	ctx := context.Background()
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	addressName := "test-global-address-" + fmt.Sprint(seededRand.Int())

	buf := &bytes.Buffer{}

	defer func() {
		if err := deleteIPAddress(ctx, tc.ProjectID, "", addressName); err != nil {
			// cleanup in case of unexpected result
		}
	}()

	_, err := reserveNewGlobalExternal(buf, tc.ProjectID, addressName, true)
	if err != nil {
		t.Errorf("reserveNewGlobalExternal got err: %v", err)
	}
	buf.Reset()

	if err := releaseGlobalStaticExternal(buf, tc.ProjectID, addressName); err != nil {
		t.Errorf("releaseGlobalStaticExternal got err: %v", err)
	}

	expectedResult := "Static external IP address released"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("releaseGlobalStaticExternal got %q, want %q", got, expectedResult)
	}

	// List global IP addresses and verify the new address is not in the list
	addresses, err := listIPAddresses(ctx, tc.ProjectID, "")
	if err != nil {
		t.Errorf("listGlobalIPAddresses got err: %v", err)
	}

	for _, address := range addresses {
		if address == addressName {
			t.Errorf("Global address %v found in the list of reserved addresses", addressName)
		}
	}

}

func TestReleaseRegionalExternal(t *testing.T) {
	ctx := context.Background()
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	addressName := "test-address-" + fmt.Sprint(seededRand.Int())
	region := "us-central1"

	buf := &bytes.Buffer{}

	defer func() {
		if err := deleteIPAddress(ctx, tc.ProjectID, region, addressName); err != nil {
			// cleanup in case of unexpected result
		}
	}()

	_, err := reserveNewRegionalExternal(buf, tc.ProjectID, region, addressName, true)
	if err != nil {
		t.Errorf("reserveNewRegionalExternal got err: %v", err)
	}
	buf.Reset()

	if err := releaseRegionalStaticExternal(buf, tc.ProjectID, region, addressName); err != nil {
		t.Errorf("releaseRegionalStaticExternal got err: %v", err)
	}

	expectedResult := "Static external IP address released"
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("releaseRegionalStaticExternal got %q, want %q", got, expectedResult)
	}

	// List regional IP addresses and verify the new address is not in the list
	addresses, err := listIPAddresses(ctx, tc.ProjectID, region)
	if err != nil {
		t.Errorf("listIPAddresses got err: %v", err)
	}

	for _, address := range addresses {
		if address == addressName {
			t.Errorf("Regional address %v found in the list of reserved addresses", addressName)
		}
	}

}

func TestReleaseNonExistentExternal(t *testing.T) {
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	addressName := "test-address-" + fmt.Sprint(seededRand.Int())
	region := "us-central1"

	buf := &bytes.Buffer{}

	err := releaseRegionalStaticExternal(buf, tc.ProjectID, region, addressName)
	if err == nil {
		t.Errorf("releaseRegionalStaticExternal should have returned an error")
	}
	buf.Reset()

	err = releaseGlobalStaticExternal(buf, tc.ProjectID, addressName)
	if err == nil {
		t.Errorf("releaseGlobalStaticExternal should have returned an error")
	}

}

func TestGetGlobalExternal(t *testing.T) {
	ctx := context.Background()
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	addressName := "test-global-address-" + fmt.Sprint(seededRand.Int())

	buf := &bytes.Buffer{}

	defer func() {
		if err := deleteIPAddress(ctx, tc.ProjectID, "", addressName); err != nil {
			t.Errorf("deleteGlobalIPAddress got err: %v", err)
		}
	}()

	haveAddress, err := reserveNewGlobalExternal(buf, tc.ProjectID, addressName, true)
	if err != nil {
		t.Errorf("reserveNewGlobalExternal got err: %v", err)
	}
	buf.Reset()

	gotAddress, err := getGlobalExternal(buf, tc.ProjectID, addressName)
	if err != nil {
		t.Errorf("getGlobalExternal got err: %v", err)
	}

	if gotAddress.Region != nil {
		t.Error("returned global IP address region is not global")
	}

	expectedResult := fmt.Sprintf("Global address %v has external IP address: %v", *haveAddress.Name, *haveAddress.Address)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("getGlobalExternal got %q, want %q", got, expectedResult)
	}
}

func TestGetRegionalExternal(t *testing.T) {
	ctx := context.Background()
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	addressName := "test-address-" + fmt.Sprint(seededRand.Int())
	region := "us-central1"

	defer func() {
		if err := deleteIPAddress(ctx, tc.ProjectID, region, addressName); err != nil {
			t.Errorf("deleteIPAddress got err: %v", err)
		}
	}()

	buf := &bytes.Buffer{}

	haveAddress, err := reserveNewRegionalExternal(buf, tc.ProjectID, region, addressName, true)
	if err != nil {
		t.Errorf("reserveNewRegionalExternal got err: %v", err)
	}
	buf.Reset()

	gotAddress, err := getRegionalExternal(buf, tc.ProjectID, region, addressName)
	if err != nil {
		t.Errorf("getRegionalExternal got err: %v", err)
	}

	if !strings.Contains(*gotAddress.Region, *haveAddress.Region) {
		t.Errorf("returned regional IP address region is not equal to requested: %s != %s", *gotAddress.Region, *haveAddress.Region)
	}

	expectedResult := fmt.Sprintf("Regional address %v has external IP address: %v", *haveAddress.Name, *haveAddress.Address)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("getRegionalExternal got %q, want %q", got, expectedResult)
	}

}

func TestAssignUnassignStaticAddressToExistingVM(t *testing.T) {
	ctx := context.Background()
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	instanceName := "test-instance-" + fmt.Sprint(seededRand.Int())
	addressName := "test-address-" + fmt.Sprint(seededRand.Int())
	zone := "us-central1-a"
	region := "us-central1"
	buf := &bytes.Buffer{}

	// initiate instance
	err := createTestInstance(tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Fatalf("createTestInstance got err: %v", err)
		return
	}

	defer func() {
		if err := deleteInstance(tc.ProjectID, zone, instanceName); err != nil {
			t.Errorf("deleteInstance got err: %v", err)
		}

	}()

	// create and retrieve address
	address, err := reserveNewRegionalExternal(buf, tc.ProjectID, region, addressName, true)
	if err != nil {
		t.Fatalf("reserveNewRegionalExternal got err: %v", err)
		return
	}

	defer func() {

		if err := deleteIPAddress(ctx, tc.ProjectID, region, addressName); err != nil {
			t.Errorf("deleteIPAddress got err: %v", err)
		}
	}()

	// assign address to test
	if err := assignStaticAddressToExistingVM(buf, tc.ProjectID, zone, instanceName, address.GetAddress(), "nic0"); err != nil {
		t.Errorf("assignStaticAddressToExistingVM got err: %v", err)
	}

	// verify output
	expectedResult := fmt.Sprintf("Static address %s assigned to the instance %s", address.GetAddress(), instanceName)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("assignStaticAddressToExistingVM got %q, want %q", got, expectedResult)
	}

	reqGet := &computepb.GetInstanceRequest{
		Project:  tc.ProjectID,
		Zone:     zone,
		Instance: instanceName,
	}

	// verify address assign
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	instance, err := instancesClient.Get(ctx, reqGet)
	if err != nil {
		t.Errorf("instancesClient.Get got err: %v", err)
	}

	for _, ni := range instance.NetworkInterfaces {
		if *ni.Name != "nic0" {
			continue
		}
		for _, ac := range ni.AccessConfigs {

			if ac.NatIP != nil && *ac.NatIP == address.GetAddress() {
				return // address assign verified
			}
		}

	}
	t.Error("IP address did not assigned properly") // address assign not verified

	// unassign address
	if err := unassignStaticAddressFromExistingVM(buf, tc.ProjectID, zone, instanceName, "nic0"); err != nil {
		t.Errorf("unassignStaticAddressFromExistingVM got err: %v", err)
	}

	// verify output
	expectedResult = fmt.Sprintf("Static address %s unassigned from the instance %s", address.GetAddress(), instanceName)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("unassignStaticAddressFromExistingVM got %q, want %q", got, expectedResult)
	}

	// verify address unassign
	instance, err = instancesClient.Get(ctx, reqGet)
	if err != nil {
		t.Errorf("instancesClient.Get got err: %v", err)
	}

	for _, ni := range instance.NetworkInterfaces {
		if *ni.Name != "nic0" {
			continue
		}
		for _, ac := range ni.AccessConfigs {
			if ac.NatIP != nil && *ac.NatIP == address.GetAddress() {
				t.Errorf("address %v still found in the list of assigned addresses", addressName)
				return
			}
		}
	}
}

func TestAssignStaticExternalToNewVM(t *testing.T) {
	ctx := context.Background()
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	instanceName := "test-instance-" + fmt.Sprint(seededRand.Int())
	addressName := "test-address-" + fmt.Sprint(seededRand.Int())
	zone := "us-central1-a"
	region := "us-central1"
	buf := &bytes.Buffer{}

	// create and retrieve address
	address, err := reserveNewRegionalExternal(buf, tc.ProjectID, region, addressName, true)
	if err != nil {
		t.Errorf("reserveNewRegionalExternal got err: %v", err)
		return
	}

	defer func() {
		if err := deleteIPAddress(ctx, tc.ProjectID, region, addressName); err != nil {
			t.Errorf("deleteIPAddress got err: %v", err)
		}
	}()

	// assign address to test
	if err := assignStaticExternalToNewVM(buf, tc.ProjectID, zone, instanceName, address.GetAddress()); err != nil {
		t.Errorf("assignStaticExternalToNewVM got err: %v", err)
	}

	defer func() {
		if err := deleteInstance(tc.ProjectID, zone, instanceName); err != nil {
			t.Errorf("deleteInstance got err: %v", err)
		}

	}()

	// verify output
	expectedResult := fmt.Sprintf("Static address %s assigned to new VM", address.GetAddress())
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("assignStaticExternalToNewVM got %q, want %q", got, expectedResult)
	}

	// verify address assign
	reqGet := &computepb.GetInstanceRequest{
		Project:  tc.ProjectID,
		Zone:     zone,
		Instance: instanceName,
	}

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	instance, err := instancesClient.Get(ctx, reqGet)
	if err != nil {
		t.Errorf("instancesClient.Get got err: %v", err)
	}

	for _, ni := range instance.NetworkInterfaces {
		if *ni.Name != "nic0" {
			continue
		}
		for _, ac := range ni.AccessConfigs {

			if ac.NatIP != nil && *ac.NatIP == address.GetAddress() {
				return // address assign verified
			}
		}

	}
	t.Error("IP address did not assigned properly") // address assign not verified

}

func TestPromoteEphemeralAddress(t *testing.T) {
	ctx := context.Background()
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	tc := testutil.SystemTest(t)
	instanceName := "test-instance-" + fmt.Sprint(seededRand.Int())
	addressName := "address-promoted-" + fmt.Sprint(seededRand.Int())
	zone := "us-central1-a"
	region := "us-central1"
	buf := &bytes.Buffer{}

	// initiate instance
	err := createTestInstance(tc.ProjectID, zone, instanceName)
	if err != nil {
		t.Errorf("createTestInstance got err: %v", err)
		return
	}

	defer func() {
		if err := deleteInstance(tc.ProjectID, zone, instanceName); err != nil {
			t.Errorf("deleteInstance got err: %v", err)
		}
	}()

	reqGet := &computepb.GetInstanceRequest{
		Project:  tc.ProjectID,
		Zone:     zone,
		Instance: instanceName,
	}

	// verify address assign
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	instance, err := instancesClient.Get(ctx, reqGet)
	if err != nil {
		t.Errorf("instancesClient.Get got err: %v", err)
	}

	ephemeralIP := ""

	for _, ni := range instance.NetworkInterfaces {
		if *ni.Name != "nic0" {
			continue
		}
		for _, ac := range ni.AccessConfigs {
			if *ac.Type == computepb.AccessConfig_ONE_TO_ONE_NAT.String() {
				ephemeralIP = *ac.NatIP
				break
			}
		}
	}

	if ephemeralIP == "" {
		t.Error("no ephemeral IP found in instance")
		return
	}

	// List IP addresses to verify the new address is NOT in the list yet
	addresses, err := listIPAddresses(ctx, tc.ProjectID, region)
	if err != nil {
		t.Errorf("listIPAddresses got err: %v", err)
	}

	for _, address := range addresses {
		if strings.Contains(address, addressName) {
			t.Error("ephemeral IP address already promoted")
			return
		}
	}

	// promote ephemeral address
	if err := promoteEphemeralAddress(buf, tc.ProjectID, region, ephemeralIP, addressName); err != nil {
		t.Errorf("promoteEphemeralAddress got err: %v", err)
	}

	// release static ip
	defer func() {
		if err = releaseRegionalStaticExternal(buf, tc.ProjectID, region, addressName); err != nil {
			t.Errorf("releaseRegionalStaticExternal got err: %v", err)
		}
	}()

	// verify output
	expectedResult := fmt.Sprintf("Ephemeral IP %s address promoted successfully", ephemeralIP)
	if got := buf.String(); !strings.Contains(got, expectedResult) {
		t.Errorf("promoteEphemeralAddress got %q, want %q", got, expectedResult)
	}

	// List IP addresses to verify the new address is in the list already
	addresses, err = listIPAddresses(ctx, tc.ProjectID, region)
	if err != nil {
		t.Errorf("listIPAddresses got err: %v", err)
	}

	for _, address := range addresses {
		if strings.Contains(address, addressName) {
			return
		}
	}
	t.Error("ephemeral IP was not promoted")
}

func TestGetInstanceIPAddresses(t *testing.T) {
	buf := &bytes.Buffer{}
	instance := &computepb.Instance{
		NetworkInterfaces: []*computepb.NetworkInterface{
			{
				NetworkIP: proto.String("10.128.0.1"),
				AccessConfigs: []*computepb.AccessConfig{
					{
						Type:  proto.String(computepb.AccessConfig_ONE_TO_ONE_NAT.String()),
						NatIP: proto.String("34.68.123.45"),
					},
					{
						Type:  proto.String(computepb.AccessConfig_ONE_TO_ONE_NAT.String()),
						NatIP: proto.String("34.68.123.46"),
					},
				},
				Ipv6AccessConfigs: []*computepb.AccessConfig{
					{
						Type:         proto.String(computepb.AccessConfig_DIRECT_IPV6.String()),
						ExternalIpv6: proto.String("2600:1901:0:1234::"),
					},
				},
				Ipv6Address: proto.String("2600:1901:0:5678::"),
			},
		},
	}

	tests := []struct {
		name        string
		addressType computepb.Address_AddressType
		isIPV6      bool
		want        []string
	}{
		{
			name:        "External IPv4",
			addressType: computepb.Address_EXTERNAL,
			isIPV6:      false,
			want:        []string{"34.68.123.45", "34.68.123.46"},
		},
		{
			name:        "Internal IPv4",
			addressType: computepb.Address_INTERNAL,
			isIPV6:      false,
			want:        []string{"10.128.0.1"},
		},
		{
			name:        "External IPv6",
			addressType: computepb.Address_EXTERNAL,
			isIPV6:      true,
			want:        []string{"2600:1901:0:1234::"},
		},
		{
			name:        "Internal IPv6",
			addressType: computepb.Address_INTERNAL,
			isIPV6:      true,
			want:        []string{"2600:1901:0:5678::"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getInstanceIPAddresses(buf, instance, tt.addressType, tt.isIPV6)
			if len(got) != len(tt.want) {
				t.Errorf("getInstanceIPAddresses() = %v, want %v", got, tt.want)
				return
			}
			for i, ip := range got {
				if ip != tt.want[i] {
					t.Errorf("getInstanceIPAddresses() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
