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

package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iam/v1"
)

func TestQuickstart(t *testing.T) {
	tc := testutil.SystemTest(t)
	m := testutil.BuildMain(t)

	if !m.Built() {
		t.Fatalf("failed to build app")
	}

	// create a service account to use in the test
	testServiceAccount, err := createServiceAccount(tc.ProjectID, "iam-quickstart-service-account", "IAM test account")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	t.Cleanup(func() {
		if err := deleteServiceAccount(testServiceAccount.Email); err != nil {
			t.Fatalf("cleanup failed: %v", err)
		}
	})

	testMember := "serviceAccount:" + testServiceAccount.Email

	stdOut, stdErr, err := m.Run(nil, 10*time.Minute,
		"--project_id", tc.ProjectID,
		"--member_id", testMember,
	)

	if err != nil {
		t.Errorf("stdout: %v", string(stdOut))
		t.Errorf("stderr: %v", string(stdErr))
		t.Errorf("execution failed: %v", err)
	}

	if got := string(stdOut); !strings.Contains(got, testMember) {
		t.Errorf("got %q, want to contain %q", got, testMember)
	}
}

// createServiceAccount creates a service account.
func createServiceAccount(projectID, name, displayName string) (*iam.ServiceAccount, error) {
	ctx := context.Background()
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %w", err)
	}

	request := &iam.CreateServiceAccountRequest{
		AccountId: name,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: displayName,
		},
	}
	account, err := service.Projects.ServiceAccounts.Create("projects/"+projectID, request).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.Create: %w", err)
	}
	return account, nil
}

// deleteServiceAccount deletes a service account.
func deleteServiceAccount(email string) error {
	ctx := context.Background()
	service, err := iam.NewService(ctx)
	if err != nil {
		return fmt.Errorf("iam.NewService: %w", err)
	}

	if _, err := service.Projects.ServiceAccounts.Delete("projects/-/serviceAccounts/" + email).Do(); err != nil {
		return fmt.Errorf("Projects.ServiceAccounts.Delete: %w", err)
	}
	return nil
}
