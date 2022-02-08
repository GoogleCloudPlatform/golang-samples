// Copyright 2022 Google LLC
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

package main

// [START asset_quickstart_analyze_iam_policy_longrunning_gcs]
import (
	"context"
	"fmt"

	asset "cloud.google.com/go/asset/apiv1"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
)

// analyzeIAMPolicyGCS analyzes accessible IAM policies that match a request.
func analyzeIAMPolicyGCS(scope, fullResourceName, gcsURI string) error {
	// scope := "projects/my-project-id"
	// fullResourceName := "//cloudresourcemanager.googleapis.com/projects/project/my-project-id"
	// gcsURI := "gs://bucket_name/object_name"

	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("asset.NewClient: %v", err)
	}
	defer client.Close()

	req := &assetpb.AnalyzeIamPolicyLongrunningRequest{
		AnalysisQuery: &assetpb.IamPolicyAnalysisQuery{
			Scope: scope,
			ResourceSelector: &assetpb.IamPolicyAnalysisQuery_ResourceSelector{
				FullResourceName: fullResourceName,
			},
			Options: &assetpb.IamPolicyAnalysisQuery_Options{
				ExpandGroups:     true,
				OutputGroupEdges: true,
			},
		},
		OutputConfig: &assetpb.IamPolicyAnalysisOutputConfig{
			Destination: &assetpb.IamPolicyAnalysisOutputConfig_GcsDestination_{
				GcsDestination: &assetpb.IamPolicyAnalysisOutputConfig_GcsDestination{
					Uri: gcsURI,
				},
			},
		},
	}

	op, err := client.AnalyzeIamPolicyLongrunning(ctx, req)
	if err != nil {
		return fmt.Errorf("client.AnalyzeIamPolicyLongrunning: %v", err)
	}
	fmt.Print(op.Metadata())

	// Wait for the longrunning operation complete.
	resp, err := op.Wait(ctx)
	if err != nil && !op.Done() {
		return fmt.Errorf("failed to fetch operation status: %v", err)
	}
	if err != nil && op.Done() {
		return fmt.Errorf("operation completed with error: %v", err)
	}
	fmt.Printf("operation completed successfully: %v\n", resp)
	return nil
}

// [END asset_quickstart_analyze_iam_policy_longrunning_gcs]
