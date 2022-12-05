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
 
// [START asset_quickstart_create_saved_query]
 
// Sample create-saved-query create saved query.
package main
 
import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
 
	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
)
 
// Command-line flags.
var (
	savedQueryID = flag.String("saved_query_id", "YOUR-QUERY-ID", "Identifier of Saved Query.")
)
 
func main() {
	flag.Parse()
	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		log.Fatalf("asset.NewClient: %v", err)
	}
	defer client.Close()
 
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	parent := fmt.Sprintf("projects/%s", projectID)
	req := &assetpb.CreateSavedQueryRequest{
		Parent: parent,
		SavedQueryId: *savedQueryID,
		SavedQuery: &assetpb.SavedQuery{
			Content: &assetpb.SavedQuery_QueryContent{
				QueryContent: &assetpb.SavedQuery_QueryContent_IamPolicyAnalysisQuery {
					IamPolicyAnalysisQuery: &assetpb.IamPolicyAnalysisQuery {
						Scope: parent,
						AccessSelector: &assetpb.IamPolicyAnalysisQuery_AccessSelector {
							Permissions: []string{"iam.serviceAccount.actAs"},
						},
					},
				},
			},
		}}
	response, err := client.CreateSavedQuery(ctx, req)
	if err != nil {
		log.Fatalf("client.CreateSavedQuery: %v", err)
	}
	fmt.Print(response)
}
 
// [END asset_quickstart_create_saved_query]
