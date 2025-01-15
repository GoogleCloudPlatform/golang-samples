// Copyright 2023 Google LLC
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

package chat

// [START generativeaionvertexai_gemini_multiturn_chat]
// [START aiplatform_gemini_multiturn_chat]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

func makeChatRequests(ctx context.Context, w io.Writer, projectID, region, modelName string) error {
	client, err := genai.NewClient(ctx, projectID, region)

	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	defer client.Close()

	gemini := client.GenerativeModel(modelName)
	chat := gemini.StartChat()

	send := func(message string) error {
		r, err := chat.SendMessage(ctx, genai.Text(message))
		if err != nil {
			return err
		}
		rb, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(w, string(rb))
		return nil
	}

	if err := send("Hello"); err != nil {
		return err
	}
	if err := send("What are all the colors in a rainbow?"); err != nil {
		return err
	}
	return send("Why does it appear when it rains?")
}

// [END aiplatform_gemini_multiturn_chat]
// [END generativeaionvertexai_gemini_multiturn_chat]
