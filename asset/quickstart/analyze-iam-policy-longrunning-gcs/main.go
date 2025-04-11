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

// [START asset_quickstart_analyze_iam_policy_longrunning_gcs]

// Sample analyze_iam_policy_longrunning analyzes accessible IAM policies that match a request.
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
	uri := flag.String("uri", "", "Output GCS uri.")
	flag.Parse()
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		log.Fatalf("asset.NewClient: %v", err)
	}
	defer client.Close()

	req := &assetpb.AnalyzeIamPolicyLongrunningRequest{
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
		OutputConfig: &assetpb.IamPolicyAnalysisOutputConfig{
			Destination: &assetpb.IamPolicyAnalysisOutputConfig_GcsDestination_{
				GcsDestination: &assetpb.IamPolicyAnalysisOutputConfig_GcsDestination{
					Uri: *uri,
				},
			},
		},
	}

	op, err := client.AnalyzeIamPolicyLongrunning(ctx, req)
	if err != nil {
		log.Fatalf("AnalyzeIamPolicyLongrunning: %v", err)
	}
	fmt.Print(op.Metadata())

	// Wait for the longrunning operation complete.
	resp, err := op.Wait(ctx)
	if err != nil && !op.Done() {
		fmt.Println("failed to fetch operation status", err)
		return
	}
	if err != nil && op.Done() {
		fmt.Println("operation completed with error", err)
		return
	}
	fmt.Println("operation completed successfully", resp)
}

// [END asset_quickstart_analyze_iam_policy_longrunning_gcs]
