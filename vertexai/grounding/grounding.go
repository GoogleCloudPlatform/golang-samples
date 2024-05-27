// Copyright 2024 Google LLC
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

package grounding

// [START generativeaionvertexai_grounding_public_data_basic]
import (
	"context"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

func generateTextWithGroundingWeb(w io.Writer, projectID, modelName string) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, "us-central1")
	if err != nil {
		return err
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	
}

// [END generativeaionvertexai_grounding_public_data_basic]
