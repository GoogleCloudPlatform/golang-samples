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
	"strings"
	"testing"

	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"

	languagev1beta2 "cloud.google.com/go/language/apiv1beta2"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestSentiment(t *testing.T) {
	testutil.SystemTest(t)
	ctx, c := newClient(t)

	res, err := analyzeSentiment(ctx, c, "I am very happy.")
	if err != nil {
		t.Fatalf("got %v, want nil err", err)
	}
	if got := res.DocumentSentiment.Score; got <= 0 {
		t.Errorf("sentiment score: got %f, want positive", got)
	}
}

func TestEntitySentiment(t *testing.T) {
	testutil.SystemTest(t)
	ctx, c := newBetaClient(t)

	res, err := analyzeEntitySentiment(ctx, c, "Marvin Gaye is the best")
	if err != nil {
		t.Fatalf("got %v, want nil err", err)
	}
	for _, e := range res.Entities {
		if e.Name == "Marvin Gaye" {
			if got := e.Sentiment.Score; got <= 0 {
				t.Errorf("Marvin Gaye sentiment score: got %f, want positive", got)
			}
			return // found
		}
	}
	t.Errorf("got %+v; want Marvin Gaye in Entities", res)
}

func TestEntity(t *testing.T) {
	testutil.SystemTest(t)
	ctx, c := newClient(t)

	res, err := analyzeEntities(ctx, c, "Homer Simpson likes donuts.")
	if err != nil {
		t.Fatalf("got %v, want nil err", err)
	}
	for _, e := range res.Entities {
		if e.Name == "Homer Simpson" {
			return // found
		}
	}
	t.Errorf("got %+v; want Homer in Entities", res)
}

func TestSyntax(t *testing.T) {
	testutil.SystemTest(t)
	ctx, c := newClient(t)

	res, err := analyzeSyntax(ctx, c, "If you bend the gopher, his belly folds.")
	if err != nil {
		t.Fatalf("got %v, want nil err", err)
	}

	for _, tok := range res.Tokens {
		if tok.Lemma == "gopher" {
			if tok.PartOfSpeech.Tag != languagepb.PartOfSpeech_NOUN {
				t.Errorf("PartOfSpeech: got %+v, want NOUN", tok.PartOfSpeech.Tag)
			}
			return // found
		}
	}
	t.Errorf("got %+v; want gopher in Tokens", res)
}

func TestClassify(t *testing.T) {
	testutil.SystemTest(t)
	ctx, c := newClient(t)

	res, err := classifyText(ctx, c, "Android is a mobile operating system developed by Google, based on the Linux kernel and designed primarily for touchscreen mobile devices such as smartphones and tablets.")
	if err != nil {
		t.Fatalf("got %v, want nil err", err)
	}
	for _, category := range res.Categories {
		if strings.Contains(category.GetName(), "Computers") {
			return // found
		}
	}
	t.Errorf("got %+v; want Computers in Categories", res)
}

func newClient(t *testing.T) (context.Context, *language.Client) {
	ctx := context.Background()
	client, err := language.NewClient(ctx)
	if err != nil {
		t.Fatalf("language.NewClient: %v", err)
	}
	return ctx, client
}

func newBetaClient(t *testing.T) (context.Context, *languagev1beta2.Client) {
	ctx := context.Background()
	client, err := languagev1beta2.NewClient(ctx)
	if err != nil {
		t.Fatalf("language.NewClient: %v", err)
	}
	return ctx, client
}
