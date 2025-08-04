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

package connectors

// [START managedkafka_restart_connector]
import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/managedkafka/apiv1/managedkafkapb"
	"google.golang.org/api/option"

	managedkafka "cloud.google.com/go/managedkafka/apiv1"
)

func restartConnector(w io.Writer, projectID, region, connectClusterID, connectorID string, opts ...option.ClientOption) error {
	// projectID := "my-project-id"
	// region := "us-central1"
	// connectClusterID := "my-connect-cluster"
	// connectorID := "my-connector"
	ctx := context.Background()
	client, err := managedkafka.NewManagedKafkaConnectClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("managedkafka.NewManagedKafkaConnectClient got err: %w", err)
	}
	defer client.Close()

	connectorPath := fmt.Sprintf("projects/%s/locations/%s/connectClusters/%s/connectors/%s", projectID, region, connectClusterID, connectorID)
	req := &managedkafkapb.RestartConnectorRequest{
		Name: connectorPath,
	}
	resp, err := client.RestartConnector(ctx, req)
	if err != nil {
		return fmt.Errorf("client.RestartConnector got err: %w", err)
	}
	fmt.Fprintf(w, "Restarted connector: %#v\n", resp)
	return nil
}

// [END managedkafka_restart_connector]
