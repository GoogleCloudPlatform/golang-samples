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
		{"detectFaces", detectFaces, detectFacesURI},
		{"detectLabels", detectLabels, detectLabelsURI},
		{"detectLandmarks", detectLandmarks, detectLandmarksURI},
		{"detectText", detectText, detectTextURI},
		{"detectDocumentText", detectDocumentText, detectDocumentTextURI},
		{"detectLogos", detectLogos, detectLogosURI},
		{"detectProperties", detectProperties, detectPropertiesURI},
		{"detectCropHints", detectCropHints, detectCropHintsURI},
		{"detectWeb", detectWeb, detectWebURI},
		{"detectWebGeo", detectWebGeo, detectWebGeoURI},
		{"detectSafeSearch", detectSafeSearch, detectSafeSearchURI},
		{"localizeObjects", localizeObjects, localizeObjectsURI},
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
