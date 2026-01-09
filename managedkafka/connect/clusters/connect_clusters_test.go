// Copyright 2025 Google LLC
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

package clusters

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
	connectClusterPrefix = "connect-cluster"
	region               = "us-central1"
)

func TestConnectClusters(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	connectClusterID := fmt.Sprintf("%s-%d", connectClusterPrefix, time.Now().UnixNano())
	kafkaClusterID := fmt.Sprintf("kafka-cluster-%d", time.Now().UnixNano())
	options := fake.Options(t)

	// First create a Kafka cluster that the Connect cluster will reference
	kafkaClusterPath := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", tc.ProjectID, region, kafkaClusterID)

	t.Run("CreateConnectCluster", func(t *testing.T) {
		buf.Reset()
		if err := createConnectCluster(buf, tc.ProjectID, region, connectClusterID, kafkaClusterPath, options...); err != nil {
			t.Fatalf("failed to create a connect cluster: %v", err)
		}
		got := buf.String()
		want := "Created connect cluster"
		if !strings.Contains(got, want) {
			t.Fatalf("createConnectCluster() mismatch got: %s\nwant: %s", got, want)
		}
	})
	t.Run("GetConnectCluster", func(t *testing.T) {
		buf.Reset()
		if err := getConnectCluster(buf, tc.ProjectID, region, connectClusterID, options...); err != nil {
			t.Fatalf("failed to get connect cluster: %v", err)
		}
		got := buf.String()
		want := "Got connect cluster"
		if !strings.Contains(got, want) {
			t.Fatalf("getConnectCluster() mismatch got: %s\nwant: %s", got, want)
		}
	})
	t.Run("UpdateConnectCluster", func(t *testing.T) {
		buf.Reset()
		memoryBytes := int64(25769803776) // 24 GiB in bytes
		labels := map[string]string{"environment": "test"}
		if err := updateConnectCluster(buf, tc.ProjectID, region, connectClusterID, memoryBytes, labels, options...); err != nil {
			t.Fatalf("failed to update connect cluster: %v", err)
		}
		got := buf.String()
		want := "Updated connect cluster"
		if !strings.Contains(got, want) {
			t.Fatalf("updateConnectCluster() mismatch got: %s\nwant: %s", got, want)
		}
	})
	t.Run("ListConnectClusters", func(t *testing.T) {
		buf.Reset()
		if err := listConnectClusters(buf, tc.ProjectID, region, options...); err != nil {
			t.Fatalf("failed to list connect clusters: %v", err)
		}
		got := buf.String()
		want := "Got connect cluster"
		if !strings.Contains(got, want) {
			t.Fatalf("listConnectClusters() mismatch got: %s\nwant: %s", got, want)
		}
	})
	t.Run("DeleteConnectCluster", func(t *testing.T) {
		buf.Reset()
		if err := deleteConnectCluster(buf, tc.ProjectID, region, connectClusterID, options...); err != nil {
			t.Fatalf("failed to delete connect cluster: %v", err)
		}
		got := buf.String()
		want := "Deleted connect cluster"
		if !strings.Contains(got, want) {
			t.Fatalf("deleteConnectCluster() mismatch got: %s\nwant: %s", got, want)
		}
	})
}
