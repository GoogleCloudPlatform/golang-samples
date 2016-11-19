// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package language_snippets

import (
	"golang.org/x/net/context"

	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

// Functions are included in the documentation. Use a package level variable
// to avoid it being included in the function signature below.
var client *language.Client

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
