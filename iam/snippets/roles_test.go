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

package snippets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
)

func TestRoles(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	uuid, _ := uuid.NewV4()

	// viewGrantableRoles test.
	fullResourceName := "//cloudresourcemanager.googleapis.com/projects/" + tc.ProjectID
	grantable, err := viewGrantableRoles(buf, fullResourceName)
	if err != nil {
		t.Fatalf("viewGrantableRoles: %v", err)
	}
	if len(grantable) < 1 {
		t.Fatalf("viewGrantableRoles: expected at least 1 item")
	}

	// getRole test.
	role, err := getRole(buf, "roles/appengine.appAdmin")
	if err != nil {
		t.Fatalf("getRole: %v", err)
	}
	wantName := "roles/appengine.appAdmin"
	if role.Name != wantName {
		t.Fatalf("getRole: role.Name is %q, wanted %q", role.Name, wantName)
	}

	// Custom role variables.
	// Name must start with a letter and be 6-30 characters.
	name := "a" + strings.Replace(uuid.String(), "-", "", -1)[:29]
	title := "Test role"
	desc := "Test description"
	perms := []string{"resourcemanager.projects.getIamPolicy"}
	stage := "GA"

	// createRole test.
	role, err = createRole(buf, tc.ProjectID, name, title, desc, stage, perms)
	if err != nil {
		t.Fatalf("createRole: %v", err)
	}
	wantTitle := title
	if wantTitle != role.Title {
		t.Fatalf("createRole: role.Title is %q, wanted %q", role.Name, wantTitle)
	}

	// editRole test.
	newTitle := "Updated test role"
	role, err = editRole(buf, tc.ProjectID, name, newTitle, desc, stage, perms)
	if err != nil {
		t.Fatalf("editRole: %v", err)
	}
	wantTitle = newTitle
	if wantTitle != role.Title {
		t.Fatalf("editRole: role.Title is %q, wanted %q", role.Title, wantTitle)
	}

	// listRoles test.
	roles, err := listRoles(buf, tc.ProjectID)
	if err != nil {
		t.Fatalf("listRoles: %v", err)
	}
	if len(roles) < 1 {
		t.Fatalf("listRoles: expected at least 1 item")
	}

	// deleteRole test.
	err = deleteRole(buf, tc.ProjectID, name)
	if err != nil {
		t.Fatalf("deleteRole: %v", err)
	}
}
