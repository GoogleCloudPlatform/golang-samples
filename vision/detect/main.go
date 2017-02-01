// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command detect uses the Vision API's label detection capabilities to find a label based on an image's content.
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
		fmt.Fprintf(os.Stderr, "Prefix a path with gs:// to refer to a file on GCS.\n")
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	path := args[0]

	samples := []struct {
		name       string
		local, gcs func(io.Writer, string) error
	}{
		{"detectFaces", detectFaces, detectFacesGCS},
		{"detectLabels", detectLabels, detectLabelsGCS},
		{"detectLandmarks", detectLandmarks, detectLandmarksGCS},
		{"detectText", detectText, detectTextGCS},
		{"detectLogos", detectLogos, detectLogosGCS},
		{"detectProperties", detectProperties, detectPropertiesGCS},
		{"detectSafeSearch", detectSafeSearch, detectSafeSearchGCS},
	}

	for _, sample := range samples {
		fmt.Println("---", sample.name)
		var err error
		if strings.HasPrefix(path, "gs://") {
			err = sample.gcs(os.Stdout, path)
		} else {
			err = sample.local(os.Stdout, path)
		}
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}
