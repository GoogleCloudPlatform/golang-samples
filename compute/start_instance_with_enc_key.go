// Copyright 2021 Google LLC
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

// [START compute_start_enc_instance]
import (
	"context"
	"fmt"
	"io"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

// startInstanceWithEncKey starts a stopped Google Compute Engine instance (with encrypted disks).
func startInstanceWithEncKey(w io.Writer, projectID, zone, instanceName, key string) error {
	// projectID := "your_project_id"
	// zone := "europe-central2-b"
	// instanceName := "your_instance_name"
	// key := "your_encryption_key"

	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %v", err)
	}
	defer instancesClient.Close()

	instanceReq := &computepb.GetInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
	}

	instance, err := instancesClient.Get(ctx, instanceReq)
	if err != nil {
		return fmt.Errorf("unable to get instance: %v", err)
	}

	req := &computepb.StartWithEncryptionKeyInstanceRequest{
		Project:  projectID,
		Zone:     zone,
		Instance: instanceName,
		InstancesStartWithEncryptionKeyRequestResource: &computepb.InstancesStartWithEncryptionKeyRequest{
			Disks: []*computepb.CustomerEncryptionKeyProtectedDisk{
				{
					Source: instance.Disks[0].Source,
					DiskEncryptionKey: &computepb.CustomerEncryptionKey{
						RawKey: proto.String(key),
					},
				},
			},
		},
	}

	op, err := instancesClient.StartWithEncryptionKey(ctx, req)
	if err != nil {
		return fmt.Errorf("unable to start instance with encryption key: %v", err)
	}

	zoneOperationsClient, err := compute.NewZoneOperationsRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewZoneOperationsRESTClient: %v", err)
	}
	defer zoneOperationsClient.Close()

	for {
		waitReq := &computepb.WaitZoneOperationRequest{
			Operation: op.Proto().GetName(),
			Project:   projectID,
			Zone:      zone,
		}
		zoneOp, err := zoneOperationsClient.Wait(ctx, waitReq)
		if err != nil {
			return fmt.Errorf("unable to wait for the operation: %v", err)
		}

		if *zoneOp.Status.Enum() == computepb.Operation_DONE {
			fmt.Fprintf(w, "Instance with encryption key started\n")
			break
		}
	}

	return nil
}

// [END compute_start_enc_instance]
