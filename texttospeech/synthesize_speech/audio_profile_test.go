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
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestAudioProfile(t *testing.T) {
	testutil.SystemTest(t)

	const text = "hello"
	const outputFile = "output.mp3"

	os.Remove(outputFile)

	var buf bytes.Buffer
	if err := audioProfile(&buf, text, outputFile); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(outputFile)

	if got := buf.String(); !strings.Contains(got, "Audio content written to file") {
		t.Error("'Audio content written to file' not found")
	}

	stat, err := os.Stat(outputFile)
	if err != nil {
		t.Error(err)
	}

	if stat.Size() == 0 {
		t.Error("Empty output file")
	}
}
