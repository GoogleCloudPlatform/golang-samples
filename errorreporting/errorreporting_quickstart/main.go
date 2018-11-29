// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// [START error_reporting_setup_go]
// [START error_reporting_quickstart]

// Sample errorreporting_quickstart contains is a quickstart
// example for the Google Cloud Error Reporting API.
package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/errorreporting"
)

var errorClient *errorreporting.Client

func main() {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "YOUR_PROJECT_ID"

	var err error
	errorClient, err = errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: "myservice",
		OnError: func(err error) {
			log.Printf("Could not log error: %v", err)
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer errorClient.Close()

	resp, err := http.Get("not-a-valid-url")
	if err != nil {
		logAndPrintError(err)
		return
	}
	log.Print(resp.Status)
}

func logAndPrintError(err error) {
	errorClient.Report(errorreporting.Entry{
		Error: err,
	})
	log.Print(err)
}

// [END error_reporting_quickstart]
// [END error_reporting_setup_go]
