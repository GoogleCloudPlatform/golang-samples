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

package client

// [START bigquery_set_client_endpoint]
import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/option"
)

// setClientEndpoint creates a bigquery.Client pointing
// to a specific endpoint.
func setClientEndpoint(endpoint, projectID string) (*bigquery.Client, error) {
	// projectID := "my-project-id"
	// endpoint := "https://bigquery.googleapis.com/bigquery/v2/
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID, option.WithEndpoint(endpoint))
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	return client, nil
}

// [END bigquery_set_client_endpoint]
