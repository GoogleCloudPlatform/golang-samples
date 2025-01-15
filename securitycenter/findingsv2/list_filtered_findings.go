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

// [START securitycenter_list_filtered_findings_v2]
import (
	"context"
	"fmt"
	"io"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	"cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
	"google.golang.org/api/iterator"
)

// listFilteredFindings prints findings with category 'MEDIUM_RISK_ONE' for a
// specific source to w. sourceName is the full resource name of the source
// to search for findings under.
func listFilteredFindings(w io.Writer, sourceName string) error {
	// Specific source:
	// 		sourceName := "{parent}/sources/{sourceId}"
	// All sources:
	// 		sourceName := "{parent}/sources/-"
	// where,
	// Parent must be in one of the following formats:
	//		"organizations/{orgId}"
	//		"projects/{projectId}"
	//		"folders/{folderId}"
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	req := &securitycenterpb.ListFindingsRequest{
		Parent: sourceName,
		Filter: `category="MEDIUM_RISK_ONE"`,
	}
	it := client.ListFindings(ctx, req)
	for {
		result, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("it.Next: %w", err)
		}
		finding := result.Finding
		fmt.Fprintf(w, "Finding Name: %s, ", finding.Name)
		fmt.Fprintf(w, "Resource Name %s, ", finding.ResourceName)
		fmt.Fprintf(w, "Category: %s\n", finding.Category)
	}
	return nil
}

// [END securitycenter_list_filtered_findings_v2]
