// Copyright 2021 Google LLC
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

package connectionpool

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestConfigureConnectionPool(t *testing.T) {
	project := os.Getenv("GOLANG_SAMPLES_BIGTABLE_PROJECT")
	instance := os.Getenv("GOLANG_SAMPLES_BIGTABLE_INSTANCE")
	if project == "" || instance == "" {
		t.Skip("Skipping bigtable integration test. Set GOLANG_SAMPLES_BIGTABLE_PROJECT and GOLANG_SAMPLES_BIGTABLE_INSTANCE.")
	}

	buf := new(bytes.Buffer)
	if err := configureConnectionPool(buf, project, instance); err != nil {
		t.Errorf("configureConnectionPool: %v", err)
	}

	if got, want := buf.String(), "Connected with pool size of 10"; !strings.Contains(got, want) {
		t.Errorf("configureConnectionPool got %q, want %q", got, want)
	}
}
