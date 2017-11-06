// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"strings"
	"testing"

	"golang.org/x/net/context"

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
