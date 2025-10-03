// Copyright 2025 Google LLC
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

package tuning

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const gcsOutputBucket = "genai-tests"

func TestTuningGeneration(t *testing.T) {
	tc := testutil.SystemTest(t)

	t.Setenv("GOOGLE_GENAI_USE_VERTEXAI", "1")
	t.Setenv("GOOGLE_CLOUD_LOCATION", "us-central1")
	t.Setenv("GOOGLE_CLOUD_PROJECT", tc.ProjectID)

	prefix := fmt.Sprintf("tuning-output/%d", time.Now().UnixNano())
	outputURI := fmt.Sprintf("gs://%s/%s", gcsOutputBucket, prefix)
	buf := new(bytes.Buffer)

	t.Run("create tuning job in project", func(t *testing.T) {
		buf.Reset()
		err := createTuningJob(buf, outputURI)
		if err != nil {
			t.Fatalf("createTuningJob failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("predictWithTunedEndpoint in project", func(t *testing.T) {
		buf.Reset()
		err := predictWithTunedEndpoint(buf, "test-tuning-job")
		if err != nil {
			t.Fatalf("predictWithTunedEndpoint failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})

	t.Run("get Tuning Job in project", func(t *testing.T) {
		buf.Reset()
		err := getTuningJob(buf, "test-tuning-job")
		if err != nil {
			t.Fatalf("getTuningJob failed: %v", err)
		}

		output := buf.String()
		if output == "" {
			t.Error("expected non-empty output, got empty")
		}
	})
}
