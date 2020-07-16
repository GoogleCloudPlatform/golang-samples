// Copyright 2020 Google LLC
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

package servers

// [START cloud_game_servers_realm_list]

import (
	"context"
	"fmt"
	"io"

	gaming "cloud.google.com/go/gaming/apiv1beta"
	"google.golang.org/api/iterator"
	gamingpb "google.golang.org/genproto/googleapis/cloud/gaming/v1beta"
)

// listRealms lists the realms in a location.
func listRealms(w io.Writer, projectID, location string) error {
	// projectID := "my-project"
	// location := "global"
	ctx := context.Background()
	client, err := gaming.NewRealmsClient(ctx)
	if err != nil {
		return fmt.Errorf("NewRealmsClient: %v", err)
	}
	defer client.Close()

	req := &gamingpb.ListRealmsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
	}

	it := client.ListRealms(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Next: %v", err)
		}

		fmt.Fprintf(w, "Realm listed: %v\n", resp.Name)
	}

	return nil
}

// [END cloud_game_servers_realm_list]
