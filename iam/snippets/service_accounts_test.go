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
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/gofrs/uuid"
)

func TestServiceAccounts(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := &bytes.Buffer{}
	uuid, _ := uuid.NewV4()
	// Name must start with a letter and be 6-30 characters.
	name := "a" + strings.Replace(uuid.String(), "-", "", -1)[:29]

	// createServiceAccount test.
	account, err := createServiceAccount(buf, tc.ProjectID, name, "Test")
	if err != nil {
		t.Fatalf("createServiceAccount: %v", err)
	}
	wantEmail := name + "@" + tc.ProjectID + ".iam.gserviceaccount.com"
	if wantEmail != account.Email {
		t.Fatalf("createServiceAccount: account.Email is %q, wanted %q", account.Email, wantEmail)
	}

	// renameServiceAccount test.

	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		newAccount, err := renameServiceAccount(buf, account.Email, "Updated Test")
		if err != nil {
			r.Errorf("renameServiceAccount: %v", err)
			return
		}
		wantDispName := "Updated Test"
		if wantDispName != newAccount.DisplayName {
			r.Errorf("renameServiceAccount: account.DisplayName is %q, wanted %q", newAccount.Name, wantDispName)
		}
	})

	// disableServiceAccount test.
	err = disableServiceAccount(buf, account.Email)
	if err != nil {
		t.Fatalf("disableServiceAccount: %v", err)
	}

	// enableServiceAccount test.
	err = enableServiceAccount(buf, account.Email)
	if err != nil {
		t.Fatalf("enableServiceAccount: %v", err)
	}

	// listServiceAccounts test.
	accounts, err := listServiceAccounts(buf, tc.ProjectID)
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
