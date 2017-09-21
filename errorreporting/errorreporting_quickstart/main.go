// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START errorreporting_quickstart]

// Sample errorreporting_quickstart contains is a quickstart
// example for the Google Cloud Error Reporting API.
package errorreporting_quickstart

import (
	"log"

	"cloud.google.com/go/errorreporting"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "YOUR_PROJECT_ID"

	errorClient, err := errorreporting.NewClient(ctx, projectID, "myservice", "v1.0", false)
	if err != nil {
		log.Fatal(err)
	}
	defer errorClient.Close()

	// Report panics.
	defer errorClient.Catch(ctx)

	// Your program here...
}

// [END errorreporting_quickstart]
