// Copyright 2022 Google LLC
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

package clientendpoint

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestSetClientEndpoint(t *testing.T) {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		t.Skip("GOOGLE_APPLICATION_CREDENTIALS not set")
	}

	customEndpoint := "http://localhost:8080/storage/v1/"

	var buf bytes.Buffer
	if err := setClientEndpoint(&buf, customEndpoint); err != nil {
		t.Errorf("setClientEndpoint: %s", err)
	}

	if got, want := buf.String(), "request endpoint set for the client"; !strings.Contains(got, want) {
		t.Errorf("setClientEndpoint: got %q; want to contain %q", got, want)
	}
}
