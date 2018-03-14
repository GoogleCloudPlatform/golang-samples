/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"
	"os"
	"testing"

	dlp "cloud.google.com/go/dlp/apiv2"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"golang.org/x/net/context"
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
