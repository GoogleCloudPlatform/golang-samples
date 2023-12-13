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

	"cloud.google.com/go/vertexai/genai"
)

var projectId = "PROJECT_ID"
var region = "us-central1"
var modelName = "gemini-pro-vision"

func makeChatRequests(projectId string, region string, modelName string) error {
	client, err := genai.NewClient(context.TODO(), projectId, region)
	if err != nil {
		return fmt.Errorf("error creating client: %v", err)
	}
	defer client.Close()

	gemini := client.GenerativeModel(modelName)
	chat := gemini.StartChat()

	r, err := chat.SendMessage(
		context.TODO(),
		genai.Text("Hello"))
	if err != nil {
		return err
	}
	rb, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(rb))

	r, err = chat.SendMessage(
		context.TODO(),
		genai.Text("What are all the colors in a rainbow?"))
	if err != nil {
		return err
	}
	rb, _ = json.MarshalIndent(r, "", "  ")
	fmt.Println(string(rb))

	r, err = chat.SendMessage(
		context.TODO(),
		genai.Text("Why does it appear when it rains?"))
	if err != nil {
		return err
	}
	rb, _ = json.MarshalIndent(r, "", "  ")
	fmt.Println(string(rb))

	return nil
}

// [END aiplatform_gemini_multiturn_chat]
