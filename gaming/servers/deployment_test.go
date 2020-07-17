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

func TestGameServerDeployments(t *testing.T) {
	testutil.KnownBadMTLS(t)
	tc := testutil.SystemTest(t)

	t.Run("create deployment", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := createGameServerDeployment(buf, tc.ProjectID, "mydeployment"); err != nil {
			t.Errorf("createGameServerDeployment: %v", err)
		}

		got := buf.String()
		want := "Deployment created: projects/" + tc.ProjectID + "/locations/global/gameServerDeployments/mydeployment"
		if got != want {
			t.Errorf("createGameServerDeployment got %q, want %q", got, want)
		}
	})

	t.Run("get created deployment", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := getGameServerDeployment(buf, tc.ProjectID, "mydeployment"); err != nil {
			t.Errorf("getGameServerDeployment: %v", err)
		}

		got := buf.String()
		want := "Deployment retrieved: projects/" + tc.ProjectID + "/locations/global/gameServerDeployments/mydeployment"
		if got != want {
			t.Errorf("getGameServerDeployment got %q, want %q", got, want)
		}
	})

	t.Run("list created deployment", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := listGameServerDeployments(buf, tc.ProjectID); err != nil {
			t.Errorf("listGameServerDeployments: %v", err)
		}

		got := buf.String()
		want := "Deployment listed: projects/" + tc.ProjectID + "/locations/global/gameServerDeployments/mydeployment\n"
		if got != want {
			t.Errorf("listGameServerDeployments got %q, want %q", got, want)
		}
	})

	t.Run("config tests", innerTestGameServerFleet)

	t.Run("delete deployment", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := deleteGameServerDeployment(buf, tc.ProjectID, "mydeployment"); err != nil {
			t.Errorf("deleteGameServerDeployment: %v", err)
		}

		got := buf.String()
		want := "Deployment deleted."
		if got != want {
			t.Errorf("deleteGameServerDeployment got %q, want %q", got, want)
		}
	})

	t.Run("list no deployments", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := listGameServerDeployments(buf, tc.ProjectID); err != nil {
			t.Errorf("listGameServerDeployments: %v", err)
		}

		got := buf.String()
		want := ""
		if got != want {
			t.Errorf("listGameServerDeployments got %q, want %q", got, want)
		}
	})
}

func innerTestGameServerFleet(t *testing.T) {
	tc := testutil.SystemTest(t)

	t.Run("create config", func(t *testing.T) {
		buf := new(bytes.Buffer)

		if err := createGameServerConfig(buf, tc.ProjectID, "mydeployment", "myconfig"); err != nil {
			t.Errorf("createGameServerConfig: %v", err)
		}

		got := buf.String()
		want := "Config created: projects/" + tc.ProjectID + "/locations/global/gameServerDeployments/mydeployment/configs/myconfig"
		if got != want {
			t.Errorf("createGameServerConfig got %q, want %q", got, want)
		}
	})

	t.Run("get created config", func(t *testing.T) {
		buf := new(bytes.Buffer)

		if err := getGameServerConfig(buf, tc.ProjectID, "mydeployment", "myconfig"); err != nil {
			t.Errorf("getGameServerConfig: %v", err)
		}

		got := buf.String()
		want := "Config retrieved: projects/" + tc.ProjectID + "/locations/global/gameServerDeployments/mydeployment/configs/myconfig"
		if got != want {
			t.Errorf("getGameServerConfig got %q, want %q", got, want)
		}
	})

	t.Run("list created config", func(t *testing.T) {
		buf := new(bytes.Buffer)

		if err := listGameServerConfigs(buf, tc.ProjectID, "mydeployment"); err != nil {
			t.Errorf("listGameServerConfigs: %v", err)
		}

		got := buf.String()
		want := "Config listed: projects/" + tc.ProjectID + "/locations/global/gameServerDeployments/mydeployment/configs/myconfig\n"
		if got != want {
			t.Errorf("listGameServerConfigs got %q, want %q", got, want)
		}
	})

	t.Run("rollout tests", innerTestGameServerDeploymentRollout)

	t.Run("delete config", func(t *testing.T) {
		buf := new(bytes.Buffer)
		if err := deleteGameServerConfig(buf, tc.ProjectID, "mydeployment", "myconfig"); err != nil {
			t.Errorf("deleteGameServerConfig: %v", err)
		}

		got := buf.String()
		want := "Config deleted."
		if got != want {
			t.Errorf("deleteGameServerConfig got %q, want %q", got, want)
		}
	})

	t.Run("list no configs", func(t *testing.T) {
		buf := new(bytes.Buffer)

		if err := listGameServerConfigs(buf, tc.ProjectID, "mydeployment"); err != nil {
			t.Errorf("listGameServerConfigs: %v", err)
		}

		got := buf.String()
		want := ""
		if got != want {
			t.Errorf("listGameServerConfigs got %q, want %q", got, want)
		}
	})
}

