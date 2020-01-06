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

package samples

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

var sentinel = regexp.MustCompile(`// Copyright \d\d\d\d Google LLC
//
// Licensed under the Apache License, Version 2\.0 \(the "License"\);
// you may not use this file except in compliance with the License\.
// You may obtain a copy of the License at
//
//     https://www\.apache\.org/licenses/LICENSE-2\.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied\.
// See the License for the specific language governing permissions and
// limitations under the License\.`)

const prefix = "// Copyright "

var skip = map[string]bool{
	// These files are based off the gRPC samples, and are under a BSD license.
	"endpoints/getting-started-grpc/client/main.go":              true,
	"endpoints/getting-started-grpc/server/main.go":              true,
	"endpoints/getting-started-grpc/helloworld/helloworld.pb.go": true,
	"run/grpc-ping/pkg/api/v1/message.pb.go":                     true,
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
