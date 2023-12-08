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

package snippets

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	speech "cloud.google.com/go/speech/apiv2"
	"cloud.google.com/go/speech/apiv2/speechpb"
)

func transcribeFileV2(w io.Writer, projectId, audioFile string) error {
	ctx := context.Background()
	client, err := speech.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("speech.NewClient: %w", err)
	}
	defer client.Close()

	content, err := os.ReadFile(audioFile)

	if err != nil {
		return fmt.Errorf("file: %s, ReadFile Error: %w", audioFile, err)
	}

	languageCodes := []string{"en-US"}
	config := &speechpb.RecognitionConfig{
		DecodingConfig: &speechpb.RecognitionConfig_AutoDecodingConfig{},
		Model:          "short",
		LanguageCodes:  languageCodes,
	}
	request := &speechpb.RecognizeRequest{
		Recognizer: fmt.Sprintf("projects/%s/locations/global/recognizers/_", projectId),
		Config:     config,
		AudioSource: &speechpb.RecognizeRequest_Content{
			Content: content,
		},
	}

	response, err := client.Recognize(
		ctx, request,
	)

	if err != nil {
		return fmt.Errorf("file: %s, Recognize: %w", audioFile, err)
	}

	for i, result := range response.Results {
		fmt.Fprintf(w, "%s\n", strings.Repeat("-", 20))
		fmt.Fprintf(w, "Result %d\n", i+1)
		for j, alternative := range result.Alternatives {
			fmt.Fprintf(w, "Alternative %d: %s\n", j+1, alternative.Transcript)
		}
	}

	return nil
}
