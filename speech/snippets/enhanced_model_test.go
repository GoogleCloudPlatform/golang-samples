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

package snippets

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestEnhancedModel(t *testing.T) {
	testutil.SystemTest(t)

	var buf bytes.Buffer
	err := enhancedModel(&buf)
	if err != nil {
		t.Fatalf("%v - You may need to enable data logging. See https://cloud.google.com/speech-to-text/docs/enable-data-logging", err)
	}

	if got := buf.String(); !strings.Contains(got, "Chrome") {
		t.Fatalf(`enhancedModel(../testdata/commercial_mono.wav) = %q; want "Chrome"`, got)
	}
}
