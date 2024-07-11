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

package consumergroups

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
	consumerGroupPrefix = "consumergroup"
	parentClusterID     = "test-cluster"
	region              = "us-central1"
)

func TestConsumerGroups(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	consumerGroupID := fmt.Sprintf("%s-%d", consumerGroupPrefix, time.Now().UnixNano())
	options := fake.Options(t)
	t.Run("GetConsumerGroup", func(t *testing.T) {
		if err := getConsumerGroup(buf, tc.ProjectID, region, parentClusterID, consumerGroupID, options...); err != nil {
			t.Fatalf("failed to get consumer group: %v", err)
		}
		got := buf.String()
		want := "Got consumer group"
		if !strings.Contains(got, want) {
			t.Fatalf("getConsumerGroup() mismatch got: %s\nwant: %s", got, want)
		}
	})
	t.Run("UpdateConsumerGroup", func(t *testing.T) {
		partitionOffset := map[int32]int64{
			1: 10,
		}
		topicPath := "fake-topic-path"
		if err := updateConsumerGroup(buf, tc.ProjectID, region, parentClusterID, consumerGroupID, topicPath, partitionOffset, options...); err != nil {
			t.Fatalf("failed to update consumer group: %v", err)
		}
		got := buf.String()
		want := "Updated consumer group"
		if !strings.Contains(got, want) {
			t.Fatalf("updateConsumerGroup() mismatch got: %s\nwant: %s", got, want)
		}
	})
	t.Run("ListConsumerGroups", func(t *testing.T) {
		if err := listConsumerGroups(buf, tc.ProjectID, region, parentClusterID, options...); err != nil {
			t.Fatalf("failed to list consumer groups: %v", err)
		}
		got := buf.String()
		want := "Got consumer group"
		if !strings.Contains(got, want) {
			t.Fatalf("listConsumerGroups() mismatch got: %s\nwant: %s", got, want)
		}
	})
	t.Run("DeleteConsumerGroup", func(t *testing.T) {
		if err := deleteConsumerGroup(buf, tc.ProjectID, region, parentClusterID, consumerGroupID, options...); err != nil {
			t.Fatalf("failed to delete consumer group: %v", err)
		}
		got := buf.String()
		want := "Deleted consumer group"
		if !strings.Contains(got, want) {
			t.Fatalf("deleteConsumerGroup() mismatch got: %s\nwant: %s", got, want)
		}
	})
}
