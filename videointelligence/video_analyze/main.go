// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// github.com/broady/preprocess
//go:generate bash -c "cat gen/template.go | preprocess | goimports > video_analyze.go"
//go:generate bash -c "cat gen/template.go | preprocess gcs | goimports > video_analyze_gcs.go"

// Command video_analyze uses the Google Cloud Video Intelligence API to analyze a video.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-image>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "Pass either a path to a local file, or a URI.\n")
		fmt.Fprintf(os.Stderr, "Prefix a path with gs:// to refer to a file on GCS.\n")
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	path := flag.Arg(0)
	match := flag.Arg(1)

	samples := []struct {
		name       string
		local, uri func(io.Writer, string) error
	}{
		{"label", label, labelURI},
		{"shotChange", shotChange, shotChangeURI},
		{"explicitContent", explicitContent, explicitContentURI},
		{"speechTranscription", speechTranscription, speechTranscriptionURI},
	}

	for _, sample := range samples {
		if !strings.Contains(sample.name, match) {
			continue
		}
		fmt.Println("---", sample.name)
		var err error
		if strings.Contains(path, "://") {
			err = sample.uri(os.Stdout, path)
		} else {
			err = sample.local(os.Stdout, path)
		}
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
