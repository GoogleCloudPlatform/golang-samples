// Copyright 2026 Google LLC
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

// Package tools shows examples of various tools that Gemini model can use to generate text.
package tools

// [START aiplatform_genai_func_calling_config]
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	genai "google.golang.org/genai"
)

func generateWithFuncCallConfig(w io.Writer) error {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		HTTPOptions: genai.HTTPOptions{APIVersion: "v1"},
	})
	if err != nil {
		return fmt.Errorf("failed to create genai client: %w", err)
	}

	getAlbumSalesFunc := &genai.FunctionDeclaration{
		Name:        "get_album_sales",
		Description: "Gets the number of albums sold",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"albums": {
					Type:        genai.TypeArray,
					Description: "List of albums",
					Items: &genai.Schema{
						Type:        genai.TypeObject,
						Description: "Album and its sales",
						Properties: map[string]*genai.Schema{
							"album_name": {
								Type:        genai.TypeString,
								Description: "Name of the music album",
							},
							"copies_sold": {
								Type:        genai.TypeInteger,
								Description: "Number of copies sold",
							},
						},
					},
				},
			},
		},
	}

	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{
				FunctionDeclarations: []*genai.FunctionDeclaration{getAlbumSalesFunc},
			},
		},
		ToolConfig: &genai.ToolConfig{
			FunctionCallingConfig: &genai.FunctionCallingConfig{
				Mode:                 genai.FunctionCallingConfigModeAny,
				AllowedFunctionNames: []string{"get_album_sales"},
			},
		},
		Temperature: genai.Ptr(float32(0.0)),
	}

	promptText := `At Stellar Sounds, a music label, 2024 was a rollercoaster. 
				"Echoes of the Night," a debut synth-pop album, surprisingly sold 350,000 copies, 
				while veteran rock band "Crimson Tide's" latest, "Reckless Hearts," lagged at 120,000. 

				Their up-and-coming indie artist, "Luna Bloom's" EP, "Whispers of Dawn," secured 75,000 sales. 
				The biggest disappointment was the highly-anticipated rap album "Street Symphony" 
				only reaching 100,000 units. 

				Overall, Stellar Sounds moved over 645,000 units this year, revealing unexpected 
				trends in music consumption.`

	modelName := "gemini-2.5-flash"

	resp, err := client.Models.GenerateContent(ctx, modelName, genai.Text(promptText), config)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	funcCalls := resp.FunctionCalls()
	if len(funcCalls) > 0 {
		for _, fc := range funcCalls {
			fmt.Fprintf(w, "Function Call Detected: %s\n", fc.Name)

			jsondata, err := json.MarshalIndent(fc.Args, "", " ")
			if err != nil {
				return fmt.Errorf("failed to marshal function call args: %w", err)
			}

			fmt.Fprintln(w, jsondata)
			// Example response
			// {
			//  "albums": [
			//   {
			//    "album_name": "Echoes of the Night",
			//    "copies_sold": 350000
			//   },
			//   {
			//    "album_name": "Reckless Hearts",
			//    "copies_sold": 120000
			//   },
			//   {
			//    "album_name": "Whispers of Dawn",
			//    "copies_sold": 75000
			//   },
			//   {
			//    "album_name": "Street Symphony",
			//    "copies_sold": 100000
			//   }
			//  ]
			// }

		}
	} else {
		fmt.Fprintln(w, "No function calls were generated")
	}

	return nil
}

// [END aiplatform_genai_func_calling_config]
