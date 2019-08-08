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

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	asset "cloud.google.com/go/asset/apiv1p2beta1"
    assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1p2beta1"
)

func TestMain(t *testing.T) {
	tc := testutil.SystemTest(t)
	os.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

    // Set a feed for delete feed
	ctx := context.Background()
    client, err := asset.NewClient(ctx)
    if err != nil {
            log.Fatal(err)
    }
    
    feedParent := fmt.Sprintf("projects/%s", tc.ProjectID) 
    feedId := "YOUR_FEED_ID"
    assetNames :=  []string{"YOUR_ASSET_NAME"}
    topic := fmt.Sprintf("projects/%s/topics/%s", tc.ProjectID, "YOUR_TOPIC_NAME")
    
    req := &assetpb.CreateFeedRequest{
            Parent: feedParent,
            FeedId: feedId,
            Feed: &assetpb.Feed{
              AssetNames: assetNames,
              FeedOutputConfig: &assetpb.FeedOutputConfig{
                Destination: &assetpb.FeedOutputConfig_PubsubDestination{
                  PubsubDestination: &assetpb.PubsubDestination{
                    Topic: topic,
                  },
                },
              },
            }}
    _, err = client.CreateFeed(ctx, req)
    if err != nil {
            log.Fatal(err)
    }

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = oldStdout

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read stdout: %v", err)
	}
	got := string(out)

	want := "Deleted Feed"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}
}
