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

package findingsv2

// [START securitycenter_group_filtered_findings_v2]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	"cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
	"google.golang.org/api/iterator"
)

// groupFindingsWithFilter groups findings by category with a filter.
func groupFindingsWithFilter(w io.Writer, sourceName string) error {
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close()

	req := &securitycenterpb.GroupFindingsRequest{
		Parent:  sourceName,
		GroupBy: "category",
		Filter:  `state="ACTIVE"`,
	}

	it := client.GroupFindings(ctx, req)
	for i := 0; ; i++ {
		groupResult, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("it.Next: %w", err)
		}
		fmt.Fprintf(w, "Grouped Finding %d: State: %v, Count: %d\n", i+1, groupResult.Properties["state"], groupResult.Count)
	}
	return nil
}

// [END securitycenter_group_filtered_findings_v2]
