// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package samples

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

var sentinel = regexp.MustCompile(`// Copyright \d\d\d\d Google (Inc|LLC)\. All rights reserved\.
// Use of this source code is governed by the Apache 2\.0
// license that can be found in the LICENSE file\.`)

const prefix = "// Copyright "

var skip = map[string]bool{
	// These files are based off the gRPC samples, and are under a BSD license.
	"endpoints/getting-started-grpc/client/main.go":              true,
	"endpoints/getting-started-grpc/server/main.go":              true,
	"endpoints/getting-started-grpc/helloworld/helloworld.pb.go": true,
}

func TestLicense(t *testing.T) {
	err := filepath.Walk(".", func(path string, fi os.FileInfo, err error) error {
		if skip[path] {
			return nil
		}

		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		src, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}

		// Find license
		if !sentinel.Match(src) {
			t.Errorf("%v: license header not present", path)
			return nil
		}

		// Also check it is at the top of the file.
		if !bytes.HasPrefix(src, []byte(prefix)) {
			t.Errorf("%v: license header not at the top", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
