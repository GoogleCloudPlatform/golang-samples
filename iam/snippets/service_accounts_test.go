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

func TestServiceAccounts(t *testing.T) {
	buf := &bytes.Buffer{}
	uuid, _ := uuid.NewV4()
	// Name must start with a letter and be 6-30 characters.
	name := "a" + strings.Replace(uuid.String(), "-", "", -1)[:29]
	project := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")

	// createServiceAccount test.
	account, err := createServiceAccount(buf, project, name, "Test")
	if err != nil {
		t.Fatalf("createServiceAccount: %v", err)
	}
	wantEmail := name + "@" + project + ".iam.gserviceaccount.com"
	if wantEmail != account.Email {
		t.Fatalf("createServiceAccount: account.Email is %q, wanted %q", account.Email, wantEmail)
	}

	// renameServiceAccount test.
	account, err = renameServiceAccount(buf, account.Email, "Updated Test")
	if err != nil {
		t.Fatalf("renameServiceAccount: %v", err)
	}
	wantDispName := "Updated Test"
	if wantDispName != account.DisplayName {
		t.Fatalf("renameServiceAccount: account.DisplayName is %q, wanted %q", account.Name, wantDispName)
	}

	// listServiceAccounts test
	accounts, err := listServiceAccounts(buf, project)
	if err != nil {
		t.Fatalf("listServiceAccounts: %v", err)
	}
	if len(accounts) < 1 {
		t.Fatalf("listServiceAccounts: expected at least 1 item")
	}

	// createKey test.
	key, err := createKey(buf, account.Email)
	if err != nil {
		t.Fatalf("createKey: %v", err)
	}
	if key == nil {
		t.Fatalf("createKey: wanted a key but got nil")
	}

	// listKeys test.
	keys, err := listKeys(buf, account.Email)
	if err != nil {
		t.Fatalf("listKeys: %v", err)
	}
	if len(keys) < 1 {
		t.Fatalf("listKeys: expected at least 1 item")
	}

	// deleteKey test.
	err = deleteKey(buf, key.Name)
	if err != nil {
		t.Fatalf("deleteKey: %v", err)
	}

	// deleteServiceAccount test.
	err = deleteServiceAccount(buf, account.Email)
	if err != nil {
		t.Fatalf("deleteServiceAccount: %v", err)
	}
}
