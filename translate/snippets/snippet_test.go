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

	in := "The Go Gopher is cute"
	got, err := translateText("ja", in)
	if err != nil {
		t.Fatalf("translateText: %v", err)
	}
	if jaGo, jaKawaii := "ゴー", "かわいい"; !(strings.Contains(got, jaGo) || strings.Contains(got, jaKawaii)) {
		t.Errorf("translateText(%q)=%q; want to contain %q or %q", in, got, jaGo, jaKawaii)
	}
}

func TestTranslateWithModel(t *testing.T) {
	t.Skip("Project must be whitelisted")

	testutil.SystemTest(t)

	in := "The Go Gopher is cute"
	got, err := translateTextWithModel("ja", in, "nmt")
	if err != nil {
		t.Fatalf("translateText: %v", err)
	}
	if want := "ゴー"; !strings.Contains(got, want) {
		t.Errorf("translateText(%q)=%q; want to contain %q", in, got, want)
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
