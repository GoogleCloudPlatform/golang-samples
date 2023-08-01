// Copyright 2023 Google LLC
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

package muteconfig

// [START securitycenter_list_mute_configs]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"cloud.google.com/go/securitycenter/apiv1/securitycenterpb"
	"google.golang.org/api/iterator"
)

// listMuteRules lists mute configs at the organization level will return all the configs
// at the org, folder, and project levels.
// Similarly, listing configs at folder level will list all the configs
// at the folder and project levels.
func listMuteRules(w io.Writer, parent string) error {
	// Use any one of the following resource paths to list mute configurations:
	//         - organizations/{organization_id}
	//         - folders/{folder_id}
	//         - projects/{project_id}
	// parent := fmt.Sprintf("projects/%s", "your-google-cloud-project-id")
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycenterpb.ListMuteConfigsRequest{Parent: parent}

	// List all mute configs present in the resource.
	it := client.ListMuteConfigs(ctx, req)
	for {
		muteconfig, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("it.Next: %w", err)
		}
		fmt.Fprintf(w, "Muteconfig Name: %s, ", muteconfig.Name)
	}
	return nil
}

// [END securitycenter_list_mute_configs]
