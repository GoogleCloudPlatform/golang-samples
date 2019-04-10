// Copyright 2019 Google LLC
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

package findings

// [START list_sources]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"google.golang.org/api/iterator"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1"
)

// listSources prints all sources in  orgID to w.  orgID is the numeric
// identifier of the organization.
func listSources(w io.Writer, orgID string) error {
	// orgID := "12321311"
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %v", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	req := &securitycenterpb.ListSourcesRequest{
		Parent: fmt.Sprintf("organizations/%s", orgID),
	}
	it := client.ListSources(ctx, req)
	for {
		source, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("it.Next: %v", err)
		}
		fmt.Fprintf(w, "Source Name: %s, ", source.Name)
		fmt.Fprintf(w, "Display name: %s, ", source.DisplayName)
		fmt.Fprintf(w, "Description: %s\n", source.Description)
	}
	return nil
}

// [END list_sources]
