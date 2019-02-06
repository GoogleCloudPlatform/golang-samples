// Copyright 2019 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package snippets

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestViewGrantableRoles(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	resource := "//cloudresourcemanager.googleapis.com/projects/" + projectID
	viewGrantableRoles(resource)
}

func TestQueryTestablePermissions(t *testing.T) {
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	resource := "//cloudresourcemanager.googleapis.com/projects/" + projectID
	queryTestablePermissions(resource)
}

func TestRoles(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	roleName := "gotest" + strconv.Itoa(rand.Intn(100000))
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")

	createRole(roleName, projectID, "Go Test", "Go Test",
		[]string{"appengine.versions.get"}, "GA")
	editRole(roleName, projectID, "Go Test 2", "Go Test 2",
		[]string{"appengine.versions.get"}, "GA")
	disableRole(roleName, projectID)
	listRoles(projectID)
	deleteRole(roleName, projectID)
}

func TestServiceAccountsAndKeys(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	accountName := "gotest" + strconv.Itoa(rand.Intn(100000))
	projectID := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")

	account := createServiceAccount(projectID, accountName, "Test Account")
	account = renameServiceAccount(account.Email, "Updated Test Account")
	listServiceAccounts(projectID)
	key := createKey(account.Email)
	listKeys(account.Email)
	deleteKey(key.Name)
	deleteServiceAccount(account.Email)
}
