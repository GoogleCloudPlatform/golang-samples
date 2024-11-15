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

package embeddings

import (
	"bytes"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestEmbeddings(t *testing.T) {
	tc := testutil.SystemTest(t)

	var buf bytes.Buffer
	location := "us-central1"

	t.Run("generate embeddings with lower dimension", func(t *testing.T) {
		buf.Reset()
		res, err := generateWithLowerDimension(&buf, tc.ProjectID, location)
		if err != nil {
			t.Fatalf("generateWithLowerDimension failed: %v", err)
		}

		if res == nil {
			t.Error("received empty response")
		}

		expEmbeddingLen := 128
		for _, embedding := range res {
			if len(embedding) != expEmbeddingLen {
				t.Errorf("expected an embedding of len %d, got len %d", expEmbeddingLen, len(embedding))
			}
		}
	})

	t.Run("generate embeddings for image and text", func(t *testing.T) {
		buf.Reset()
		res, err := generateForTextAndImage(&buf, tc.ProjectID, location)
		if err != nil {
			t.Fatalf("generateForImageAndText failed: %v", err)
		}

		if res == nil {
			t.Error("received empty response")
		}

		expEmbeddingLen := 1408
		for _, embedding := range res {
			if len(embedding) != expEmbeddingLen {
				t.Errorf("expected an embedding of len %d, got len %d", expEmbeddingLen, len(embedding))
			}
		}
	})

	t.Run("generate embedding for image", func(t *testing.T) {
		buf.Reset()
		res, err := generateForImage(&buf, tc.ProjectID, location)
		if err != nil {
			t.Fatalf("generateForImage failed: %v", err)
		}

		if res == nil {
			t.Error("received empty response")
		}

		expEmbeddingLen := 1408
		if len(res) != expEmbeddingLen {
			t.Errorf("expected an embedding of len %d, got len %d", expEmbeddingLen, len(res))
		}
	})

	t.Run("generate embedding for video", func(t *testing.T) {
		buf.Reset()
		res, err := generateForVideo(&buf, tc.ProjectID, location)
		if err != nil {
			t.Fatalf("generateForVideo failed: %v", err)
		}

		if res == nil {
			t.Error("received empty response")
		}

		expEmbeddingLen := 1408
		if len(res) != expEmbeddingLen {
			t.Errorf("expected an embedding of len %d, got len %d", expEmbeddingLen, len(res))
		}
	})

	t.Run("generate embeddings for image text and video", func(t *testing.T) {
		buf.Reset()
		res, err := generateForImageTextAndVideo(&buf, tc.ProjectID, location)
		if err != nil {
			t.Fatalf("generateForImageTextAndVideo failed: %v", err)
		}

		if res == nil {
			t.Error("received empty response")
		}

		expEmbeddingLen := 1408
		for _, embedding := range res {
			if len(embedding) != expEmbeddingLen {
				t.Errorf("expected an embedding of len %d, got len %d", expEmbeddingLen, len(embedding))
			}
		}
	})
}
