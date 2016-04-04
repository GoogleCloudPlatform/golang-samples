// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package samples

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Break up string so the current file is also tested.
const sentinel = "Google Inc. All rights reserved.\n" +
	"// Use of this source code is governed by the Apache 2.0\n" +
	"// license that can be found in the LICENSE file."

const prefix = "// Copyright "

func TestLicense(t *testing.T) {
	err := filepath.Walk(".", func(path string, fi os.FileInfo, err error) error {
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
		if !bytes.Contains(src, []byte(sentinel)) {
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
