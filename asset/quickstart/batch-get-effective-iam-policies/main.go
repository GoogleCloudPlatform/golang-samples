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

// [START asset_quickstart_batch_get_effective_iam_policies]

// Sample asset-quickstart batch get effective iam policies.
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
	projectId := flag.String("projectId", "", "Project Id to construct the scope under which the effective IAM policies to get.")
	fullResourceName := flag.String("fullResourceName", "", "Resource on which the IAM policies are effective.")

	flag.Parse()
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		log.Fatalf("asset.NewClient: %v", err)
	}
	defer client.Close()

	req := &assetpb.BatchGetEffectiveIamPoliciesRequest{
		Scope: fmt.Sprintf("projects/%s", *projectId),
		Names: []string{*fullResourceName},
	}

	op, err := client.BatchGetEffectiveIamPolicies(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	for index, result := range op.PolicyResults {
		fmt.Println(index, result)
	}
}

// [END asset_quickstart_batch_get_effective_iam_policies]
