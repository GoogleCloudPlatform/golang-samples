// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command analyze performs sentiment, entity, and syntax analysis
// on a string of text via the Cloud Natural Language API.
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"

	"golang.org/x/net/context"

	// [START imports]
	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
	// [END imports]
)

func main() {
	if len(os.Args) < 2 {
		usage("Missing command.")
	}

	// [START init]
	ctx := context.Background()
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// [END init]

	text := strings.Join(os.Args[2:], " ")
	if text == "" {
		usage("Missing text.")
	}

	switch os.Args[1] {
	case "entities":
		printResp(analyzeEntities(ctx, client, text))
	case "sentiment":
		printResp(analyzeSentiment(ctx, client, text))
	case "syntax":
		printResp(analyzeSyntax(ctx, client, text))
	default:
		usage("Unknown command.")
	}
}

func usage(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Fprintln(os.Stderr, "usage: analyze [entities|sentiment|syntax] <text>")
	os.Exit(2)
}

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

func analyzeSyntax(ctx context.Context, client *language.Client, text string) (*languagepb.AnnotateTextResponse, error) {
	return client.AnnotateText(ctx, &languagepb.AnnotateTextRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		Features: &languagepb.AnnotateTextRequest_Features{
			ExtractSyntax: true,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}

func printResp(v proto.Message, err error) {
	if err != nil {
		log.Fatal(err)
	}
	proto.MarshalText(os.Stdout, v)
}
