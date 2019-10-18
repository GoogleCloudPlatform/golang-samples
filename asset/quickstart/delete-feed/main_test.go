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
	"strconv"
	"strings"
	"testing"
	"time"

	asset "cloud.google.com/go/asset/apiv1p2beta1"
	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1p2beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMain(t *testing.T) {
	tc := testutil.SystemTest(t)
	env := map[string]string{"GOOGLE_CLOUD_PROJECT": tc.ProjectID}
	feedID := fmt.Sprintf("FEED-%s", strconv.FormatInt(time.Now().UnixNano(), 10))

	ctx := context.Background()
	client, err := asset.NewClient(ctx)
	if err != nil {
		t.Fatalf("asset.NewClient: %v", err)
	}

	feedParent := fmt.Sprintf("projects/%s", tc.ProjectID)
	assetNames := []string{"YOUR_ASSET_NAME"}
	topic := fmt.Sprintf("projects/%s/topics/%s", tc.ProjectID, "YOUR_TOPIC_NAME")

	createTopic(ctx, t, tc.ProjectID, "YOUR_TOPIC_NAME")

	req := &assetpb.CreateFeedRequest{
		Parent: feedParent,
		FeedId: feedID,
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
		t.Fatalf("client.CreateFeed: %v", err)
	}

	m := testutil.BuildMain(t)
	defer m.Cleanup()

	if !m.Built() {
		t.Errorf("failed to build app")
	}

	stdOut, stdErr, err := m.Run(env, 2*time.Minute, fmt.Sprintf("--feed_id=%s", feedID))
	if err != nil {
		t.Errorf("execution failed: %v", err)
	}
	if len(stdErr) > 0 {
		t.Errorf("did not expect stderr output, got %d bytes: %s", len(stdErr), string(stdErr))
	}
	got := string(stdOut)
	want := "Deleted Feed"
	if !strings.Contains(got, want) {
		t.Errorf("stdout returned %s, wanted to contain %s", got, want)
	}
}

func createTopic(ctx context.Context, t *testing.T, projectID, topicName string) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("pubsub.NewClient: %v", err)
	}

	topic := client.Topic(topicName)
	ok, err := topic.Exists(ctx)
	if err != nil {
		t.Fatalf("failed to check if topic exists: %v", err)
	}
	if !ok {
		_, err := client.CreateTopic(ctx, topicName)
		// In case the topic was created in the meantime.
		if status.Code(err) == codes.AlreadyExists {
			return
		}
		if err != nil {
			t.Fatalf("CreateTopic: %v", err)
		}
	}
}
