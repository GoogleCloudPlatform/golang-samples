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

// Command analyze performs sentiment, entity, and classfication
// on a string of text via the Cloud Natural Language API.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"

	// [START language_imports]
	language "cloud.google.com/go/language/apiv2"
	"cloud.google.com/go/language/apiv2/languagepb"
	// [END language_imports]
)

func main() {
	if len(os.Args) < 2 {
		usage("Missing command.")
	}

	// [START language_init]
	ctx := context.Background()
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	// [END language_init]

	text := strings.Join(os.Args[2:], " ")
	if text == "" {
		usage("Missing text.")
	}

	switch os.Args[1] {
	case "entities":
		printResp(analyzeEntities(ctx, client, text))
	case "sentiment":
		printResp(analyzeSentiment(ctx, client, text))
	case "classify":
		printResp(classifyText(ctx, client, text))
	default:
		usage("Unknown command.")
	}
}

func usage(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Fprintln(os.Stderr, "usage: analyze [entities|sentiment|classify] <text>")
	os.Exit(2)
}

// [START language_entities_text]

func analyzeEntities(ctx context.Context, client *language.Client, text string) (*languagepb.AnalyzeEntitiesResponse, error) {
	return client.AnalyzeEntities(ctx, &languagepb.AnalyzeEntitiesRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}

// [END language_entities_text]

// [START language_sentiment_text]

func analyzeSentiment(ctx context.Context, client *language.Client, text string) (*languagepb.AnalyzeSentimentResponse, error) {
	return client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
}

// [END language_sentiment_text]

// [START language_classify_text]

func classifyText(ctx context.Context, client *language.Client, text string) (*languagepb.ClassifyTextResponse, error) {
	return client.ClassifyText(ctx, &languagepb.ClassifyTextRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
}

// [END language_classify_text]

func printResp(v proto.Message, err error) {
	if err != nil {
		log.Fatal(err)
	}
	proto.MarshalText(os.Stdout, v)
}