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

// Package metadata contains example snippets using the DLP info types API.
package metadata

import (
	"bytes"
	"context"
	"log"
	"os"
	"strings"
	"testing"

	dlp "cloud.google.com/go/dlp/apiv2"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

var client *dlp.Client
var projectID string

func TestMain(m *testing.M) {
	ctx := context.Background()
	if c, ok := testutil.ContextMain(m); ok {
		var err error
		client, err = dlp.NewClient(ctx)
		if err != nil {
			log.Fatalf("datastore.NewClient: %v", err)
		}
		projectID = c.ProjectID
		defer client.Close()
	}
	os.Exit(m.Run())
}

func TestInfoTypes(t *testing.T) {
	testutil.SystemTest(t)
	tests := []struct {
		language string
		filter   string
		want     string
	}{
		{
			want: "TIME",
		},
		{
			language: "en-US",
			want:     "TIME",
		},
		{
			language: "es",
			want:     "DATE",
		},
		{
			filter: "supported_by=INSPECT",
			want:   "GENDER",
		},
	}
	for _, test := range tests {
		t.Run(test.language, func(t *testing.T) {
			t.Parallel()
			buf := new(bytes.Buffer)
			err := infoTypes(buf, client, test.language, test.filter)
			if err != nil {
				t.Errorf("infoTypes(%s, %s) = error %q, want substring %q", test.language, test.filter, err, test.want)
			}
			if got := buf.String(); !strings.Contains(got, test.want) {
				t.Errorf("infoTypes(%s, %s) = %s, want substring %q", test.language, test.filter, got, test.want)
			}
		})
	}
}
