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
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestBatchPredict(t *testing.T) {
	tc := testutil.SystemTest(t)
	var buf bytes.Buffer
	var r *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	ctx := context.Background()
	bucketName := testutil.TestBucket(ctx, t, tc.ProjectID, "golang-samples-batch")
	location := "us-central1"
	outputURI := fmt.Sprintf("gs://%s/", bucketName)
	inputURIs := []string{"gs://cloud-samples-data/generative-ai/embeddings/embeddings_input.jsonl"}
	name := fmt.Sprintf("test-job-go-batch-%v-%v", time.Now().Format("2006-01-02"), r.Int())

	err := embedBatchPredict(&buf, tc.ProjectID, location, name, outputURI, inputURIs)
	if err != nil {
		t.Error(err)
	}

	output := buf.String()
	if output != name {
		t.Errorf("job name doesn't match. Got: %s, want: %s", output, name)
	}
}
