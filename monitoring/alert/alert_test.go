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

package alert

import (
	"bytes"
	"context"
	"strings"
	"testing"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"google.golang.org/api/iterator"
)

func TestListAlertPolicies(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	if err := listAlertPolicies(buf, tc.ProjectID); err != nil {
		t.Fatalf("listAlertPolicies got error: %v", err)
	}
	want := "Done"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("listAlertPolicies got %v, want substring %q", got, want)
	}
}

func TestBackupPolicies(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	if err := backupPolicies(buf, tc.ProjectID); err != nil {
		t.Fatalf("backupPolicies got error: %v", err)
	}
	want := "ProjectID"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("backupPolicies got %v, want substring %q", got, want)
	}
}

func TestRestorePolicies(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	if err := backupPolicies(buf, tc.ProjectID); err != nil {
		t.Fatalf("backupPolicies got error: %v", err)
	}
	backup := strings.NewReader(buf.String())
	buf.Reset()
	if err := restorePolicies(buf, tc.ProjectID, backup); err != nil {
		t.Fatalf("restorePolicies got error: %v", err)
	}
	want := "Successfully"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("restorePolicies got %v, want substring %q", got, want)
	}
}

func TestReplaceChannels(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := monitoring.NewAlertPolicyClient(ctx)
	if err != nil {
		t.Fatalf("%v", err)
	}
	req := &monitoringpb.ListAlertPoliciesRequest{
		Name: "projects/" + tc.ProjectID,
	}
	it := client.ListAlertPolicies(ctx, req)
	a, err := it.Next()
	if err == iterator.Done {
		t.Skip("No alert policies")
	}
	if err != nil {
		t.Fatalf("ListAlertPolicies error: %v", err)
	}

	buf := new(bytes.Buffer)
	name := a.GetName()
	if err := replaceChannels(buf, tc.ProjectID, name[strings.LastIndex(name, "/")+1:], a.GetNotificationChannels()); err != nil {
		t.Fatalf("replaceChannels got error: %v", err)
	}
	want := "Successfully"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("replaceChannels got %v, want substring %q", got, want)
	}
}

func TestEnablePolicies(t *testing.T) {
	tc := testutil.SystemTest(t)
	tests := []struct{ enable bool }{{true}, {false}}
	for _, test := range tests {
		buf := new(bytes.Buffer)
		if err := enablePolicies(buf, tc.ProjectID, test.enable); err != nil {
			t.Fatalf("enablePolicies got error: %v", err)
		}
		want := "Successfully"
		if got := buf.String(); !strings.Contains(got, want) {
			t.Fatalf("enablePolicies(%v) got %v, want substring %q", test.enable, got, want)
		}
	}
}
