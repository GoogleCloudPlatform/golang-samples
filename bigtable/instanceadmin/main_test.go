// Copyright 2026 Google LLC
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
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"cloud.google.com/go/bigtable"
)

func TestInstanceAdmin(t *testing.T) {
	ctx := context.Background()
	project := os.Getenv("GOLANG_SAMPLES_BIGTABLE_PROJECT")
	instance := os.Getenv("GOLANG_SAMPLES_BIGTABLE_INSTANCE")
	zone := os.Getenv("GOLANG_SAMPLES_BIGTABLE_ZONE")
	if project == "" || instance == "" || zone == "" {
		t.Skip("Skipping bigtable integration test. Set GOLANG_SAMPLES_BIGTABLE_PROJECT and GOLANG_SAMPLES_BIGTABLE_INSTANCE and GOLANG_SAMPLES_BIGTABLE_ZONE.")
	}

	cluster := "my-cluster"

	// Initialize the Instance Admin Client
	instanceAdminClient, err := bigtable.NewInstanceAdminClient(ctx, project)
	if err != nil {
		t.Fatalf("bigtable.NewInstanceAdminClient: %v", err)
	}
	defer instanceAdminClient.Close()

	buf := new(bytes.Buffer)
	if err = createInstance(buf, project, instance, cluster, zone); err != nil {
		t.Fatalf("Error creating instance: %v", err)
	}
	want := fmt.Sprintf("Instance %s created successfully.", instance)
	got := buf.String()
	if !strings.Contains(got, want) {
		t.Errorf("Unexpected output string: %q", got)
	}
	t.Logf("Instance %s created successfully.\n", instance)

	if err = instanceAdminClient.DeleteInstance(ctx, instance); err != nil {
		t.Errorf("Error deleting instance: %v", err)
	}
	t.Logf("Instance %s deleted successfully.\n", instance)
}
