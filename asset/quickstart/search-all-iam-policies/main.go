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

// [START asset_quickstart_search_all_iam_policies]

// Sample search-all-iam-policies searches all IAM policies within the given scope.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	asset "cloud.google.com/go/asset/apiv1"
	"google.golang.org/api/iterator"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
)

func main() {
	scope := flag.String("scope", "", "Scope of the search.")
	query := flag.String("query", "", "Query statement.")
	flag.Parse()
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		log.Fatalf("asset.NewClient: %v", err)
	}
	defer client.Close()

	pageSize := 20
	pageToken := ""

	req := &assetpb.SearchAllIamPoliciesRequest{
		Scope:     *scope,
		Query:     *query,
		PageSize:  int32(pageSize),
		PageToken: pageToken,
	}
	it := client.SearchAllIamPolicies(ctx, req)
	var policies []*assetpb.IamPolicySearchResult
	_, err = iterator.NewPager(it, pageSize, "").NextPage(&policies)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(policies)
}

// [END asset_quickstart_search_all_iam_policies]
