// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package translate_snippets contains snippet code for the Translate API.
// The code is not runnable.
package translate_snippets

import (
	"bytes"
	"strings"

	"golang.org/x/text/language"

	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

func TestTranslateText(t *testing.T) {
	testutil.SystemTest(t)

	got, err := translateText("ja", "The Go Gopher is cute")
	if err != nil {
		t.Fatalf("translateText: %v", err)
	}
	if want := "可愛い"; !strings.Contains(got, want) {
		t.Errorf("translateText(Gophers are cute)=%q; want to contain %q", got, want)
	}
}

func TestDetect(t *testing.T) {
	testutil.SystemTest(t)

	detection, err := detectLanguage("可愛い")
	if err != nil {
		t.Fatalf("detectLanguage: %v", err)
	}

	if got, want := detection.Language, language.Japanese; got != want {
		t.Errorf("detection.Language: got %q; want %q", got, want)
	}
}

func TestListSupported(t *testing.T) {
	testutil.SystemTest(t)

	buf := &bytes.Buffer{}
	if err := listSupportedLanguages(buf, "th"); err != nil {
		t.Fatal(err)
	}
	if got, want := buf.String(), `"en":`; !strings.Contains(got, want) {
		t.Errorf("listSupportedLanguages(th): want %q in %q", want, got)
	}
}
