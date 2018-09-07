// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"log"
	"os"
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
