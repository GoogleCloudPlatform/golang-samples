// Copyright 2019 Google LLC
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

package main

import (
	"context"
	"log"

	language "cloud.google.com/go/language/apiv1beta2"
	"cloud.google.com/go/language/apiv1beta2/languagepb"
)

func betaClient() *language.Client {
	ctx := context.Background()
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	return client
}

// [START language_entity_sentiment_text]

func analyzeEntitySentiment(ctx context.Context, client *language.Client, text string) (*languagepb.AnalyzeEntitySentimentResponse, error) {
	return client.AnalyzeEntitySentiment(ctx, &languagepb.AnalyzeEntitySentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
}

// [END language_entity_sentiment_text]
