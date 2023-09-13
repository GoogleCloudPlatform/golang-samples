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

// [START error_reporting_setup_go]
// [START error_reporting_quickstart]

// Sample errorreporting_quickstart contains is a quickstart
// example for the Google Cloud Error Reporting API.
package main

import (
	"context"
	"errors"
	"log"
	"os"

	"cloud.google.com/go/errorreporting"
)

var errorClient *errorreporting.Client

func main() {
	ctx := context.Background()

	// TODO: Sets your Google Cloud Platform project ID via environment or explicitly
	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	if projectID == "" {
		projectID = "YOUR_PROJECT_ID"
	}

	var err error
	errorClient, err = errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName:    "errorreporting_quickstart",
		ServiceVersion: "0.0.0",
		OnError: func(err error) {
			log.Printf("Could not repport the error: %v", err)
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer errorClient.Close()

	err = errors.New("Something went wrong")
	if err != nil {
		logAndPrintError(err)
		return
	}
}

func logAndPrintError(err error) {
	// error context is autopopulated
	errorClient.Report(errorreporting.Entry{
		Error: err,
	})
	log.Print(err)
}

// [END error_reporting_quickstart]
// [END error_reporting_setup_go]
