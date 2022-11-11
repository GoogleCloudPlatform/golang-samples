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

// [START asset_quickstart_analyze_iam_policy]

// Sample analyze_iam_policy analyzes accessible IAM policies that match a request.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
)

func main() {
	scope := flag.String("scope", "", "Scope of the analysis.")
	fullResourceName := flag.String("fullResourceName", "", "Query resource.")
	flag.Parse()
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		log.Fatalf("asset.NewClient: %v", err)
	}
	defer client.Close()

	req := &assetpb.AnalyzeIamPolicyRequest{
		AnalysisQuery: &assetpb.IamPolicyAnalysisQuery{
			Scope: *scope,
			ResourceSelector: &assetpb.IamPolicyAnalysisQuery_ResourceSelector{
				FullResourceName: *fullResourceName,
			},
			Options: &assetpb.IamPolicyAnalysisQuery_Options{
				ExpandGroups:     true,
				OutputGroupEdges: true,
			},
		},
	}

	op, err := client.AnalyzeIamPolicy(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	for index, result := range op.MainAnalysis.AnalysisResults {
		fmt.Println(index, result)
	}
}

// [END asset_quickstart_analyze_iam_policy]
