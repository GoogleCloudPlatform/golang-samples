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
	clusterPrefix = "cluster"
	region        = "us-central1"
)

func TestClusters(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	clusterID := fmt.Sprintf("%s-%d", clusterPrefix, time.Now().UnixNano())
	options := fake.Options(t)
	t.Run("CreateCluster", func(t *testing.T) {
		subnet := fmt.Sprintf("projects/%s/regions/%s/subnetworks/default", tc.ProjectID, region)
		vcpuCount := 3
		memoryBytes := 3221225472
		if err := createCluster(buf, tc.ProjectID, region, clusterID, subnet, int64(vcpuCount), int64(memoryBytes), options...); err != nil {
			t.Fatalf("failed to create a cluster: %v", err)
		}
		got := buf.String()
		want := "Created cluster"
		if !strings.Contains(got, want) {
			t.Fatalf("createCluster() mismatch got: %s\nwant: %s", got, want)
		}
	})
	t.Run("GetCluster", func(t *testing.T) {
		if err := getCluster(buf, tc.ProjectID, region, clusterID, options...); err != nil {
			t.Fatalf("failed to get cluster: %v", err)
		}
		got := buf.String()
		want := "Got cluster"
		if !strings.Contains(got, want) {
			t.Fatalf("getCluster() mismatch got: %s\nwant: %s", got, want)
		}
	})
	t.Run("UpdateCluster", func(t *testing.T) {
		memoryBytes := 3221225475
		if err := updateCluster(buf, tc.ProjectID, region, clusterID, int64(memoryBytes), options...); err != nil {
			t.Fatalf("failed to update cluster: %v", err)
		}
		got := buf.String()
		want := "Updated cluster"
		if !strings.Contains(got, want) {
			t.Fatalf("updateCluster() mismatch got: %s\nwant: %s", got, want)
		}
	})
	t.Run("ListClusters", func(t *testing.T) {
		if err := listClusters(buf, tc.ProjectID, region, options...); err != nil {
			t.Fatalf("failed to list clusters: %v", err)
		}
		got := buf.String()
		want := "Got cluster"
		if !strings.Contains(got, want) {
			t.Fatalf("listClusters() mismatch got: %s\nwant: %s", got, want)
		}
	})
	t.Run("DeleteCluster", func(t *testing.T) {
		if err := deleteCluster(buf, tc.ProjectID, region, clusterID, options...); err != nil {
			t.Fatalf("failed to delete cluster: %v", err)
		}
		got := buf.String()
		want := "Deleted cluster"
		if !strings.Contains(got, want) {
			t.Fatalf("deleteCluster() mismatch got: %s\nwant: %s", got, want)
		}
	})
}
