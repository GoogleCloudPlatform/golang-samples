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

package language_snippets

import (
	"context"

	language "cloud.google.com/go/language/apiv1"
	"cloud.google.com/go/language/apiv1/languagepb"
)

// Functions are included in the documentation. Use a package level variable
// to avoid it being included in the function signature below.
var client *language.Client

// [START language_entities_gcs]

func analyzeEntitiesFromGCS(ctx context.Context, gcsURI string) (*languagepb.AnalyzeEntitiesResponse, error) {
	return client.AnalyzeEntities(ctx, &languagepb.AnalyzeEntitiesRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_GcsContentUri{
				GcsContentUri: gcsURI,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}

// [END language_entities_gcs]

// [START language_sentiment_gcs]

func analyzeSentimentFromGCS(ctx context.Context, gcsURI string) (*languagepb.AnalyzeSentimentResponse, error) {
	return client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_GcsContentUri{
				GcsContentUri: gcsURI,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
}

// [END language_sentiment_gcs]

// [START language_syntax_gcs]

func analyzeSyntaxFromGCS(ctx context.Context, gcsURI string) (*languagepb.AnnotateTextResponse, error) {
	return client.AnnotateText(ctx, &languagepb.AnnotateTextRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_GcsContentUri{
				GcsContentUri: gcsURI,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		Features: &languagepb.AnnotateTextRequest_Features{
			ExtractSyntax: true,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}

// [END language_syntax_gcs]

// [START language_classify_gcs]

func classifyTextFromGCS(ctx context.Context, gcsURI string) (*languagepb.ClassifyTextResponse, error) {
	return client.ClassifyText(ctx, &languagepb.ClassifyTextRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_GcsContentUri{
				GcsContentUri: gcsURI,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
}

// [END language_classify_gcs]
