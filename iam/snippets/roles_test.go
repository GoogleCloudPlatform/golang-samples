// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
)

func TestRoles(t *testing.T) {
	buf := &bytes.Buffer{}
	uuid, _ := uuid.NewV4()
	project := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")

	// viewGrantableRoles test.
	fullResourceName := "//cloudresourcemanager.googleapis.com/projects/" + project
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
	role, err = createRole(buf, name, project, title, desc, perms, stage)
	if err != nil {
		t.Fatalf("createRole: %v", err)
	}
	wantTitle := title
	if wantTitle != role.Title {
		t.Fatalf("createRole: role.Title is %q, wanted %q", role.Name, wantTitle)
	}

	// editRole test.
	newTitle := "Updated test role"
	role, err = editRole(buf, name, project, newTitle, desc, perms, stage)
	if err != nil {
		t.Fatalf("editRole: %v", err)
	}
	wantTitle = newTitle
	if wantTitle != role.Title {
		t.Fatalf("editRole: role.Title is %q, wanted %q", role.Title, wantTitle)
	}

	// listRoles test.
	roles, err := listRoles(buf, project)
	if err != nil {
		t.Fatalf("listRoles: %v", err)
	}
	if len(roles) < 1 {
		t.Fatalf("listRoles: expected at least 1 item")
	}

	// deleteRole test.
	err = deleteRole(buf, name, project)
	if err != nil {
		t.Fatalf("deleteRole: %v", err)
	}
}
