// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package gentest

import (
	"bytes"
	"io/ioutil"
	"testing"

	"golang.org/x/tools/imports"

	"github.com/broady/preprocess/lib/preprocess"
)

func TestGen(t *testing.T) {
	in, err := ioutil.ReadFile("template.go")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		outFile string
		labels  []string
	}{
		{
			outFile: "../video_analyze.go",
			labels:  []string{},
		},
		{
			outFile: "../video_analyze_gcs.go",
			labels:  []string{"gcs"},
		},
	}
	for _, tc := range tests {
		want, err := ioutil.ReadFile(tc.outFile)
		if err != nil {
			t.Error(err)
			continue
		}

		got, err := preprocess.Process(bytes.NewReader(in), tc.labels, "//#")
		if err != nil {
			t.Errorf("%q: %v", tc.outFile, err)
			continue
		}
		got, err = imports.Process(tc.outFile, got, nil)
		if err != nil {
			t.Errorf("gofmt %q: %v", tc.outFile, err)
		}

		if !bytes.Equal(got, want) {
			t.Errorf("%s does not match. did you edit the generated file instead of the template, or forget to run `go generate`?", tc.outFile)
		}
	}
}
