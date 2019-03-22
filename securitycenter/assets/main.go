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

// Contains examples calls to Cloud Security Center ListAssets API method.

package main

import (
	securitycenter "cloud.google.com/go/securitycenter/apiv1beta1"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/api/iterator"
	securitycenterpb "google.golang.org/genproto/googleapis/cloud/securitycenter/v1beta1"
	"os"
	"time"
)

// Lists all assets currently in an organization.
// [START list_all_assets]
func ListAllAssets(orgId string) (int, error) {
	// Instantiate a context and a security service client to make API calls.
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		fmt.Printf("Error instantiating client %v\n", err)
		return 0, err
	}

	listAssetsRequest := securitycenterpb.ListAssetsRequest{}
	// orgId is the numeric organization ID.  e.g. 01231231
	listAssetsRequest.Parent = fmt.Sprintf("organizations/%s", orgId)
	var assetsFound = 0

	// Call the service
	it := client.ListAssets(ctx, &listAssetsRequest)
loop:
	for {
		r, err := it.Next()
		switch err {
		case nil:
			fmt.Printf("Asset Name: %s, Resource Name %s, Resource Type %s\n", r.Asset.Name, r.Asset.SecurityCenterProperties.ResourceName, r.Asset.SecurityCenterProperties.ResourceType)
			assetsFound++
		case iterator.Done:
			break loop
		default:
			fmt.Printf("Error listing assets: %v", err)
			return assetsFound, err
		}
	}
	return assetsFound, nil
}

// [END list_all_assets]

// List all current assets in an organization that are GCP projects.
// [START list_project_assets]
func ListAllProjectAssets(orgId string) (int, error) {
	// Initialize a new context and client and return any errors (e.g
	// credentials were not found)
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		fmt.Printf("Error instantiating client %v\n", err)
		return 0, err
	}

	listAssetsRequest := &securitycenterpb.ListAssetsRequest{}
	// orgId is the numeric organization ID.  e.g. 01231231
	listAssetsRequest.Parent = fmt.Sprintf("organizations/%s", orgId)
	listAssetsRequest.Filter = "security_center_properties.resource_type=\"google.cloud.resourcemanager.Project\""

	it := client.ListAssets(ctx, listAssetsRequest)
	var assetsFound = 0
	// Process the results from the iterator.
loop:
	for {
		r, err := it.Next()
		switch err {
		case nil:
			fmt.Printf("Asset Name: %s, Resource Name %s, Resource Type %s\n", r.Asset.Name, r.Asset.SecurityCenterProperties.ResourceName, r.Asset.SecurityCenterProperties.ResourceType)

			assetsFound++
		case iterator.Done:
			break loop
		default:
			fmt.Printf("Error listing assets: %v", err)
			return 0, err
		}
	}
	return assetsFound, nil
}

// [END list_project_assets]

// List all assets in an organization that were GCP projects at asOf.
// [START list_project_assets_at_time]
func ListAllProjectAssetsAtTime(orgId string, asOf time.Time) (int, error) {
	// Initialize a new context and client and return any errors (e.g
	// credentials were not found)
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		fmt.Printf("Error instantiating client %v\n", err)
		return 0, err
	}

	listAssetsRequest := &securitycenterpb.ListAssetsRequest{}
	// orgId is the numeric organization ID.  e.g. 01231231
	listAssetsRequest.Parent = fmt.Sprintf("organizations/%s", orgId)
	listAssetsRequest.Filter = "security_center_properties.resource_type=\"google.cloud.resourcemanager.Project\""
	// Convert the time to a Timestamp protobuf
	readTime, err := ptypes.TimestampProto(asOf)
	if err != nil {
		fmt.Printf("Error converting %v: %v", asOf, err)
		return 0, err
	}
	listAssetsRequest.ReadTime = readTime

	it := client.ListAssets(ctx, listAssetsRequest)
	var assetsFound = 0
	// Process the results from the iterator.
loop:
	for {
		r, err := it.Next()
		switch err {
		case nil:
			fmt.Printf("Asset Name: %s, Resource Name %s, Resource Type %s\n", r.Asset.Name, r.Asset.SecurityCenterProperties.ResourceName, r.Asset.SecurityCenterProperties.ResourceType)

			assetsFound++
		case iterator.Done:
			break loop
		default:
			fmt.Printf("Error listing assets: %v", err)
			return 0, err
		}
	}
	return assetsFound, nil
}

// [END list_project_assets_at_time]

// [START list_project_assets_with_state_changes]
func ListAllProjectAssetsWithStateChanges(orgId string) (int, error) {
	// Initialize a new context and client and return any errors (e.g
	// credentials were not found)
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		fmt.Printf("Error instantiating client %v\n", err)
		return 0, err
	}

	listAssetsRequest := &securitycenterpb.ListAssetsRequest{}
	// orgId is the numeric organization ID.  e.g. 01231231
	listAssetsRequest.Parent = fmt.Sprintf("organizations/%s", orgId)
	listAssetsRequest.Filter = "security_center_properties.resource_type=\"google.cloud.resourcemanager.Project\""
	// Convert the time (10 days) to a Duration protobuf.
	daysAgo, _ := time.ParseDuration("240h")
	fmt.Printf("Days ago: %v\n", daysAgo)
	// Compare all current assets to there state as of 10 days ago.
	listAssetsRequest.CompareDuration = ptypes.DurationProto(daysAgo)

	it := client.ListAssets(ctx, listAssetsRequest)
	var assetsFound = 0
	// Process the results from the iterator.
loop:
	for {
		r, err := it.Next()
		switch err {
		case nil:
			fmt.Printf("Asset Name: %s, Resource Name %s, Resource Type %s, State Change %s\n", r.Asset.Name, r.Asset.SecurityCenterProperties.ResourceName, r.Asset.SecurityCenterProperties.ResourceType, r.State)

			assetsFound++
		case iterator.Done:
			break loop
		default:
			fmt.Printf("Error listing assets: %v", err)
			return 0, err
		}
	}
	return assetsFound, nil
}

// [END list_project_assets_with_state_changes]

func main() {
	orgId := os.Getenv("GCLOUD_ORGANIZATION")
	if orgId == "" {
		fmt.Fprintf(os.Stderr, "GCLOUD_ORGANIZATION environment variable must be set.\n")
	}
	fmt.Println("Project All Assets:")
	ListAllAssets(orgId)
	fmt.Println("Project Assets:")
	ListAllProjectAssets(orgId)
	fmt.Println("List all projects at a certain time:")
	ListAllProjectAssetsAtTime(orgId, time.Date(2019, time.March, 18, 0, 0, 0, 0, time.UTC))
	fmt.Println("List all projects with change from a day ago:")
	ListAllProjectAssetsWithStateChanges(orgId)

}
