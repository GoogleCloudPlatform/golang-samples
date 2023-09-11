// Copyright 2023 Google LLC
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

package language_v2

// [START language_classify_text]
import (
	"context"
	"fmt"
	"io"

	language "cloud.google.com/go/language/apiv2"
	"cloud.google.com/go/language/apiv2/languagepb"
)

// classifyText sends a string of text to the Cloud Natural Language API to
// classify the category of the text.
func classifyText(w io.Writer, text string) error {
	ctx := context.Background()

	// Initialize client.
	client, err := language.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	resp, err := client.ClassifyText(ctx, &languagepb.ClassifyTextRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})

	if err != nil {
		return fmt.Errorf("ClassifyText: %w", err)
	}
	fmt.Fprintf(w, "Response: %q\n", resp)

	return nil
}

// [END language_classify_text]
