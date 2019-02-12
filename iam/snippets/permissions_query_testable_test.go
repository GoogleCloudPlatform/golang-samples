// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package snippets

import (
	"bytes"
	"os"
	"testing"
)

func TestQueryTestablePermissions(t *testing.T) {
	buf := &bytes.Buffer{}
	project := os.Getenv("GOLANG_SAMPLES_PROJECT_ID")
	name := "//cloudresourcemanager.googleapis.com/projects/" + project
	permissions, err := queryTestablePermissions(buf, name)
	if err != nil {
		t.Fatalf("queryTestablePermissions: %v", err)
	}
	if len(permissions) < 1 {
		t.Fatalf("queryTestablePermissions: expected at least 1 item")
	}
}
