// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const (
	text       = "Hello there."
	ssml       = "<speak>Hello there.</speak>"
	outputFile = "output.mp3"
)

func TestSynthesizeText(t *testing.T) {
	testutil.SystemTest(t)

	os.Remove(outputFile)

	var buf bytes.Buffer
	err := SynthesizeText(&buf, text, outputFile)
	if err != nil {
		t.Error(err)
	}
	got := buf.String()

	if !strings.Contains(got, "Audio content written to file") {
		t.Error("'Audio content written to file' not found")
	}

	stat, err := os.Stat(outputFile)
	if err != nil {
		t.Error(err)
	}

	if stat.Size() == 0 {
		t.Error("Empty output file")
	}

	os.Remove(outputFile)
}

func TestSynthesizeSSML(t *testing.T) {
	testutil.SystemTest(t)

	os.Remove(outputFile)

	var buf bytes.Buffer
	err := SynthesizeSSML(&buf, ssml, outputFile)
	if err != nil {
		t.Error(err)
	}
	got := buf.String()

	if !strings.Contains(got, "Audio content written to file") {
		t.Error("'Audio content written to file' not found")
	}

	stat, err := os.Stat(outputFile)
	if err != nil {
		t.Error(err)
	}

	if stat.Size() == 0 {
		t.Error("Empty output file")
	}

	os.Remove(outputFile)
}
