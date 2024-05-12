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

// [START aiplatform_gemini_multiturn_chat]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"cloud.google.com/go/vertexai/genai"
)

func makeChatRequests(w io.Writer, projectID string, location string, modelName string) error {
	// location := "us-central1"
	// modelName := "gemini-1.0-pro-002"
	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	defer client.Close()

	gemini := client.GenerativeModel(modelName)
	chat := gemini.StartChat()

	r, err := chat.SendMessage(
		ctx,
		genai.Text("Hello"))
	if err != nil {
		return err
	}
	rb, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}
	fmt.Fprintln(w, string(rb))

	r, err = chat.SendMessage(
		ctx,
		genai.Text("What are all the colors in a rainbow?"))
	if err != nil {
		return err
	}
	rb, err = json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}
	fmt.Fprintln(w, string(rb))

	r, err = chat.SendMessage(
		ctx,
		genai.Text("Why does it appear when it rains?"))
	if err != nil {
		return fmt.Errorf("chat.SendMessage: %w", err)
	}
	rb, err = json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}
	fmt.Fprintln(w, string(rb))

	return nil
}

// [END aiplatform_gemini_multiturn_chat]
