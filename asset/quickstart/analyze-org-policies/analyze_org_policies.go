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

// [START asset_quickstart_analyze_org_policies]

package create

import (
	"context"
	"fmt"
	"io"
	"log"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
	"google.golang.org/api/iterator"
)

func analyzeOrgPolicies(w io.Writer, scope string, constraint string) error {
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	req := &assetpb.AnalyzeOrgPoliciesRequest{
		Scope:      scope,
		Constraint: constraint,
	}
	it := client.AnalyzeOrgPolicies(ctx, req)

	// Traverse and print the first 10 org policy results in response
	for {
		response, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(w, response)
	}
	return nil
}

// [END asset_quickstart_analyze_org_policies]
