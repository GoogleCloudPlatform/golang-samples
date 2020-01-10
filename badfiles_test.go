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

package samples

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bmatcuk/doublestar"
	"github.com/h2non/filetype"
)

// alwaysBad contains glob patterns that should always be rejected.
var alwaysBad = []string{
	"**/*.swp",
}

// allowList contains glob patterns that are acceptable.
var allowList = []string{
	// Files that are always good!
	"**/*.go",
	"**/*.md",
	"**/*.yaml",
	"**/*.sh",
	"**/*.bash",
	"**/*.mod",
	"**/*.sum",
	"**/*.svg",
	"**/*.tmpl",
	"**/*.css",
	"**/*.html",
	"**/*.js",
	"**/*.sql",
	"**/*.dot",
	"**/*.proto",

	"LICENSE",
	"**/*Dockerfile*",
	"**/.dockerignore",
	"**/Makefile",

	// Primarily ML APIs.
	"**/testdata/**/*.jpg",
	"**/testdata/**/*.wav",
	"**/testdata/**/*.raw",
	"**/testdata/**/*.png",
	"**/testdata/**/*.txt",
	"**/testdata/**/*.csv",

	// Healthcare data.
	"healthcare/testdata/dicom_00000001_000.dcm",
	"healthcare/testdata/hl7v2message.dat",

	// Webrisk samples.
	"webrisk/non_existing_path.path",
	"webrisk/internal/webrisk_proto/*.proto",
	"webrisk/testdata/hashes.gob",

	// Endpoints samples.
	"endpoints/**/*.proto",

	// Cloud Functions codelab picture.
	"functions/codelabs/gopher/gophercolor.png",

	// Cloud Functions configs.
	"functions/ocr/app/config.json",
	"functions/slack/config.json",

	// Samples that aren't really code. Legacy.
	"**/appengine/**/*.txt",

	// Test output and configs.
	"gotest.out",
	"testing/kokoro/*.cfg",

	// TODO: cruft that should probably be under "testdata".
	"appengine_flexible/pubsub/sample_message.json",
	"dialogflow/resources/**/*",
	"texttospeech/**/*",
	"storage/objects/notes.txt",
	"videointelligence/resources/**/*",

	// Renovate configuration.
	".github/renovate.json",

	// Getting Started on GCE systemd service file.
	"**/gce/**/*.service",
}

// Check whether accidental binary files have been checked in.
func TestBadFiles(t *testing.T) {
	err := filepath.Walk(".", func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			if fi.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		for _, pattern := range alwaysBad {
			match, err := doublestar.PathMatch(pattern, path)
			if err != nil {
				t.Fatalf("bad pattern: %q", pattern)
			}
			if match {
				t.Errorf("Bad file checked in: %v", path)
				return nil
			}
		}

		for _, pattern := range allowList {
			match, err := doublestar.PathMatch(pattern, path)
			if err != nil {
				t.Fatalf("bad pattern: %q", pattern)
			}
			if match {
				return nil
			}
		}

		ft, _ := filetype.MatchFile(path)
		t.Errorf("Likely bad file checked in: %v. MIME type: %s", path, ft.MIME)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
