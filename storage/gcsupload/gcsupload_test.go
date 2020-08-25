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
	"bytes"
	"context"
	"crypto/md5"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestUpload(t *testing.T) {
	tc := testutil.SystemTest(t)

	ctx := context.Background()
	projectID := tc.ProjectID
	bucket := projectID + "-gcsupload"

	testutil.CleanBucket(ctx, t, projectID, bucket)

	input := strings.Repeat("GCS test\n", 30)
	r := strings.NewReader(input)

	name := "atest.txt"
	obj, objAttrs, err := upload(ctx, r, projectID, bucket, name, true)
	if err != nil {
		t.Fatalf("expected to successfully upload: %v", err)
	}
	if objAttrs == nil {
		t.Fatal("expected back a non-nil object")
	}
	defer obj.Delete(ctx)

	if g, w := objAttrs.Name, name; g != w {
		t.Errorf("name: got=%q want=%q", g, w)
	}
	if g, w := objAttrs.Size, int64(len(input)); g != w {
		t.Errorf("size: got=%d want=%d", g, w)
	}
	h := md5.New()
	h.Write([]byte(input))
	if g, w := objAttrs.MD5, h.Sum(nil); !bytes.Equal(g, w) {
		t.Errorf("md5: got=%x want=%x", g, w)
	}
}
