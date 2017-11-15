// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

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
