// Copyright 2024 Google LLC
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

package topics

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/managedkafka/fake"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	topicPrefix     = "topic"
	parentClusterID = "test-cluster"
	region          = "us-central1"
)

func TestTopics(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	topicID := fmt.Sprintf("%s-%d", topicPrefix, time.Now().UnixNano())
	options := fake.Options(t)
	t.Run("CreateTopic", func(t *testing.T) {
		partitionCount := 10
		replicationFactor := 3
		configs := map[string]string{
			"min.insync.replicas": "1",
		}
		if err := createTopic(buf, tc.ProjectID, region, parentClusterID, topicID, int32(partitionCount), int32(replicationFactor), configs, options...); err != nil {
			t.Fatalf("failed to create a topic: %v", err)
		}
		got := buf.String()
		want := "Created topic"
		if !strings.Contains(got, want) {
			t.Fatalf("createTopic() mismatch got: %s\nwant: %s", got, want)
		}
	})
	t.Run("GetTopic", func(t *testing.T) {
		if err := getTopic(buf, tc.ProjectID, region, parentClusterID, topicID, options...); err != nil {
			t.Fatalf("failed to get topic: %v", err)
		}
		got := buf.String()
		want := "Got topic"
		if !strings.Contains(got, want) {
			t.Fatalf("getTopic() mismatch got: %s\nwant: %s", got, want)
		}
	})
	t.Run("UpdateTopic", func(t *testing.T) {
		partitionCount := 20
		configs := map[string]string{
			"min.insync.replicas": "2",
		}
		if err := updateTopic(buf, tc.ProjectID, region, parentClusterID, topicID, int32(partitionCount), configs, options...); err != nil {
			t.Fatalf("failed to update topic: %v", err)
		}
		got := buf.String()
		want := "Updated topic"
		if !strings.Contains(got, want) {
			t.Fatalf("updateTopic() mismatch got: %s\nwant: %s", got, want)
		}
	})
	t.Run("ListTopics", func(t *testing.T) {
		if err := listTopics(buf, tc.ProjectID, region, parentClusterID, options...); err != nil {
			t.Fatalf("failed to list topics: %v", err)
		}
		got := buf.String()
		want := "Got topic"
		if !strings.Contains(got, want) {
			t.Fatalf("listTopics() mismatch got: %s\nwant: %s", got, want)
		}
	})
	t.Run("DeleteTopic", func(t *testing.T) {
		if err := deleteTopic(buf, tc.ProjectID, region, parentClusterID, topicID, options...); err != nil {
			t.Fatalf("failed to delete topic: %v", err)
		}
		got := buf.String()
		want := "Deleted topic"
		if !strings.Contains(got, want) {
			t.Fatalf("deleteTopic() mismatch got: %s\nwant: %s", got, want)
		}
	})
}
