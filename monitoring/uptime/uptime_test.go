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

package uptime

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreate(t *testing.T) {
	c := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	config, err := create(buf, c.ProjectID)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	want := "Successfully"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("%q not found in output: %q", want, got)
	}
	delete(ioutil.Discard, config.GetName())
}

func TestList(t *testing.T) {
	c := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	if err := list(buf, c.ProjectID); err != nil {
		t.Fatalf("list: %v", err)
	}
	want := "Done"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("%q not found in output: %q", want, got)
	}
}

func TestListIPs(t *testing.T) {
	testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	if err := listIPs(buf); err != nil {
		t.Fatalf("listIPs: %v", err)
	}
	want := "Done"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("%q not found in output: %q", want, got)
	}
}

func TestGet(t *testing.T) {
	c := testutil.SystemTest(t)
	config, err := create(ioutil.Discard, c.ProjectID)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer delete(ioutil.Discard, config.GetName())
	buf := new(bytes.Buffer)
	got, err := get(buf, config.GetName())
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.GetDisplayName() != config.GetDisplayName() {
		t.Fatalf("display names not equal: want %q, got %q", config.GetDisplayName(), got.GetDisplayName())
	}
}

func TestUpdate(t *testing.T) {
	c := testutil.SystemTest(t)
	config, err := create(ioutil.Discard, c.ProjectID)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer delete(ioutil.Discard, config.GetName())
	buf := new(bytes.Buffer)
	displayName := "New display name"
	path := "/example.com/example"
	updated, err := update(buf, config.GetName(), displayName, path)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	want := "Successfully"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("%q not found in output: %q", want, got)
	}

	if got := updated.GetDisplayName(); got != displayName {
		t.Errorf("Display name not updated: got %q, want %q", got, displayName)
	}
	if got := updated.GetHttpCheck().GetPath(); got != path {
		t.Errorf("HTTP path not updated: got %q, want %q", got, path)
	}
}

func TestDelete(t *testing.T) {
	c := testutil.SystemTest(t)
	config, err := create(ioutil.Discard, c.ProjectID)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	buf := new(bytes.Buffer)
	delete(buf, config.GetName())
	want := "Successfully"
	if got := buf.String(); !strings.Contains(got, want) {
		t.Errorf("%q not found in output: %q", want, got)
	}
}
