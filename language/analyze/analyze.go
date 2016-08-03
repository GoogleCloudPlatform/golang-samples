// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command analyze performs sentiment, entity, and syntax analysis
// on a string of text via the Cloud Natural Language API.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	language "google.golang.org/api/language/v1beta1"
)

func main() {
	if len(os.Args) < 2 {
		usage("Missing command.")
	}

	ctx := context.Background()
	hc, err := google.DefaultClient(ctx, language.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}
	client, err := language.New(hc)
	if err != nil {
		log.Fatal(err)
	}

	text := strings.Join(os.Args[2:], " ")
	if text == "" {
		usage("Missing text.")
	}

	switch os.Args[1] {
	case "entities":
		printResp(analyzeEntities(client, text))
	case "sentiment":
		printResp(analyzeSentiment(client, text))
	case "syntax":
		printResp(analyzeSyntax(client, text))
	default:
		usage("Unknown command.")
	}
}

func usage(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Fprintln(os.Stderr, "usage: analyze [entities|sentiment|syntax] <text>")
	os.Exit(2)
}

func analyzeEntities(s *language.Service, text string) (*language.AnalyzeEntitiesResponse, error) {
	req := s.Documents.AnalyzeEntities(&language.AnalyzeEntitiesRequest{
		Document: &language.Document{
			Content: text,
			Type:    "PLAIN_TEXT",
		},
		EncodingType: "UTF8",
	})
	return req.Do()
}

func analyzeSentiment(s *language.Service, text string) (*language.AnalyzeSentimentResponse, error) {
	req := s.Documents.AnalyzeSentiment(&language.AnalyzeSentimentRequest{
		Document: &language.Document{
			Content: text,
			Type:    "PLAIN_TEXT",
		},
	})
	return req.Do()
}

func analyzeSyntax(s *language.Service, text string) (*language.AnnotateTextResponse, error) {
	req := s.Documents.AnnotateText(&language.AnnotateTextRequest{
		Document: &language.Document{
			Content: text,
			Type:    "PLAIN_TEXT",
		},
		Features: &language.Features{
			ExtractSyntax: true,
		},
		EncodingType: "UTF8",
	})
	return req.Do()
}

func printResp(v interface{}, err error) {
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(b)
}
