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

package servers

import (
	"bytes"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestInspectString(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)

	{
		if err := createRealm(buf, tc.ProjectID, "global", "myrealm"); err != nil {
			t.Errorf("createRealm: %v", err)
		}

		got := buf.String()
		want := "Realm created: projects/" + tc.ProjectID + "/locations/global/realms/myrealm"
		if got != want {
			t.Errorf("createRealm got %q, want %q", got, want)
		}
	}

	buf.Reset()

	{
		if err := getRealm(buf, tc.ProjectID, "global", "myrealm"); err != nil {
			t.Errorf("getRealm: %v", err)
		}

		got := buf.String()
		want := "Realm retrieved: projects/" + tc.ProjectID + "/locations/global/realms/myrealm"
		if got != want {
			t.Errorf("getRealm got %q, want %q", got, want)
		}
	}

	buf.Reset()

	{
		if err := listRealms(buf, tc.ProjectID, "global"); err != nil {
			t.Errorf("listRealms: %v", err)
		}

		got := buf.String()
		want := "Realm listed: projects/" + tc.ProjectID + "/locations/global/realms/myrealm"
		if got != want {
			t.Errorf("listRealms got %q, want %q", got, want)
		}
	}

	buf.Reset()

	{
		if err := deleteRealm(buf, tc.ProjectID, "global", "myrealm"); err != nil {
			t.Errorf("deleteRealm: %v", err)
		}

		got := buf.String()
		want := "Realm deleted."
		if got != want {
			t.Errorf("deleteRealm got %q, want %q", got, want)
		}
	}
}
