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

package clientv2

// [START securitycenter_set_client_endpoint_v2]
import (
	"context"
	"fmt"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	"google.golang.org/api/option"
)

// createClientWithEndpoint creates a Security Command Center client for a
// regional endpoint.
func createClientWithEndpoint(repLocation string) error {
	// Assemble the regional endpoint URL using provided location.
	repEndpoint := fmt.Sprintf("securitycenter.%s.rep.googleapis.com:443", repLocation)
	// Instantiate client for regional endpoint. Use this client to access resources that
	// are subject to data residency controls, and that are located in the region
	// specified in repLocation.
	repCtx := context.Background()
	repClient, err := securitycenter.NewClient(repCtx, option.WithEndpoint(repEndpoint))
	if err != nil {
		return err
	}
	defer repClient.Close()

	return nil
}

// [END securitycenter_set_client_endpoint_v2]
