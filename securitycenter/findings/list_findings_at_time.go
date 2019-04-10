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

// [START list_findings_at_time]
import (
	"context"
	"fmt"
	"io"
	"time"

	securitycenter "cloud.google.com/go/securitycenter/apiv1"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/api/iterator"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1"
)

// listFindingsAtTime prints findings that where present for a specific source
// as of five days ago to w. sourceName is the full resource name of the
// source to search for findings under.
func listFindingsAtTime(w io.Writer, sourceName string) error {
	// Specific source.
	// sourceName := "organizations/111122222444/sources/1234"
	// All sources.
	// sourceName := "organizations/111122222444/sources/-"
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %v", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.
	fiveDaysAgo, err := ptypes.TimestampProto(time.Now().AddDate(0, 0, -5))
	if err != nil {
		return fmt.Errorf("Error converting five days ago: %v", err)
	}

	req := &securitycenterpb.ListFindingsRequest{
		Parent:   sourceName,
		ReadTime: fiveDaysAgo,
	}
	it := client.ListFindings(ctx, req)
	for {
		result, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("it.Next: %v", err)
		}
		finding := result.Finding
		fmt.Fprintf(w, "Finding Name: %s, ", finding.Name)
		fmt.Fprintf(w, "Resource Name %s, ", finding.ResourceName)
		fmt.Fprintf(w, "Category: %s\n", finding.Category)
	}
	return nil
}

// [END list_findings_at_time]
