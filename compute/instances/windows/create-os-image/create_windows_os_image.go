// Copyright 2022 Google LLC
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

// [START compute_windows_image_create]
import (
	"context"
	"fmt"
	"io"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/protobuf/proto"
)

// createWindowsOSImage creates a new Windows image from the specified source disk.
func createWindowsOSImage(
	w io.Writer,
	projectID, zone, sourceDiskName, imageName, storageLocation string,
	forceCreate bool,
) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// sourceDiskName := "your_source_disk_name"
	// imageName := "your_image_name"
	// storageLocation := "eu"
	// forceCreate := false

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
	disksClient, err := compute.NewDisksRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewDisksRESTClient: %w", err)
	}
	defer disksClient.Close()

	// Getting instances where source disk is attached
	diskRequest := &computepb.GetDiskRequest{
		Project: projectID,
		Zone:    zone,
		Disk:    sourceDiskName,
	}

	sourceDisk, err := disksClient.Get(ctx, diskRequest)
	if err != nil {
		return fmt.Errorf("unable to get disk: %w", err)
	}

	// Сhecking whether the instances is stopped
	for _, fullInstanceName := range sourceDisk.GetUsers() {
		parsedName := strings.Split(fullInstanceName, "/")
		l := len(parsedName)
		if l < 5 {
			return fmt.Errorf(
				"API returned instance name with unexpected format",
			)
		}
		instanceReq := &computepb.GetInstanceRequest{
			Project:  parsedName[l-5],
			Zone:     parsedName[l-3],
			Instance: parsedName[l-1],
		}
		instance, err := instancesClient.Get(ctx, instanceReq)
		if err != nil {
			return fmt.Errorf("unable to get instance: %w", err)
		}

		if instance.GetStatus() != "TERMINATED" && instance.GetStatus() != "STOPPED" {
			if !forceCreate {
				return fmt.Errorf("instance %s should be stopped. "+
					"Please stop the instance using "+
					"GCESysprep command or set forceCreate parameter to true "+
					"(not recommended). More information here: "+
					"https://cloud.google.com/compute/docs/instances/windows/creating-windows-os-image#api",
					parsedName[l-1],
				)
			}
		}
	}

	if forceCreate {
		fmt.Fprintf(w, "Warning: ForceCreate option compromise the integrity of your image. "+
			"Stop the instance before you create the image if possible.",
		)
	}

	req := &computepb.InsertImageRequest{
		Project:     projectID,
		ForceCreate: &forceCreate,
		ImageResource: &computepb.Image{
			Name:             proto.String(imageName),
			SourceDisk:       proto.String(fmt.Sprintf("zones/%s/disks/%s", zone, sourceDiskName)),
			StorageLocations: []string{storageLocation},
		},
	}

	op, err := imagesClient.Insert(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to create image: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("unable to wait for the operation: %w", err)
	}

	fmt.Fprintf(w, "Image created\n")

	return nil
}

// [END compute_windows_image_create]
