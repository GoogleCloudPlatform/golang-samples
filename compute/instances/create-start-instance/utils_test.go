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

import (
	"context"

	compute "cloud.google.com/go/compute/apiv1"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

func waitZoneOp(ctx context.Context, projectId, zone string, op compute.Operation) error {
	zoneOperationsClient, err := compute.NewZoneOperationsRESTClient(ctx)
	if err != nil {
		return err
	}
	defer zoneOperationsClient.Close()

	for {
		waitReq := &computepb.WaitZoneOperationRequest{
			Operation: op.Proto().GetName(),
			Project:   projectId,
			Zone:      zone,
		}
		zoneOp, err := zoneOperationsClient.Wait(ctx, waitReq)
		if err != nil {
			return err
		}

		if *zoneOp.Status.Enum() == computepb.Operation_DONE {
			break
		}
	}

	return nil
}

func deleteInstance(ctx context.Context, projectId, zone, instanceName string) error {
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return err
	}
	req := &computepb.DeleteInstanceRequest{
		Project:  projectId,
		Zone:     zone,
		Instance: instanceName,
	}

	op, err := instancesClient.Delete(ctx, req)
	if err != nil {
		return err
	}

	return waitZoneOp(ctx, projectId, zone, *op)
}
