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

package testutil

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/broady/preprocess/lib/preprocess"
	"golang.org/x/tools/imports"
)

type genTest struct {
	t         *testing.T
	template  string
	labels    []string
	goimports bool
}

func Generated(t *testing.T, templateFile string) genTest {
	return genTest{
		t:        t,
		template: templateFile,
	}
}

func (g genTest) Labels(labels ...string) genTest {
	g.labels = labels
	return g
}

func (g genTest) Goimports() genTest {
	g.goimports = true
	return g
}

func (g genTest) Matches(outFile string) {
	tmpl, err := ioutil.ReadFile(g.template)
	if err != nil {
		g.t.Errorf("ReadFile(%v): %v", g.template, err)
		return
	}

	want, err := ioutil.ReadFile(outFile)
	if err != nil {
		g.t.Errorf("ReadFile(%v): %v", outFile, err)
		return
	}

	got, err := preprocess.Process(bytes.NewReader(tmpl), g.labels, "//#")
	if err != nil {
		g.t.Errorf("Preprocess(%v, %v): %v", g.template, g.labels, err)
		return
	}

	if g.goimports {
		got, err = imports.Process(outFile, got, nil)
		if err != nil {
			g.t.Errorf("Goimports(%v): %v", outFile, err)
			return
		}
	}

	if !bytes.Equal(got, want) {
		g.t.Errorf("Generated output for %s doesn't match file on disk. Did you edit the generated file instead of the template, or forget to run `go generate`?", outFile)
	}
}
