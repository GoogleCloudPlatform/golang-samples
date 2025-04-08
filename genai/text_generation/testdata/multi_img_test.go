// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package text_generation_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/genai/tools/text_generation"
)

func TestGenerateWithMultiImg(t *testing.T) {
	
	dummyImagePath := filepath.Join(os.TempDir(), "latte.jpg")
	err := os.WriteFile(dummyImagePath, []byte("This is a dummy JPEG image."), 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy latte.jpg: %v", err)
	}
	defer os.Remove(dummyImagePath)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	err = os.Chdir("./")
	if err != nil {
		t.Fatalf("Failed to change directory to testdata: %v", err)
	}
	defer os.Chdir(originalDir)

	var buf bytes.Buffer
	err = text_generation.generateWithMultiImg(&buf)
	if err != nil {
		t.Errorf("generateWithMultiImg failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Errorf("generateWithMultiImg produced empty output")
	}

	t.Logf("Output:\n%s", buf.String())
}