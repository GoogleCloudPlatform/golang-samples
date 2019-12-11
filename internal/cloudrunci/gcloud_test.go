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

package cloudrunci

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestCreateIDToken(t *testing.T) {
	testutil.EndToEndTest(t)
	// TODO assign to token
	_, err := CreateIDToken("http://example.com")
	if err != nil {
		t.Errorf("CreateIDToken: %q", err)
	}

	// validate token
}

func TestGcloud(t *testing.T) {
	testutil.EndToEndTest(t)
	out, err := gcloud("label", exec.Command(gcloudBin, "help"))
	if err != nil {
		t.Errorf("CreateIDToken: %q", err)
	}

	want := "gcloud - manage Google Cloud Platform resources and developer workflow"
	if got := string(out); !strings.Contains(got, want) {
		t.Errorf("gcloud: got (%s), want (%s)", got, want)
	}
}
