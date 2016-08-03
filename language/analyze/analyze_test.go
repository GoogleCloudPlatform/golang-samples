// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	language "google.golang.org/api/language/v1beta1"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestSentiment(t *testing.T) {
	testutil.SystemTest(t)
	c := newClient(t)

	res, err := analyzeSentiment(c, "I am very happy.")
	if err != nil {
		t.Fatalf("got %v, want nil err", err)
	}
	if got := res.DocumentSentiment.Polarity; got <= 0 {
		t.Errorf("sentiment polarity: got %f, want positive", got)
	}
}

func TestEntity(t *testing.T) {
	testutil.SystemTest(t)
	c := newClient(t)

	res, err := analyzeEntities(c, "Homer Simpson likes donuts.")
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
	c := newClient(t)

	res, err := analyzeSyntax(c, "If you bend the gopher, his belly folds.")
	if err != nil {
		t.Fatalf("got %v, want nil err", err)
	}

	for _, tok := range res.Tokens {
		if tok.Lemma == "gopher" {
			if tok.PartOfSpeech.Tag != "NOUN" {
				t.Errorf("PartOfSpeech: got %+v, want NOUN", tok.PartOfSpeech.Tag)
			}
			return // found
		}
	}
	t.Errorf("got %+v; want gopher in Tokens", res)
}

func newClient(t *testing.T) *language.Service {
	ctx := context.Background()
	hc, err := google.DefaultClient(ctx, language.CloudPlatformScope)
	if err != nil {
		t.Fatalf("DefaultClient: %v", err)
	}
	client, err := language.New(hc)
	if err != nil {
		t.Fatalf("language.New: %v", err)
	}
	return client
}