func innerTestGameServerDeploymentRollout(t *testing.T) {
	tc := testutil.SystemTest(t)

	t.Run("rollout set default", func(t *testing.T) {
		buf := new(bytes.Buffer)

		if err := updateRolloutDefaultConfig(buf, tc.ProjectID, "mydeployment", "myconfig"); err != nil {
			t.Errorf("updateRolloutDefaultConfig: %v", err)
		}

		got := buf.String()
		want := "Deployment rollout updated: projects/" + tc.ProjectID + "/locations/global/gameServerDeployments/mydeployment"
		if got != want {
			t.Errorf("updateRolloutDefaultConfig got %q, want %q", got, want)
		}
	})

	t.Run("get rollout with default", func(t *testing.T) {
		buf := new(bytes.Buffer)

		if err := getGameServerDeploymentRollout(buf, tc.ProjectID, "mydeployment"); err != nil {
			t.Errorf("getGameServerDeploymentRollout: %v", err)
		}

		got := buf.String()
		want := "Rollout default: projects/" + tc.ProjectID + "/locations/global/gameServerDeployments/mydeployment/configs/myconfig\n"
		if got != want {
			t.Errorf("getGameServerDeploymentRollout got %q, want %q", got, want)
		}
	})

	t.Run("rollout remove default", func(t *testing.T) {
		buf := new(bytes.Buffer)

		if err := updateRolloutClearDefaultConfig(buf, tc.ProjectID, "mydeployment"); err != nil {
			t.Errorf("updateRolloutClearDefaultConfig: %v", err)
		}

		got := buf.String()
		want := "Deployment rollout updated: projects/" + tc.ProjectID + "/locations/global/gameServerDeployments/mydeployment"
		if got != want {
			t.Errorf("updateRolloutClearDefaultConfig got %q, want %q", got, want)
		}
	})

	t.Run("get rollout with no default", func(t *testing.T) {
		buf := new(bytes.Buffer)

		if err := getGameServerDeploymentRollout(buf, tc.ProjectID, "mydeployment"); err != nil {
			t.Errorf("getGameServerDeploymentRollout: %v", err)
		}

		got := buf.String()
		want := "Rollout default: \n"
		if got != want {
			t.Errorf("getGameServerDeploymentRollout got %q, want %q", got, want)
		}
	})

	t.Run("rollout set override", func(t *testing.T) {
		buf := new(bytes.Buffer)

		if err := updateRolloutOverrideConfig(buf, tc.ProjectID, "global", "myrealm", "mydeployment", "myconfig"); err != nil {
			t.Errorf("updateRolloutOverrideConfig: %v", err)
		}

		got := buf.String()
		want := "Deployment rollout updated: projects/" + tc.ProjectID + "/locations/global/gameServerDeployments/mydeployment"
		if got != want {
			t.Errorf("updateRolloutOverrideConfig got %q, want %q", got, want)
		}
	})

	t.Run("get rollout with override", func(t *testing.T) {
		buf := new(bytes.Buffer)

		if err := getGameServerDeploymentRollout(buf, tc.ProjectID, "mydeployment"); err != nil {
			t.Errorf("getGameServerDeploymentRollout: %v", err)
		}

		got := buf.String()
		want := "Rollout default: \nOverride these realms [projects/" + tc.ProjectID + "/locations/global/realms/myrealm]with this config: projects/" + tc.ProjectID + "/locations/global/gameServerDeployments/mydeployment/configs/myconfig\n"
		if got != want {
			t.Errorf("getGameServerDeploymentRollout got %q, want %q", got, want)
		}
	})
	t.Run("rollout remove override", func(t *testing.T) {
		buf := new(bytes.Buffer)

		if err := updateRolloutClearOverrideConfig(buf, tc.ProjectID, "mydeployment"); err != nil {
			t.Errorf("updateRolloutClearOverrideConfig: %v", err)
		}

		got := buf.String()
		want := "Deployment rollout updated: projects/" + tc.ProjectID + "/locations/global/gameServerDeployments/mydeployment"
		if got != want {
			t.Errorf("updateRolloutClearOverrideConfig got %q, want %q", got, want)
		}
	})

	t.Run("get rollout with no override", func(t *testing.T) {
		buf := new(bytes.Buffer)

		if err := getGameServerDeploymentRollout(buf, tc.ProjectID, "mydeployment"); err != nil {
			t.Errorf("getGameServerDeploymentRollout: %v", err)
		}

		got := buf.String()
		want := "Rollout default: \n"
		if got != want {
			t.Errorf("getGameServerDeploymentRollout got %q, want %q", got, want)
		}
	})
}
