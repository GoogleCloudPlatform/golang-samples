// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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

func TestPredictGemma2(t *testing.T) {
	tc := testutil.SystemTest(t)

	projectID := tc.ProjectID
	var buf bytes.Buffer
	client := PredictionsClient{}

	t.Run("GPU predict", func(t *testing.T) {
		buf.Reset()
		// Mock ID used to check if GPU was called
		if err := predictGPU(&buf, client, projectID, GPUEndpointRegion, GPUEndpointID); err != nil {
			t.Fatal(err)
		}

		if got := buf.String(); !strings.Contains(got, "Rayleigh scattering") {
			t.Error("generated text content not found in response")
		}
	})

	t.Run("TPU predict", func(t *testing.T) {
		buf.Reset()
		// Mock ID used to check if TPU was called
		if err := predictTPU(&buf, client, projectID, TPUEndpointRegion, TPUEndpointID); err != nil {
			t.Fatal(err)
		}

		if got := buf.String(); !strings.Contains(got, "Rayleigh scattering") {
			t.Error("generated text content not found in response")
		}
	})
}
