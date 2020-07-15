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

package servers

import (
	"bytes"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestRealms(t *testing.T) {
	tc := testutil.SystemTest(t)

	t.Run("create realm", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := createRealm(buf, tc.ProjectID, "global", "myrealm"); err != nil {
			t.Errorf("createRealm: %v", err)
		}

		got := buf.String()
		want := "Realm created: projects/" + tc.ProjectID + "/locations/global/realms/myrealm"
		if got != want {
			t.Errorf("createRealm got %q, want %q", got, want)
		}
	})

	t.Run("get created realm", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := getRealm(buf, tc.ProjectID, "global", "myrealm"); err != nil {
			t.Errorf("getRealm: %v", err)
		}

		got := buf.String()
		want := "Realm retrieved: projects/" + tc.ProjectID + "/locations/global/realms/myrealm"
		if got != want {
			t.Errorf("getRealm got %q, want %q", got, want)
		}
	})

	t.Run("list created realm", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := listRealms(buf, tc.ProjectID, "global"); err != nil {
			t.Errorf("listRealms: %v", err)
		}

		got := buf.String()
		want := "Realm listed: projects/" + tc.ProjectID + "/locations/global/realms/myrealm\n"
		if got != want {
			t.Errorf("listRealms got %q, want %q", got, want)
		}
	})

	t.Run("cluster tests", innerTestGameServerCluster)

	t.Run("delete realm", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := deleteRealm(buf, tc.ProjectID, "global", "myrealm"); err != nil {
			t.Errorf("deleteRealm: %v", err)
		}

		got := buf.String()
		want := "Realm deleted."
		if got != want {
			t.Errorf("deleteRealm got %q, want %q", got, want)
		}
	})

	t.Run("list no realms", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := listRealms(buf, tc.ProjectID, "global"); err != nil {
			t.Errorf("listRealms: %v", err)
		}

		got := buf.String()
		want := ""
		if got != want {
			t.Errorf("listRealms got %q, want %q", got, want)
		}
	})
}

func innerTestGameServerCluster(t *testing.T) {
	tc := testutil.SystemTest(t)

	t.Run("create cluster", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := createCluster(buf, tc.ProjectID, "global", "myrealm", "mycluster", "projects/217093627905/locations/us-central1/clusters/gke-shared-default"); err != nil {
			t.Errorf("createCluster: %v", err)
		}

		got := buf.String()
		want := "Cluster created: projects/" + tc.ProjectID + "/locations/global/realms/myrealm/gameServerClusters/mycluster"
		if got != want {
			t.Errorf("createCluster got %q, want %q", got, want)
		}
	})

	t.Run("get created cluster", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := getGameServerCluster(buf, tc.ProjectID, "global", "myrealm", "mycluster"); err != nil {
			t.Errorf("getGameServerCluster: %v", err)
		}

		got := buf.String()
		want := "Cluster retrieved: projects/" + tc.ProjectID + "/locations/global/realms/myrealm/gameServerClusters/mycluster"
		if got != want {
			t.Errorf("getGameServerCluster got %q, want %q", got, want)
		}
	})

	t.Run("list created cluster", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := listGameServerClusters(buf, tc.ProjectID, "global", "myrealm"); err != nil {
			t.Errorf("listGameServerClusters: %v", err)
		}

		got := buf.String()
		want := "Cluster listed: projects/" + tc.ProjectID + "/locations/global/realms/myrealm/gameServerClusters/mycluster\n"
		if got != want {
			t.Errorf("listGameServerClusters got %q, want %q", got, want)
		}
	})

	t.Run("delete cluster", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := deleteCluster(buf, tc.ProjectID, "global", "myrealm", "mycluster"); err != nil {
			t.Errorf("deleteCluster: %v", err)
		}

		got := buf.String()
		want := "Cluster deleted."
		if got != want {
			t.Errorf("deleteCluster got %q, want %q", got, want)
		}
	})
}
