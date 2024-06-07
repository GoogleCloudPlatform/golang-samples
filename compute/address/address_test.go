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
	"google.golang.org/api/iterator"

	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

// helper functions
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
