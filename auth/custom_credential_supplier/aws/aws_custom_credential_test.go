// Copyright 2025 Google LLC
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
	"testing"

	"cloud.google.com/go/auth/credentials/externalaccount"
)

// TestCustomAwsSupplier_AwsRegion verifies that the supplier correctly
// resolves the Region from environment variables, respecting precedence.
func TestCustomAwsSupplier_AwsRegion(t *testing.T) {
	ctx := context.Background()
	supplier := &customAwsSupplier{}
	opts := &externalaccount.RequestOptions{}

	tests := []struct {
		name       string
		env        map[string]string
		wantRegion string
	}{
		{
			name: "AWS_REGION is set (Highest Priority)",
			env: map[string]string{
				"AWS_REGION":         "us-west-1",
				"AWS_DEFAULT_REGION": "us-east-1",
			},
			wantRegion: "us-west-1",
		},
		{
			name: "Only AWS_DEFAULT_REGION is set (Fallback)",
			env: map[string]string{
				"AWS_REGION":         "",
				"AWS_DEFAULT_REGION": "eu-central-1",
			},
			wantRegion: "eu-central-1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables for this test case
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			got, err := supplier.AwsRegion(ctx, opts)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if got != tc.wantRegion {
				t.Errorf("AwsRegion() = %v, want %v", got, tc.wantRegion)
			}
		})
	}
}

// TestCustomAwsSupplier_AwsSecurityCredentials verifies that the supplier
// correctly extracts credentials from the AWS environment variables.
func TestCustomAwsSupplier_AwsSecurityCredentials(t *testing.T) {
	ctx := context.Background()
	supplier := &customAwsSupplier{}
	opts := &externalaccount.RequestOptions{}

	// Mock AWS Credentials via Environment Variables
	// The AWS SDK v2 automatically picks these up.
	expectedID := "AKIA_TEST_ACCESS_KEY"
	expectedSecret := "TEST_SECRET_KEY"
	expectedToken := "TEST_SESSION_TOKEN"

	t.Setenv("AWS_ACCESS_KEY_ID", expectedID)
	t.Setenv("AWS_SECRET_ACCESS_KEY", expectedSecret)
	t.Setenv("AWS_SESSION_TOKEN", expectedToken)
	t.Setenv("AWS_REGION", "us-east-1") // Required for SDK config to load successfully

	creds, err := supplier.AwsSecurityCredentials(ctx, opts)
	if err != nil {
		t.Fatalf("AwsSecurityCredentials failed: %v", err)
	}

	if creds.AccessKeyID != expectedID {
		t.Errorf("AccessKeyID = %v, want %v", creds.AccessKeyID, expectedID)
	}
	if creds.SecretAccessKey != expectedSecret {
		t.Errorf("SecretAccessKey = %v, want %v", creds.SecretAccessKey, expectedSecret)
	}
	if creds.SessionToken != expectedToken {
		t.Errorf("SessionToken = %v, want %v", creds.SessionToken, expectedToken)
	}
}