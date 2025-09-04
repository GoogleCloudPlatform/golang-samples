// Copyright 2025 Google LLC
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

package clusters

// [START managedkafka_delete_connect_cluster]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/managedkafka/apiv1/managedkafkapb"
	"google.golang.org/api/option"

	managedkafka "cloud.google.com/go/managedkafka/apiv1"
)

func deleteConnectCluster(w io.Writer, projectID, region, clusterID string, opts ...option.ClientOption) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// clusterID := "my-connect-cluster"
	ctx := context.Background()
	client, err := managedkafka.NewManagedKafkaConnectClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("managedkafka.NewManagedKafkaConnectClient got err: %w", err)
	}
	defer client.Close()

	clusterPath := fmt.Sprintf("projects/%s/locations/%s/connectClusters/%s", projectID, region, clusterID)
	req := &managedkafkapb.DeleteConnectClusterRequest{
		Name: clusterPath,
	}
	op, err := client.DeleteConnectCluster(ctx, req)
	if err != nil {
		return fmt.Errorf("client.DeleteConnectCluster got err: %w", err)
	}
	err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("op.Wait got err: %w", err)
	}
	fmt.Fprint(w, "Deleted connect cluster\n")
	return nil
}

// [END managedkafka_delete_connect_cluster]
